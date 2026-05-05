package setup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"k8s-manager/internal/config"
	"k8s-manager/internal/credentials"
	"k8s-manager/internal/helm"
	"k8s-manager/internal/jenkins"
)

type Runner struct {
	cfg         *config.ClusterConfig
	creds       *credentials.Credentials
	chartsDir   string
	outputDir   string
	dryRun      bool
	secretsFile string
}

func NewRunner(cfg *config.ClusterConfig, creds *credentials.Credentials, chartsDir, outputDir string, dryRun bool) *Runner {
	return &Runner{
		cfg:       cfg,
		creds:     creds,
		chartsDir: chartsDir,
		outputDir: outputDir,
		dryRun:    dryRun,
	}
}

func (r *Runner) writeSecretsFile() error {
	if r.dryRun {
		fmt.Printf("      [dry-run] would write values-secrets.yaml to %s\n", r.outputDir)
		r.secretsFile = filepath.Join(r.outputDir, "values-secrets.yaml")
		return nil
	}
	path, err := WriteSecretsValuesFile(r.outputDir, r.creds, r.cfg.NexusRegistry)
	if err != nil {
		return err
	}
	r.secretsFile = path
	fmt.Printf("      wrote %s\n", path)
	return nil
}

// Deploy runs helm upgrade --install for the given targets.
// If targets is empty all charts are deployed.
func (r *Runner) Deploy(ctx context.Context, targets []string) error {
	all := len(targets) == 0
	want := make(map[string]bool, len(targets))
	for _, t := range targets {
		want[t] = true
	}

	fmt.Println("=== Deploy Charts ===")
	if r.dryRun {
		fmt.Println("DRY RUN — commands will be printed, not executed")
	}
	fmt.Println()

	fmt.Println("[0] Generate values-secrets.yaml from credentials")
	if err := r.writeSecretsFile(); err != nil {
		return fmt.Errorf("generate secrets values file: %w", err)
	}
	fmt.Println("    done")
	fmt.Println()

	type step struct {
		name   string
		target string
		fn     func(context.Context) error
	}
	steps := []step{
		{"Deploy jeeb-infra    (Vault, Jenkins, Nexus, SonarQube, Kong)", "infra", r.deployInfra},
		{"Deploy jeeb-data     (MongoDB, Keycloak)", "data", r.deployData},
		{"Deploy jeeb-app      (backend, frontend)", "app", r.deployApp},
		{"Deploy jeeb-learning (learning services)", "learning", r.deployLearning},
		{"Deploy jeeb-obs      (Prometheus, Loki, Grafana)", "obs", r.deployObs},
	}

	n := 0
	for _, s := range steps {
		if !all && !want[s.target] {
			continue
		}
		n++
		fmt.Printf("[%d] %s\n", n, s.name)
		if err := s.fn(ctx); err != nil {
			return fmt.Errorf("deploy %s: %w", s.target, err)
		}
		fmt.Println("    done")
		fmt.Println()
	}

	fmt.Println("=== Done ===")
	r.printAccessTable()
	return nil
}

func (r *Runner) Run(ctx context.Context) error {
	steps := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"Deploy jeeb-infra (Vault, Jenkins, Nexus, SonarQube, Kong)", r.deployInfra},
		{"Deploy jeeb-data (MongoDB, Keycloak)", r.deployData},
		{"Deploy jeeb-app (backend, frontend)", r.deployApp},
		{"Deploy jeeb-learning", r.deployLearning},
		{"Deploy jeeb-obs (Prometheus, Loki, Grafana)", r.deployObs},
		{"Wait for Vault pod ready", r.waitForVault},
		{"Initialize Vault", r.initVault},
		{"Store unseal keys in Kubernetes secret", r.storeUnsealKeysApply},
		{"Unseal Vault", r.unsealVault},
		{"Configure Vault (KV engine, policies, K8s auth roles)", r.configureVault},
		{"Patch CoreDNS for .local DNS", r.patchCoreDNS},
		{"Seed Jenkins (create seed job, generate pipeline jobs)", r.seedJenkins},
	}

	fmt.Println("=== New Cluster Setup ===")
	fmt.Printf("Charts directory : %s\n", r.chartsDir)
	fmt.Printf("Output directory : %s\n", r.outputDir)
	if r.dryRun {
		fmt.Println("DRY RUN — commands will be printed, not executed")
	}
	fmt.Println()

	fmt.Printf("[0/%d] Generate values-secrets.yaml from credentials\n", len(steps))
	if err := r.writeSecretsFile(); err != nil {
		return fmt.Errorf("generate secrets values file: %w", err)
	}
	fmt.Printf("      done\n\n")

	for i, step := range steps {
		fmt.Printf("[%d/%d] %s\n", i+1, len(steps), step.name)
		if err := step.fn(ctx); err != nil {
			return fmt.Errorf("step %d (%s): %w", i+1, step.name, err)
		}
		fmt.Printf("      done\n\n")
	}

	fmt.Println("=== Setup complete ===")
	r.printAccessTable()
	return nil
}

// DeployInfra is exported so kongkey can reuse it after updating credentials.
func (r *Runner) DeployInfra(ctx context.Context) error {
	return r.deployInfra(ctx)
}

// SetSecretsFile allows callers to inject a pre-written secrets file path,
// skipping the writeSecretsFile step (used by kongkey after it writes its own).
func (r *Runner) SetSecretsFile(path string) {
	r.secretsFile = path
}

// ── deploy steps ─────────────────────────────────────────────────────────────

func (r *Runner) deployInfra(ctx context.Context) error {
	return helm.Run(ctx, r.dryRun,
		"upgrade", "--install", r.cfg.ReleaseInfra,
		filepath.Join(r.chartsDir, "jeeb-infra"),
		"--namespace", r.cfg.NamespaceInfra,
		"--create-namespace",
		"-f", r.secretsFile,
	)
}

func (r *Runner) deployData(ctx context.Context) error {
	return helm.Run(ctx, r.dryRun,
		"upgrade", "--install", r.cfg.ReleaseData,
		filepath.Join(r.chartsDir, "jeeb-data"),
		"--namespace", r.cfg.NamespaceDev,
		"--create-namespace",
		"-f", filepath.Join(r.chartsDir, "jeeb-data", "values-dev.yaml"),
		"-f", r.secretsFile,
	)
}

func (r *Runner) deployApp(ctx context.Context) error {
	return helm.Run(ctx, r.dryRun,
		"upgrade", "--install", r.cfg.ReleaseDev,
		filepath.Join(r.chartsDir, "jeeb-app"),
		"--namespace", r.cfg.NamespaceDev,
		"--create-namespace",
		"-f", filepath.Join(r.chartsDir, "jeeb-app", "values-dev.yaml"),
		"-f", r.secretsFile,
	)
}

func (r *Runner) deployLearning(ctx context.Context) error {
	return helm.Run(ctx, r.dryRun,
		"upgrade", "--install", r.cfg.ReleaseLearning,
		filepath.Join(r.chartsDir, "jeeb-learning"),
		"--namespace", r.cfg.NamespaceDev,
		"--create-namespace",
		"-f", filepath.Join(r.chartsDir, "jeeb-learning", "values-dev.yaml"),
		"-f", r.secretsFile,
	)
}

func (r *Runner) deployObs(ctx context.Context) error {
	return helm.Run(ctx, r.dryRun,
		"upgrade", "--install", r.cfg.ReleaseObs,
		filepath.Join(r.chartsDir, "jeeb-obs"),
		"--namespace", r.cfg.NamespaceObs,
		"--create-namespace",
		"-f", r.secretsFile,
	)
}

// ── Vault steps ───────────────────────────────────────────────────────────────

func (r *Runner) waitForVault(ctx context.Context) error {
	return r.kubectl(ctx, "wait",
		fmt.Sprintf("pod/%s", r.cfg.VaultPod),
		"-n", r.cfg.NamespaceInfra,
		"--for=condition=Ready",
		"--timeout=120s",
	)
}

type vaultInitOutput struct {
	UnsealKeysB64 []string `json:"unseal_keys_b64"`
	RootToken     string   `json:"root_token"`
}

func (r *Runner) initVault(ctx context.Context) error {
	outPath := filepath.Join(r.outputDir, "vault-init.json")

	if _, err := os.Stat(outPath); err == nil {
		fmt.Printf("      vault-init.json already exists — skipping init\n")
		return nil
	}

	if r.dryRun {
		fmt.Printf("      kubectl exec -n %s %s -- vault operator init -format=json > %s\n",
			r.cfg.NamespaceInfra, r.cfg.VaultPod, outPath)
		return nil
	}

	out, err := exec.CommandContext(ctx,
		"kubectl", "exec", "-n", r.cfg.NamespaceInfra, r.cfg.VaultPod, "--",
		"vault", "operator", "init", "-format=json",
	).Output()
	if err != nil {
		return fmt.Errorf("vault init: %w", err)
	}

	if err := os.WriteFile(outPath, out, 0600); err != nil {
		return fmt.Errorf("write vault-init.json: %w", err)
	}

	fmt.Printf("      saved to %s — keep this file safe!\n", outPath)
	return nil
}

func (r *Runner) storeUnsealKeysApply(ctx context.Context) error {
	init, err := r.loadVaultInit()
	if err != nil {
		return err
	}

	if len(init.UnsealKeysB64) < 3 {
		return fmt.Errorf("expected at least 3 unseal keys, got %d", len(init.UnsealKeysB64))
	}

	yaml := fmt.Sprintf(`apiVersion: v1
kind: Secret
metadata:
  name: vault-unseal-keys
  namespace: %s
stringData:
  key1: %s
  key2: %s
  key3: %s
`, r.cfg.NamespaceInfra, init.UnsealKeysB64[0], init.UnsealKeysB64[1], init.UnsealKeysB64[2])

	if r.dryRun {
		fmt.Printf("      kubectl apply -f - (vault-unseal-keys secret in %s)\n", r.cfg.NamespaceInfra)
		return nil
	}

	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(yaml)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *Runner) unsealVault(ctx context.Context) error {
	init, err := r.loadVaultInit()
	if err != nil {
		return err
	}

	for i, key := range init.UnsealKeysB64[:3] {
		fmt.Printf("      unsealing with key %d/3...\n", i+1)
		if err := r.vaultExec(ctx, init.RootToken, "operator", "unseal", key); err != nil {
			return fmt.Errorf("unseal key %d: %w", i+1, err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

type vaultKV struct {
	path string
	key  string
	val  string
}

func (r *Runner) configureVault(ctx context.Context) error {
	init, err := r.loadVaultInit()
	if err != nil {
		return err
	}
	token := init.RootToken

	mongoURI := func(db string) string {
		return fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin",
			r.creds.MongoDBUsername, r.creds.MongoDBPassword, r.cfg.MongoHost, db)
	}
	keycloakURL := fmt.Sprintf("http://%s", r.cfg.KeycloakHost)
	keycloakPublic := fmt.Sprintf("http://localhost:%d", r.cfg.KeycloakNodePort)
	backendPublic := fmt.Sprintf("http://localhost:%d", r.cfg.BackendNodePort)
	learningPublic := fmt.Sprintf("http://localhost:%d", r.cfg.LearningNodePort)

	secrets := []vaultKV{
		// backend
		{r.cfg.VaultPathBackend, "PORT", "8080"},
		{r.cfg.VaultPathBackend, "LOG_LEVEL", "INFO"},
		{r.cfg.VaultPathBackend, "MONGO_DATABASE", "jeeb"},
		{r.cfg.VaultPathBackend, "MONGO_URI", mongoURI("jeeb")},
		{r.cfg.VaultPathBackend, "KEYCLOAK_URL", keycloakURL},
		{r.cfg.VaultPathBackend, "KEYCLOAK_REALM", r.cfg.KeycloakRealm},
		{r.cfg.VaultPathBackend, "KEYCLOAK_CLIENT_ID", r.cfg.KeycloakClientID},
		// frontend
		{r.cfg.VaultPathFrontend, "VITE_KEYCLOAK_URL", keycloakPublic},
		{r.cfg.VaultPathFrontend, "VITE_KEYCLOAK_REALM", r.cfg.KeycloakRealm},
		{r.cfg.VaultPathFrontend, "VITE_KEYCLOAK_CLIENT_ID", r.cfg.KeycloakClientID},
		{r.cfg.VaultPathFrontend, "VITE_API_URL", backendPublic},
		// learning backend
		{r.cfg.VaultPathLearningBackend, "PORT", "8080"},
		{r.cfg.VaultPathLearningBackend, "LOG_LEVEL", "INFO"},
		{r.cfg.VaultPathLearningBackend, "MONGO_DATABASE", "jeeb_learning"},
		{r.cfg.VaultPathLearningBackend, "MONGO_URI", mongoURI("jeeb_learning")},
		{r.cfg.VaultPathLearningBackend, "KEYCLOAK_URL", keycloakURL},
		{r.cfg.VaultPathLearningBackend, "KEYCLOAK_REALM", r.cfg.KeycloakRealm},
		{r.cfg.VaultPathLearningBackend, "KEYCLOAK_CLIENT_ID", r.cfg.KeycloakClientID},
		// learning frontend
		{r.cfg.VaultPathLearningFrontend, "VITE_KEYCLOAK_URL", keycloakPublic},
		{r.cfg.VaultPathLearningFrontend, "VITE_KEYCLOAK_REALM", r.cfg.KeycloakRealm},
		{r.cfg.VaultPathLearningFrontend, "VITE_KEYCLOAK_CLIENT_ID", r.cfg.KeycloakClientID},
		{r.cfg.VaultPathLearningFrontend, "VITE_API_URL", learningPublic},
	}

	type serviceRole struct {
		name      string
		sa        string
		policyTpl string
		vaultPath string
	}
	roles := []serviceRole{
		{"backend", "backend", "backend-policy", r.cfg.VaultPathBackend},
		{"frontend", "frontend", "frontend-policy", r.cfg.VaultPathFrontend},
		{"learning-backend", "learning-backend", "learning-backend-policy", r.cfg.VaultPathLearningBackend},
		{"learning-frontend", "learning-frontend", "learning-frontend-policy", r.cfg.VaultPathLearningFrontend},
	}

	steps := []func() error{
		func() error {
			return r.vaultExec(ctx, token, "secrets", "enable", "-path=secret", "kv-v2")
		},
		func() error {
			written := map[string]bool{}
			for _, s := range secrets {
				if written[s.path] {
					continue
				}
				// vault kv put uses CLI path format (secret/X), not API format (secret/data/X)
				args := []string{"kv", "put", config.KVCLIPath(s.path)}
				for _, kv := range secrets {
					if kv.path == s.path {
						args = append(args, kv.key+"="+kv.val)
					}
				}
				if err := r.vaultExec(ctx, token, args...); err != nil {
					return err
				}
				written[s.path] = true
			}
			return nil
		},
		func() error { return r.vaultExec(ctx, token, "auth", "enable", "kubernetes") },
		func() error {
			return r.vaultExec(ctx, token, "write", "auth/kubernetes/config",
				"kubernetes_host=https://kubernetes.default.svc:443")
		},
	}

	for _, role := range roles {
		role := role
		// VaultPath is already in API format (secret/data/...) — use directly in policy
		policy := fmt.Sprintf(`path "%s" { capabilities = ["read"] }`, role.vaultPath)
		steps = append(steps,
			func() error {
				return r.vaultWritePolicy(ctx, token, role.policyTpl, policy)
			},
			func() error {
				return r.vaultExec(ctx, token, "write",
					fmt.Sprintf("auth/kubernetes/role/%s", role.name),
					fmt.Sprintf("bound_service_account_names=%s", role.sa),
					fmt.Sprintf("bound_service_account_namespaces=%s", r.cfg.NamespaceDev),
					fmt.Sprintf("policies=%s", role.policyTpl),
					"ttl=1h",
				)
			},
		)
	}

	for _, s := range steps {
		if err := s(); err != nil {
			if !strings.Contains(err.Error(), "already enabled") &&
				!strings.Contains(err.Error(), "path is already in use") {
				return err
			}
		}
	}
	return nil
}

// ── CoreDNS ───────────────────────────────────────────────────────────────────

func (r *Runner) patchCoreDNS(ctx context.Context) error {
	coreDNSPatch := filepath.Join(filepath.Dir(r.chartsDir), "coredns-patch.yaml")

	if r.dryRun {
		fmt.Printf("      kubectl get svc -n %s -l %s -o jsonpath='{.items[0].spec.clusterIP}'\n",
			r.cfg.NamespaceDev, r.cfg.IngressLabel)
		fmt.Printf("      [substitute IP in %s and pipe to] kubectl apply -f -\n", coreDNSPatch)
		return nil
	}

	out, err := exec.CommandContext(ctx,
		"kubectl", "get", "svc",
		"-n", r.cfg.NamespaceDev,
		"-l", r.cfg.IngressLabel,
		"-o", "jsonpath={.items[0].spec.clusterIP}",
	).Output()

	ip := strings.TrimSpace(string(out))
	if err != nil || ip == "" {
		fmt.Printf("      WARNING: could not auto-detect ingress ClusterIP (%v)\n", err)
		fmt.Println("      Update the IP in k8s/coredns-patch.yaml manually, then:")
		fmt.Println("      kubectl apply -f k8s/coredns-patch.yaml")
		return nil
	}

	fmt.Printf("      ingress ClusterIP: %s\n", ip)

	patchYAML, err := os.ReadFile(coreDNSPatch)
	if err != nil {
		fmt.Printf("      WARNING: could not read %s: %v\n", coreDNSPatch, err)
		fmt.Println("      Apply manually: kubectl apply -f k8s/coredns-patch.yaml")
		return nil
	}

	updated := substituteCoreDNSIP(string(patchYAML), ip)

	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Stdin = bytes.NewReader([]byte(updated))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("kubectl apply coredns patch: %w", err)
	}
	return nil
}

func substituteCoreDNSIP(yaml, newIP string) string {
	lines := strings.Split(yaml, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 && isIPLike(trimmed) {
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
				lines[i] = indent + newIP + " " + strings.Join(parts[1:], " ")
			}
		}
	}
	return strings.Join(lines, "\n")
}

func isIPLike(s string) bool {
	parts := strings.SplitN(s, " ", 2)
	if len(parts) == 0 {
		return false
	}
	octets := strings.Split(parts[0], ".")
	return len(octets) == 4
}

// ── helpers ──────────────────────────────────────────────────────────────────

func (r *Runner) kubectl(ctx context.Context, args ...string) error {
	if r.dryRun {
		fmt.Printf("      kubectl %s\n", strings.Join(args, " "))
		return nil
	}
	cmd := exec.CommandContext(ctx, "kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *Runner) vaultExec(ctx context.Context, token string, args ...string) error {
	base := []string{
		"exec", "-i", "-n", r.cfg.NamespaceInfra, r.cfg.VaultPod, "--",
		"env",
		fmt.Sprintf("VAULT_ADDR=%s", r.cfg.VaultAddr),
		"VAULT_TOKEN=" + token,
		"vault",
	}
	return r.kubectl(ctx, append(base, args...)...)
}

func (r *Runner) vaultWritePolicy(ctx context.Context, token, name, policy string) error {
	if r.dryRun {
		fmt.Printf("      vault policy write %s -\n", name)
		return nil
	}
	cmd := exec.CommandContext(ctx,
		"kubectl", "exec", "-i", "-n", r.cfg.NamespaceInfra, r.cfg.VaultPod, "--",
		"env",
		fmt.Sprintf("VAULT_ADDR=%s", r.cfg.VaultAddr),
		"VAULT_TOKEN="+token,
		"vault", "policy", "write", name, "-",
	)
	cmd.Stdin = strings.NewReader(policy)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *Runner) loadVaultInit() (*vaultInitOutput, error) {
	outPath := filepath.Join(r.outputDir, "vault-init.json")
	data, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("read vault-init.json (run step 5 first): %w", err)
	}
	var v vaultInitOutput
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("parse vault-init.json: %w", err)
	}
	return &v, nil
}

func (r *Runner) seedJenkins(ctx context.Context) error {
	// Derive repo root from chartsDir (k8s/charts → k8s → repo root).
	repoRoot := filepath.Dir(filepath.Dir(r.chartsDir))
	groovyPath := filepath.Join(repoRoot, "jenkins", "jobs", "seed.groovy")

	if _, err := os.Stat(groovyPath); err != nil {
		fmt.Printf("      WARNING: seed.groovy not found at %s\n", groovyPath)
		fmt.Println("      Run 'k8s-manager seed --groovy-path <path>' after Jenkins is ready")
		return nil
	}

	return jenkins.NewSeeder(r.cfg, r.creds, groovyPath, "", r.dryRun).Run(ctx)
}

func (r *Runner) printAccessTable() {
	fmt.Printf(`
  %s (infra):
    Jenkins    http://localhost:%d
    Nexus      http://localhost:%d
    SonarQube  http://localhost:%d
    Kong       http://localhost:%d
    Vault      http://localhost:%d

  %s (dev):
    Frontend          http://localhost:%d
    Backend           http://localhost:%d
    Keycloak          http://localhost:%d
    MongoDB           localhost:%d
    Learning backend  http://localhost:%d
    Learning frontend http://localhost:%d

  %s (obs):
    Grafana    http://localhost:%d
    Prometheus http://localhost:%d
`,
		r.cfg.NamespaceInfra,
		r.cfg.JenkinsNodePort,
		r.cfg.NexusUINodePort,
		r.cfg.SonarQubeNodePort,
		r.cfg.KongNodePort,
		r.cfg.VaultNodePort,
		r.cfg.NamespaceDev,
		r.cfg.FrontendNodePort,
		r.cfg.BackendNodePort,
		r.cfg.KeycloakNodePort,
		r.cfg.MongoNodePort,
		r.cfg.LearningNodePort,
		r.cfg.LearningFrontPort,
		r.cfg.NamespaceObs,
		r.cfg.GrafanaNodePort,
		r.cfg.PrometheusNodePort,
	)
}
