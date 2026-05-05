package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"k8s-manager/internal/config"
	"k8s-manager/internal/credentials"
	"k8s-manager/internal/jenkins"
	"k8s-manager/internal/kongkey"
	"k8s-manager/internal/kube"
	"k8s-manager/internal/rancher"
	"k8s-manager/internal/setup"
	"k8s-manager/internal/validate"

	"github.com/spf13/cobra"
)

func main() {
	var (
		kubeconfig string
		credsFile  string
		chartsDir  string
		outputDir  string
		dryRun     bool
	)

	root := &cobra.Command{
		Use:   "k8s-manager",
		Short: "CLI for managing the jeeb Kubernetes cluster",
	}

	root.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig (defaults to ~/.kube/config)")
	root.PersistentFlags().StringVar(&credsFile, "credentials", "", "path to credentials.yaml (default: env/credentials.yaml relative to binary)")
	root.PersistentFlags().StringVar(&chartsDir, "charts-dir", defaultChartsDir(), "path to k8s/charts directory")
	root.PersistentFlags().StringVar(&outputDir, "output-dir", defaultOutputDir(), "directory to write vault-init.json and values-secrets.yaml")
	root.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "print commands without executing them")

	root.AddCommand(
		newStatusCmd(&kubeconfig),
		newRestartCmd(&kubeconfig),
		newLogsCmd(&kubeconfig),
		newSetupCmd(&credsFile, &chartsDir, &outputDir, &dryRun),
		newDeployCmd(&credsFile, &chartsDir, &outputDir, &dryRun),
		newSeedCmd(&credsFile, &dryRun),
		newRancherCmd(&credsFile, &dryRun),
		newValidateCmd(&credsFile),
		newKongKeyCmd(&credsFile, &chartsDir, &outputDir, &dryRun),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newSetupCmd(credsFile, chartsDir, outputDir *string, dryRun *bool) *cobra.Command {
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

Copy env/credentials.yaml.example to env/credentials.yaml and fill all values before running.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := resolvedCredsFile(*credsFile)
			creds, err := loadCreds(path)
			if err != nil {
				return err
			}
			cfg := config.LoadFromFile(path)
			return setup.NewRunner(cfg, creds, *chartsDir, *outputDir, *dryRun).Run(cmd.Context())
		},
	}
}

func newDeployCmd(credsFile, chartsDir, outputDir *string, dryRun *bool) *cobra.Command {
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
			path := resolvedCredsFile(*credsFile)
			creds, err := loadCreds(path)
			if err != nil {
				return err
			}
			cfg := config.LoadFromFile(path)
			return setup.NewRunner(cfg, creds, *chartsDir, *outputDir, *dryRun).Deploy(cmd.Context(), args)
		},
	}
}

func newSeedCmd(credsFile *string, dryRun *bool) *cobra.Command {
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
			path := resolvedCredsFile(*credsFile)
			creds, err := loadCreds(path)
			if err != nil {
				return err
			}
			cfg := config.LoadFromFile(path)
			return jenkins.NewSeeder(cfg, creds, groovyPath, jenkinsURL, *dryRun).Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&groovyPath, "groovy-path", defaultGroovyPath(), "path to seed.groovy")
	cmd.Flags().StringVar(&jenkinsURL, "jenkins-url", "", "Jenkins URL (default: http://localhost:<nodeport>)")
	return cmd
}

func newRancherCmd(credsFile *string, dryRun *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "rancher",
		Short: "Install cert-manager and Rancher (one-time, optional)",
		Long: `Installs cert-manager and Rancher via external Helm repos.
Run this once if you want the Rancher UI in addition to the main cluster.
Rancher is NOT required for jeeb services to run.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.LoadFromFile(resolvedCredsFile(*credsFile))
			return rancher.NewDeployer(cfg, *dryRun).Run(cmd.Context())
		},
	}
}

func newValidateCmd(credsFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Check that all credentials in credentials.yaml are filled",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := resolvedCredsFile(*credsFile)
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

func newKongKeyCmd(credsFile, chartsDir, outputDir *string, dryRun *bool) *cobra.Command {
	var pemKey string

	cmd := &cobra.Command{
		Use:   "kong-key",
		Short: "Update Kong's Keycloak public key and redeploy jeeb-infra",
		Long:  "Fetches the RS256 public key from Keycloak's JWKS endpoint, saves it to credentials.yaml, and redeploys jeeb-infra so Kong picks up the change.\n\nRun this once after Keycloak is up. You can also supply the key directly with --key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := resolvedCredsFile(*credsFile)
			cfg := config.LoadFromFile(path)
			return kongkey.NewUpdater(cfg, path, *chartsDir, *outputDir, *dryRun).
				Run(cmd.Context(), pemKey)
		},
	}

	cmd.Flags().StringVar(&pemKey, "key", "", "PEM public key to use instead of fetching from Keycloak")
	return cmd
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

// ── helpers ───────────────────────────────────────────────────────────────────

func loadCreds(credsFile string) (*credentials.Credentials, error) {
	path := resolvedCredsFile(credsFile)
	creds, err := credentials.Load(path)
	if err != nil {
		return nil, fmt.Errorf("load credentials from %s: %w\n\nCopy env/credentials.yaml.example → env/credentials.yaml and fill all values", path, err)
	}
	return creds, nil
}

func resolvedCredsFile(flag string) string {
	if flag != "" {
		return flag
	}
	return defaultCredsPath()
}

func defaultCredsPath() string {
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "env", "credentials.yaml")
}

func defaultChartsDir() string {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	for i := 0; i < 4; i++ {
		candidate := filepath.Join(dir, "k8s", "charts")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		dir = filepath.Dir(dir)
	}
	return filepath.Join("k8s", "charts")
}

func defaultGroovyPath() string {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	for i := 0; i < 4; i++ {
		candidate := filepath.Join(dir, "jenkins", "jobs", "seed.groovy")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		dir = filepath.Dir(dir)
	}
	return filepath.Join("jenkins", "jobs", "seed.groovy")
}

func defaultOutputDir() string {
	exe, _ := os.Executable()
	_ = runtime.GOOS
	return filepath.Dir(exe)
}
