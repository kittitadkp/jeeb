package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"k8s-manager/internal/credentials"
	"k8s-manager/internal/kongkey"
	"k8s-manager/internal/kube"
	"k8s-manager/internal/rancher"
	"k8s-manager/internal/setup"
	"k8s-manager/internal/validate"

	"github.com/spf13/cobra"
)

func main() {
	var kubeconfig string

	root := &cobra.Command{
		Use:   "k8s-manager",
		Short: "CLI for managing the jeeb Kubernetes cluster",
	}

	root.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig (defaults to ~/.kube/config)")

	root.AddCommand(
		newStatusCmd(&kubeconfig),
		newRestartCmd(&kubeconfig),
		newLogsCmd(&kubeconfig),
		newSetupCmd(),
		newDeployCmd(),
		newValidateCmd(),
		newKongKeyCmd(),
		newRancherCmd(),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newSetupCmd() *cobra.Command {
	var credsFile string
	var chartsDir string
	var outputDir string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Bootstrap a new cluster end-to-end (deploy charts, init Vault, patch CoreDNS)",
		Long: `Runs the full new-cluster setup in order:
  1. Deploy jeeb-infra  (Vault, Jenkins, Nexus, SonarQube, Kong)
  2. Deploy jeeb-app    (MongoDB, Keycloak, backend, frontend, learning)
  3. Deploy jeeb-obs    (Prometheus, Loki, Grafana)
  4. Wait for Vault pod ready
  5. Initialize Vault   → saves vault-init.json to --output-dir
  6. Store unseal keys  → Kubernetes secret vault-unseal-keys
  7. Unseal Vault
  8. Configure Vault    (KV engine, policies, K8s auth roles)
  9. Patch CoreDNS      (detects ingress ClusterIP automatically)

Copy env/credentials.yaml.example to env/credentials.yaml and fill all values before running.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if credsFile == "" {
				credsFile = defaultCredsPath()
			}

			creds, err := credentials.Load(credsFile)
			if err != nil {
				return fmt.Errorf("load credentials from %s: %w\n\nCopy env/credentials.yaml.example → env/credentials.yaml and fill all values", credsFile, err)
			}

			runner := setup.NewRunner(creds, chartsDir, outputDir, dryRun)
			return runner.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&credsFile, "credentials", "", "path to credentials.yaml (default: env/credentials.yaml relative to binary)")
	cmd.Flags().StringVar(&chartsDir, "charts-dir", defaultChartsDir(), "path to k8s/charts directory")
	cmd.Flags().StringVar(&outputDir, "output-dir", defaultOutputDir(), "directory to write vault-init.json")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print commands without executing them")

	return cmd
}

// defaultCredsPath resolves env/credentials.yaml relative to the binary location.
func defaultCredsPath() string {
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "env", "credentials.yaml")
}

func defaultChartsDir() string {
	exe, _ := os.Executable()
	// walk up to find k8s/charts (works when running from repo root too)
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

func defaultOutputDir() string {
	exe, _ := os.Executable()
	_ = runtime.GOOS
	return filepath.Dir(exe)
}

func newDeployCmd() *cobra.Command {
	var credsFile string
	var chartsDir string
	var outputDir string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "deploy [infra] [app] [learning] [obs]",
		Short: "Re-deploy one or more Helm charts (no Vault init)",
		Long: `Runs helm upgrade --install for the selected charts.
Useful after changing values or image tags without re-running full setup.

Targets: infra, app, learning, obs (default: all four)

Examples:
  k8s-manager deploy            # deploy infra + app + learning + obs
  k8s-manager deploy app        # deploy only jeeb-app
  k8s-manager deploy learning   # deploy only jeeb-learning
  k8s-manager deploy infra obs  # deploy jeeb-infra and jeeb-obs`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if credsFile == "" {
				credsFile = defaultCredsPath()
			}
			for _, t := range args {
				switch t {
				case "infra", "app", "learning", "obs":
				default:
					return fmt.Errorf("unknown target %q — valid targets: infra, app, learning, obs", t)
				}
			}
			creds, err := credentials.Load(credsFile)
			if err != nil {
				return fmt.Errorf("load credentials: %w", err)
			}
			runner := setup.NewRunner(creds, chartsDir, outputDir, dryRun)
			return runner.Deploy(cmd.Context(), args)
		},
	}

	cmd.Flags().StringVar(&credsFile, "credentials", "", "path to credentials.yaml")
	cmd.Flags().StringVar(&chartsDir, "charts-dir", defaultChartsDir(), "path to k8s/charts directory")
	cmd.Flags().StringVar(&outputDir, "output-dir", defaultOutputDir(), "directory for values-secrets.yaml")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print commands without executing them")
	return cmd
}

func newRancherCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "rancher",
		Short: "Install cert-manager and Rancher (one-time, optional)",
		Long: `Installs cert-manager and Rancher via external Helm repos.
Run this once if you want the Rancher UI in addition to the main cluster.
Rancher is NOT required for jeeb services to run.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			d := rancher.NewDeployer(dryRun)
			return d.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print commands without executing them")
	return cmd
}

func newValidateCmd() *cobra.Command {
	var credsFile string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Check that all credentials in credentials.yaml are filled",
		RunE: func(cmd *cobra.Command, args []string) error {
			if credsFile == "" {
				credsFile = defaultCredsPath()
			}
			// Load always returns creds even when validation fails, so we can
			// show the full table instead of just an error message.
			creds, _ := credentials.Load(credsFile)
			if creds == nil {
				return fmt.Errorf("could not parse %s", credsFile)
			}
			ok := validate.Run(creds)
			if !ok {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&credsFile, "credentials", "", "path to credentials.yaml")
	return cmd
}

func newKongKeyCmd() *cobra.Command {
	var credsFile string
	var chartsDir string
	var outputDir string
	var dryRun bool
	var pemKey string

	cmd := &cobra.Command{
		Use:   "kong-key",
		Short: "Update Kong's Keycloak public key and redeploy jeeb-infra",
		Long: `Fetches the RS256 public key from Keycloak's JWKS endpoint (http://localhost:30081),
saves it to credentials.yaml, and redeploys jeeb-infra so Kong picks up the change.

Run this once after Keycloak is up. You can also supply the key directly with --key.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if credsFile == "" {
				credsFile = defaultCredsPath()
			}
			u := kongkey.NewUpdater(credsFile, chartsDir, outputDir, dryRun)
			return u.Run(cmd.Context(), pemKey)
		},
	}

	cmd.Flags().StringVar(&credsFile, "credentials", "", "path to credentials.yaml")
	cmd.Flags().StringVar(&chartsDir, "charts-dir", defaultChartsDir(), "path to k8s/charts directory")
	cmd.Flags().StringVar(&outputDir, "output-dir", defaultOutputDir(), "directory for values-secrets.yaml")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print commands without executing them")
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
