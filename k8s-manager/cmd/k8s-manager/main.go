package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"k8s-manager/internal/config"
	"k8s-manager/internal/credentials"
	"k8s-manager/internal/healthcheck"
	"k8s-manager/internal/jenkins"
	"k8s-manager/internal/kongkey"
	"k8s-manager/internal/kube"
	"k8s-manager/internal/logger"
	"k8s-manager/internal/maintain"
	"k8s-manager/internal/rancher"
	"k8s-manager/internal/setup"
	"k8s-manager/internal/validate"

	"github.com/spf13/cobra"
)

func main() {
	var (
		kubeconfig  string
		secretsFile string
		configFile  string
		chartsDir   string
		outputDir   string
		dryRun      bool
		logLevel    string
	)

	root := &cobra.Command{
		Use:   "k8s-manager",
		Short: "CLI for managing the jeeb Kubernetes cluster",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logger.Init(logLevel)
			return nil
		},
	}

	root.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig (defaults to ~/.kube/config)")
	root.PersistentFlags().StringVar(&secretsFile, "secrets", "", "path to secrets.yaml (default: env/secrets.yaml relative to cwd)")
	root.PersistentFlags().StringVar(&configFile, "config", "", "path to config.yaml for cluster topology overrides (optional)")
	root.PersistentFlags().StringVar(&chartsDir, "charts-dir", defaultChartsDir(), "path to k8s/charts directory")
	root.PersistentFlags().StringVar(&outputDir, "output-dir", defaultOutputDir(), "directory to write vault-init.json and values-secrets.yaml")
	root.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "print commands without executing them")
	root.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log verbosity: debug, info, warn, error")

	// Keep --credentials as a deprecated alias for --secrets.
	root.PersistentFlags().StringVar(&secretsFile, "credentials", "", "")
	_ = root.PersistentFlags().MarkDeprecated("credentials", "use --secrets instead")
	_ = root.PersistentFlags().MarkHidden("credentials")

	root.AddCommand(
		newNamespaceCmd(&kubeconfig),
		newStatusCmd(&kubeconfig),
		newRestartCmd(&kubeconfig),
		newLogsCmd(&kubeconfig),
		newSetupCmd(&secretsFile, &configFile, &chartsDir, &outputDir, &dryRun),
		newDeployCmd(&secretsFile, &configFile, &chartsDir, &outputDir, &dryRun),
		newSeedCmd(&secretsFile, &configFile, &dryRun),
		newRancherCmd(&secretsFile, &configFile, &dryRun),
		newValidateCmd(&secretsFile),
		newKongKeyCmd(&secretsFile, &configFile, &chartsDir, &outputDir, &dryRun),
		newCheckCmd(&secretsFile, &configFile, &outputDir),
		newMaintainCmd(&secretsFile, &configFile, &outputDir),
		newPatchJenkinsCredsCmd(&secretsFile, &dryRun),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newSetupCmd(secretsFile, configFile, chartsDir, outputDir *string, dryRun *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Bootstrap a new cluster end-to-end (deploy charts, init Vault, patch CoreDNS)",
		Long: `Runs the full new-cluster setup in order:
  1.  Deploy jeeb-infra  (Vault, Jenkins, Nexus, SonarQube, Kong)
  2.  Deploy jeeb-data   (MongoDB, Keycloak)
  3.  Deploy jeeb-app    (backend, frontend)
  4.  Deploy jeeb-learning
  5.  Deploy jeeb-obs    (Prometheus, Loki, Grafana)
  6.  Wait for Vault pod ready
  7.  Initialize Vault   → saves vault-init.json to --output-dir
  8.  Store unseal keys  → Kubernetes secret vault-unseal-keys
  9.  Unseal Vault
  10. Configure Vault    (KV engine, policies, K8s auth roles)
  11. Patch CoreDNS      (detects ingress ClusterIP automatically)
  12. Seed Jenkins       (create seed job, generate all pipeline jobs)

Copy env/secrets.yaml.example to env/secrets.yaml and fill all values before running.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			secretsPath := resolvedSecretsFile(*secretsFile)
			creds, err := loadCreds(secretsPath)
			if err != nil {
				return err
			}
			cfg := config.LoadFromFile(resolvedConfigFile(*configFile))
			return setup.NewRunner(cfg, creds, *chartsDir, *outputDir, secretsPath, *dryRun).Run(cmd.Context())
		},
	}
}

func newDeployCmd(secretsFile, configFile, chartsDir, outputDir *string, dryRun *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "deploy [infra] [data] [app] [learning] [obs]",
		Short: "Re-deploy one or more Helm charts (no Vault init)",
		Long: `Runs helm upgrade --install for the selected charts.

Targets: infra, data, app, learning, obs (default: all five)

Examples:
  k8s-manager deploy              # deploy all charts
  k8s-manager deploy app          # deploy only jeeb-app
  k8s-manager deploy data         # deploy only jeeb-data (MongoDB + Keycloak)
  k8s-manager deploy infra obs    # deploy jeeb-infra and jeeb-obs`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, t := range args {
				switch t {
				case "infra", "data", "app", "learning", "obs":
				default:
					return fmt.Errorf("unknown target %q — valid targets: infra, data, app, learning, obs", t)
				}
			}
			secretsPath := resolvedSecretsFile(*secretsFile)
			creds, err := loadCreds(secretsPath)
			if err != nil {
				return err
			}
			cfg := config.LoadFromFile(resolvedConfigFile(*configFile))
			return setup.NewRunner(cfg, creds, *chartsDir, *outputDir, secretsPath, *dryRun).Deploy(cmd.Context(), args)
		},
	}
}

func newSeedCmd(secretsFile, configFile *string, dryRun *bool) *cobra.Command {
	var groovyPath string
	var jenkinsURL string

	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Create and run the Jenkins seed job to generate all pipeline jobs",
		Long: `Creates a Job DSL seed job in Jenkins and triggers it to generate the four
jeeb pipeline jobs (backend, frontend, learning-backend, learning-frontend).

Jenkins must already be running. The seed.groovy defaults to
jenkins/jobs/seed.groovy auto-detected from the repo root.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			secretsPath := resolvedSecretsFile(*secretsFile)
			creds, err := loadCreds(secretsPath)
			if err != nil {
				return err
			}
			cfg := config.LoadFromFile(resolvedConfigFile(*configFile))
			return jenkins.NewSeeder(cfg, creds, groovyPath, jenkinsURL, *dryRun).Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&groovyPath, "groovy-path", defaultGroovyPath(), "path to seed.groovy")
	cmd.Flags().StringVar(&jenkinsURL, "jenkins-url", "", "Jenkins URL (default: http://localhost:<nodeport>)")
	return cmd
}

func newRancherCmd(secretsFile, configFile *string, dryRun *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "rancher",
		Short: "Install cert-manager and Rancher (one-time, optional)",
		Long: `Installs cert-manager and Rancher via external Helm repos.
Run this once if you want the Rancher UI in addition to the main cluster.
Rancher is NOT required for jeeb services to run.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.LoadFromFile(resolvedConfigFile(*configFile))
			return rancher.NewDeployer(cfg, *dryRun).Run(cmd.Context())
		},
	}
}

func newValidateCmd(secretsFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Check that all credentials in secrets.yaml are filled",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := resolvedSecretsFile(*secretsFile)
			creds, _ := credentials.Load(path)
			if creds == nil {
				return fmt.Errorf("could not parse %s", path)
			}
			if !validate.Run(creds) {
				os.Exit(1)
			}
			return nil
		},
	}
}

func newKongKeyCmd(secretsFile, configFile, chartsDir, outputDir *string, dryRun *bool) *cobra.Command {
	var pemKey string

	cmd := &cobra.Command{
		Use:   "kong-key",
		Short: "Update Kong's Keycloak public key and redeploy jeeb-infra",
		Long:  "Fetches the RS256 public key from Keycloak's JWKS endpoint, saves it to secrets.yaml, and redeploys jeeb-infra so Kong picks up the change.\n\nRun this once after Keycloak is up. You can also supply the key directly with --key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			secretsPath := resolvedSecretsFile(*secretsFile)
			cfg := config.LoadFromFile(resolvedConfigFile(*configFile))
			return kongkey.NewUpdater(cfg, secretsPath, *chartsDir, *outputDir, *dryRun).
				Run(cmd.Context(), pemKey)
		},
	}

	cmd.Flags().StringVar(&pemKey, "key", "", "PEM public key to use instead of fetching from Keycloak")
	return cmd
}

func newCheckCmd(secretsFile, configFile, outputDir *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "check",
		Short:         "Run a health check on all cluster services",
		SilenceErrors: true,
		SilenceUsage:  true,
		Long: `Checks cluster health and prints a pass/fail table:
  - Pod health        no CrashLoopBackOff / ImagePullBackOff / OOMKilled
  - Keycloak          endpoint reachable, realm verified
  - Vault             endpoint reachable, not sealed
  - Kong              health endpoint reachable
  - Jenkins           login page reachable
  - DNS               .local hostnames resolve inside the cluster
  - Vault secrets     backend secrets readable

Exit code 0 = all healthy, 1 = one or more checks failed.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.LoadFromFile(resolvedConfigFile(*configFile))
			results := healthcheck.RunAll(cmd.Context(), cfg, *outputDir)
			healthcheck.Print(results)
			if !healthcheck.AllPassed(results) {
				return healthcheck.ErrChecksFailed
			}
			return nil
		},
	}
	return cmd
}

func newMaintainCmd(secretsFile, configFile, outputDir *string) *cobra.Command {
	return &cobra.Command{
		Use:   "maintain",
		Short: "Diagnose cluster health and print fix commands (report only, no execution)",
		Long: `Runs the health check and for each failing check prints a diagnosis block
with the exact kubectl/helm commands to resolve it.

No commands are executed — this is a read-only diagnosis report.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.LoadFromFile(resolvedConfigFile(*configFile))
			return maintain.RunReport(cmd.Context(), cfg, *outputDir)
		},
	}
}

func newNamespaceCmd(kubeconfig *string) *cobra.Command {
	return &cobra.Command{
		Use:   "namespace",
		Short: "Create jeeb namespaces (jeeb-dev, jeeb-infra, jeeb-obs)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kube.NewClient(*kubeconfig)
			if err != nil {
				return err
			}
			return client.CreateNamespaces(cmd.Context(), []string{"jeeb-dev", "jeeb-infra", "jeeb-obs"})
		},
	}
}

func newStatusCmd(kubeconfig *string) *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show pod status for jeeb namespaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kube.NewClient(*kubeconfig)
			if err != nil {
				return err
			}
			return client.PrintStatus(cmd.Context(), namespace)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "filter by namespace (default: all jeeb namespaces)")
	return cmd
}

func newRestartCmd(kubeconfig *string) *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "restart <deployment>",
		Short: "Restart a deployment by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kube.NewClient(*kubeconfig)
			if err != nil {
				return err
			}
			return client.RestartDeployment(cmd.Context(), namespace, args[0])
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "jeeb-dev", "target namespace")
	return cmd
}

func newLogsCmd(kubeconfig *string) *cobra.Command {
	var namespace string
	var follow bool
	var tail int64

	cmd := &cobra.Command{
		Use:   "logs <deployment>",
		Short: "Stream logs from a deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kube.NewClient(*kubeconfig)
			if err != nil {
				return err
			}
			return client.StreamLogs(cmd.Context(), namespace, args[0], follow, tail)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "jeeb-dev", "target namespace")
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "stream logs continuously")
	cmd.Flags().Int64Var(&tail, "tail", 100, "number of recent lines to show")
	return cmd
}

func newPatchJenkinsCredsCmd(secretsFile *string, dryRun *bool) *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "patch-jenkins-creds",
		Short: "Patch jenkins-secret from secrets.yaml and restart Jenkins",
		Long: `Updates the jenkins-secret Kubernetes secret with the current values from
secrets.yaml, then rolls out a Jenkins restart so the new credentials take effect.

All six credential keys are patched: admin-password, github-user, github-pat,
nexus-user, nexus-password, sonar-token.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			secretsPath := resolvedSecretsFile(*secretsFile)
			creds, err := loadCreds(secretsPath)
			if err != nil {
				return err
			}
			return jenkins.NewCredentialsPatcher(creds, namespace, *dryRun).Run(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "jeeb-infra", "namespace where Jenkins is deployed")
	return cmd
}

// ── helpers ───────────────────────────────────────────────────────────────────

func loadCreds(secretsPath string) (*credentials.Credentials, error) {
	creds, err := credentials.Load(secretsPath)
	if err != nil {
		return nil, fmt.Errorf("load credentials from %s: %w\n\nCopy env/secrets.yaml.example → env/secrets.yaml and fill all values", secretsPath, err)
	}
	return creds, nil
}

func resolvedSecretsFile(flag string) string {
	if flag != "" {
		return flag
	}
	return defaultSecretsPath()
}

func resolvedConfigFile(flag string) string {
	if flag != "" {
		return flag
	}
	if wd, err := os.Getwd(); err == nil {
		candidate := filepath.Join(wd, "env", "config.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func defaultSecretsPath() string {
	if wd, err := os.Getwd(); err == nil {
		if p := filepath.Join(wd, "env", "secrets.yaml"); fileExists(p) {
			return p
		}
		// backward compat: fall back to credentials.yaml with a deprecation warning
		if p := filepath.Join(wd, "env", "credentials.yaml"); fileExists(p) {
			logger.Warn("env/credentials.yaml is deprecated; rename to env/secrets.yaml")
			return p
		}
	}
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "env", "secrets.yaml")
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func defaultChartsDir() string {
	if wd, err := os.Getwd(); err == nil {
		dir := wd
		for i := 0; i < 5; i++ {
			candidate := filepath.Join(dir, "k8s", "charts")
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
			dir = filepath.Dir(dir)
		}
	}
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(dir, "k8s", "charts")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		dir = filepath.Dir(dir)
	}
	return filepath.Join("k8s", "charts")
}

func defaultGroovyPath() string {
	if wd, err := os.Getwd(); err == nil {
		dir := wd
		for i := 0; i < 5; i++ {
			candidate := filepath.Join(dir, "jenkins", "jobs", "seed.groovy")
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
			dir = filepath.Dir(dir)
		}
	}
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(dir, "jenkins", "jobs", "seed.groovy")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		dir = filepath.Dir(dir)
	}
	return filepath.Join("jenkins", "jobs", "seed.groovy")
}

func defaultOutputDir() string {
	_ = runtime.GOOS
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	exe, _ := os.Executable()
	return filepath.Dir(exe)
}
