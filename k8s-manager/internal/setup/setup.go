package setup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"k8s-manager/internal/config"
	"k8s-manager/internal/credentials"
	"k8s-manager/internal/helm"
	"k8s-manager/internal/jenkins"
	"k8s-manager/internal/logger"
	"k8s-manager/internal/rancher"
	"k8s-manager/internal/util"

	"gopkg.in/yaml.v3"
)

type Runner struct {
	cfg         *config.ClusterConfig
	creds       *credentials.Credentials
	chartsDir   string
	outputDir   string
	secretsPath string
	dryRun      bool
	secretsFile string
}

func NewRunner(cfg *config.ClusterConfig, creds *credentials.Credentials, chartsDir, outputDir, secretsPath string, dryRun bool) *Runner {
	return &Runner{
		cfg:         cfg,
		creds:       creds,
		chartsDir:   chartsDir,
		outputDir:   outputDir,
		secretsPath: secretsPath,
		dryRun:      dryRun,
	}
}

func (r *Runner) writeSecretsFile() error {
	if r.dryRun {
		logger.Step("      [dry-run] would write values-secrets.yaml to %s", r.outputDir)
		r.secretsFile = filepath.Join(r.outputDir, "values-secrets.yaml")
		return nil
	}
	path, err := WriteSecretsValuesFile(r.outputDir, r.creds, r.cfg.NexusRegistry)
	if err != nil {
		return err
	}
	r.secretsFile = path
	logger.Debug("wrote %s", path)
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

	logger.StepMsg("=== Deploy Charts ===")
	if r.dryRun {
		logger.StepMsg("DRY RUN — commands will be printed, not executed")
	}
	logger.StepMsg("")

	logger.StepMsg("[0] Generate values-secrets.yaml from credentials")
	if err := r.writeSecretsFile(); err != nil {
		return fmt.Errorf("generate secrets values file: %w", err)
	}
	logger.StepMsg("    done")
	logger.StepMsg("")

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
		logger.Step("[%d] %s", n, s.name)
		if err := s.fn(ctx); err != nil {
			return fmt.Errorf("deploy %s: %w", s.target, err)
		}
		logger.StepMsg("    done")
		logger.StepMsg("")
	}

	logger.StepMsg("=== Done ===")
	r.printAccessTable()
	return nil
}

func (r *Runner) Run(ctx context.Context) error {
	type step struct {
		name string
		fn   func(context.Context) error
	}

	steps := []step{
		{"Pre-flight checks", r.preflight},
		{"Generate values-secrets.yaml from credentials", r.writeSecretsFileStep},
		{"Install nginx ingress controller", r.ensureNginxIngress},
		{"Install Rancher + cert-manager", r.ensureRancher},
		{"Deploy jeeb-infra (Vault, Jenkins, Nexus, SonarQube, Kong)", r.deployInfra},
		{"Deploy jeeb-data (MongoDB, Keycloak)", r.deployData},
		{"Deploy jeeb-obs (Prometheus, Loki, Grafana)", r.deployObs},
		{"Initialize Nexus Docker registry", r.initNexusDockerRepo},
		{"Wait for Keycloak ready", r.waitForKeycloak},
		{"Fetch Kong RS256 key from Keycloak", r.fetchAndApplyKongKey},
		{"Re-deploy jeeb-infra with Kong key", r.deployInfra},
		{"Wait for Kong ready", r.waitForKong},
		{"Wait for Vault pod ready", r.waitForVault},
		{"Initialize Vault", r.initVault},
		{"Store unseal keys in Kubernetes secret", r.storeUnsealKeysApply},
		{"Unseal Vault", r.unsealVault},
		{"Configure Vault (KV engine, policies, K8s auth roles)", r.configureVault},
		{"Patch CoreDNS for .local DNS", r.patchCoreDNS},
		{"Wait for CoreDNS rollout", r.waitForCoreDNS},
		{"Verify DNS for all .local domains", r.verifyAllDNS},
		{"Seed Jenkins (create seed job, generate pipeline jobs)", r.seedJenkins},
	}

	logger.StepMsg("=== New Cluster Setup ===")
	logger.Step("Charts directory : %s", r.chartsDir)
	logger.Step("Output directory : %s", r.outputDir)
	if r.dryRun {
		logger.StepMsg("DRY RUN — commands will be printed, not executed")
	}
	logger.StepMsg("")

	for i, step := range steps {
		logger.Step("[%d/%d] %s", i+1, len(steps), step.name)
		if err := step.fn(ctx); err != nil {
			return fmt.Errorf("step %d (%s): %w", i+1, step.name, err)
		}
		logger.StepMsg("      done\n")
	}

	logger.StepMsg("=== Setup complete ===")
	logger.StepMsg("")
	logger.StepMsg("Next steps (manual):")
	logger.StepMsg("  1. Run Jenkins pipelines: backend, frontend, learning-backend, learning-frontend")
	logger.StepMsg("  2. Pipelines build + push images to Nexus (localhost:30050)")
	logger.StepMsg("  3. Deploy app: k8s-manager deploy app learning")
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

// ── pre-flight ────────────────────────────────────────────────────────────────

func (r *Runner) preflight(ctx context.Context) error {
	if r.dryRun {
		logger.Step("      [dry-run] skipping pre-flight checks")
		return nil
	}

	// Verify cluster is reachable
	out, err := util.RunCmdOutput(ctx, "kubectl", "cluster-info")
	if err != nil {
		return fmt.Errorf("cluster not reachable (is Docker Desktop running?): %w", err)
	}
	logger.Step("      cluster: %s", strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)[0])

	// Warn if on unexpected context
	ctxOut, _ := util.RunCmdOutput(ctx, "kubectl", "config", "current-context")
	currentCtx := strings.TrimSpace(string(ctxOut))
	if currentCtx != "docker-desktop" {
		logger.Warn("current context is %q, expected docker-desktop", currentCtx)
	} else {
		logger.Step("      context: %s", currentCtx)
	}

	// Warn if stale vault-init.json exists
	vaultInitPath := filepath.Join(r.outputDir, "vault-init.json")
	if _, err := os.Stat(vaultInitPath); err == nil {
		logger.Warn("%s exists from a previous run. Vault init will be skipped. Delete it if this is a truly fresh cluster.", vaultInitPath)
	}

	// Verify helm is on PATH
	if _, err := exec.LookPath("helm"); err != nil {
		return fmt.Errorf("helm not found on PATH: %w", err)
	}

	return nil
}

func (r *Runner) writeSecretsFileStep(ctx context.Context) error {
	return r.writeSecretsFile()
}

// ── nginx ingress ─────────────────────────────────────────────────────────────

func (r *Runner) ensureNginxIngress(ctx context.Context) error {
	// Check if already installed
	out, _ := util.RunCmdOutput(ctx, "kubectl", "get", "ns", "ingress-nginx", "--no-headers")
	if strings.Contains(string(out), "ingress-nginx") {
		logger.Step("      nginx ingress already installed — skipping")
		return nil
	}

	const manifest = "https://raw.githubusercontent.com/kubernetes/ingress-nginx/" +
		"controller-v1.12.2/deploy/static/provider/cloud/deploy.yaml"

	logger.Step("      applying %s", manifest)
	if err := r.kubectl(ctx, "apply", "-f", manifest); err != nil {
		return fmt.Errorf("install nginx ingress: %w", err)
	}

	logger.Step("      waiting for ingress controller to become available (up to 2 min)...")
	return r.kubectl(ctx, "wait",
		"-n", "ingress-nginx",
		"deployment/ingress-nginx-controller",
		"--for=condition=Available",
		"--timeout=120s",
	)
}

// ── Rancher ───────────────────────────────────────────────────────────────────

func (r *Runner) ensureRancher(ctx context.Context) error {
	// Check if already installed
	out, _ := util.RunCmdOutput(ctx, "kubectl", "get", "ns", r.cfg.RancherNamespace, "--no-headers")
	if strings.Contains(string(out), r.cfg.RancherNamespace) {
		logger.Step("      Rancher namespace %s already exists — skipping", r.cfg.RancherNamespace)
		return nil
	}
	return rancher.NewDeployer(r.cfg, r.dryRun).Run(ctx)
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

// ── Nexus ─────────────────────────────────────────────────────────────────────

func (r *Runner) initNexusDockerRepo(ctx context.Context) error {
	if r.dryRun {
		logger.Step("      [dry-run] would initialize Nexus Docker hosted repo")
		return nil
	}

	nexusURL := fmt.Sprintf("http://localhost:%d", r.cfg.NexusUINodePort)
	statusURL := nexusURL + "/service/rest/v1/status"

	logger.Debug("waiting for Nexus at %s (up to 5 min)...", nexusURL)
	if err := util.PollHTTP(ctx, util.PollConfig{Timeout: 5 * time.Minute}, nil, http.MethodGet, statusURL, http.StatusOK, nil); err != nil {
		logger.Warn("Nexus did not become ready within 5 min, continuing anyway")
	} else {
		logger.Step("      Nexus is up")
	}

	// Read initial admin password from pod
	podOut, err := util.RunCmdOutput(ctx,
		"kubectl", "get", "pod",
		"-n", r.cfg.NamespaceInfra,
		"-l", "app=nexus",
		"-o", "jsonpath={.items[0].metadata.name}",
	)
	if err != nil || strings.TrimSpace(string(podOut)) == "" {
		logger.Warn("could not find Nexus pod — skipping Docker repo init")
		return nil
	}
	nexusPod := strings.TrimSpace(string(podOut))

	initPassOut, err := util.RunCmdOutput(ctx,
		"kubectl", "exec", "-n", r.cfg.NamespaceInfra, nexusPod, "--",
		"cat", "/nexus-data/admin.password",
	)
	if err != nil {
		logger.Step("      NOTE: /nexus-data/admin.password not found — Nexus may already be configured")
		return nil
	}
	initPass := strings.TrimSpace(string(initPassOut))

	// Change admin password
	newPass := r.creds.NexusAdminPassword
	req, err := http.NewRequestWithContext(ctx, http.MethodPut,
		nexusURL+"/service/rest/v1/security/users/admin/change-password",
		strings.NewReader(newPass),
	)
	if err != nil {
		logger.Warn("could not build Nexus password request: %v", err)
	} else {
		req.SetBasicAuth("admin", initPass)
		req.Header.Set("Content-Type", "text/plain")
		resp, doErr := http.DefaultClient.Do(req)
		if doErr != nil {
			logger.Warn("could not change Nexus admin password: %v", doErr)
		} else {
			resp.Body.Close()
			logger.Step("      Nexus admin password updated")
		}
	}

	// Create Docker hosted repo (idempotent: ignore 400 if already exists)
	repoBody := `{
		"name": "jeeb",
		"online": true,
		"storage": {"blobStoreName": "default", "strictContentTypeValidation": true, "writePolicy": "ALLOW"},
		"docker": {"v1Enabled": false, "forceBasicAuth": false, "httpPort": 5000}
	}`
	req2, err := http.NewRequestWithContext(ctx, http.MethodPost,
		nexusURL+"/service/rest/v1/repositories/docker/hosted",
		strings.NewReader(repoBody),
	)
	if err != nil {
		logger.Warn("could not build Nexus repo request: %v", err)
	} else {
		req2.SetBasicAuth("admin", newPass)
		req2.Header.Set("Content-Type", "application/json")
		resp2, doErr := http.DefaultClient.Do(req2)
		if doErr != nil {
			logger.Warn("could not create Nexus Docker repo: %v", doErr)
		} else {
			resp2.Body.Close()
			if resp2.StatusCode == 201 {
				logger.Step("      created Docker hosted repo 'jeeb' on port 5000")
			} else if resp2.StatusCode == 400 {
				logger.Step("      Docker hosted repo already exists — skipping")
			}
		}
	}

	return nil
}

// ── Keycloak wait ─────────────────────────────────────────────────────────────

func (r *Runner) waitForKeycloak(ctx context.Context) error {
	url := fmt.Sprintf("http://localhost:%d/realms/%s", r.cfg.KeycloakNodePort, r.cfg.KeycloakRealm)
	logger.Debug("polling %s (up to 5 min)...", url)
	if err := util.PollHTTP(ctx, util.PollConfig{Timeout: 5 * time.Minute}, nil, http.MethodGet, url, http.StatusOK, nil); err != nil {
		return fmt.Errorf("timed out waiting for Keycloak at %s", url)
	}
	logger.Step("      Keycloak ready")
	return nil
}

// ── Kong key fetch ────────────────────────────────────────────────────────────

func (r *Runner) fetchAndApplyKongKey(ctx context.Context) error {
	if r.dryRun {
		logger.Step("      [dry-run] would fetch Kong RS256 key from Keycloak")
		return nil
	}

	jwksURL := fmt.Sprintf("http://localhost:%d/realms/%s/protocol/openid-connect/certs",
		r.cfg.KeycloakNodePort, r.cfg.KeycloakRealm)
	pemKey, err := util.FetchRS256PEM(ctx, jwksURL)
	if err != nil {
		return fmt.Errorf("fetch Kong key from Keycloak: %w", err)
	}
	logger.Debug("fetched RS256 public key from Keycloak JWKS")

	creds, err := r.writeKongKeyToCreds(pemKey)
	if err != nil {
		return fmt.Errorf("update credentials with Kong key: %w", err)
	}

	path, err := WriteSecretsValuesFile(r.outputDir, creds, r.cfg.NexusRegistry)
	if err != nil {
		return fmt.Errorf("regenerate values-secrets.yaml: %w", err)
	}
	r.secretsFile = path
	r.creds = creds
	logger.Step("      updated secrets.yaml and values-secrets.yaml")
	return nil
}

// writeKongKeyToCreds writes the PEM key into secrets.yaml and reloads it.
func (r *Runner) writeKongKeyToCreds(pemKey string) (*credentials.Credentials, error) {
	data, err := os.ReadFile(r.secretsPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", r.secretsPath, err)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse %s: %w", r.secretsPath, err)
	}

	kong, ok := raw["kong"].(map[string]interface{})
	if !ok {
		kong = map[string]interface{}{}
		raw["kong"] = kong
	}
	kong["keycloakPublicKey"] = pemKey

	updated, err := yaml.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal credentials: %w", err)
	}
	if err := os.WriteFile(r.secretsPath, updated, 0600); err != nil {
		return nil, fmt.Errorf("write %s: %w", r.secretsPath, err)
	}

	return credentials.Load(r.secretsPath)
}

// ── Kong wait ─────────────────────────────────────────────────────────────────

func (r *Runner) waitForKong(ctx context.Context) error {
	logger.Step("      waiting for Kong deployment to become available (up to 2 min)...")
	return r.kubectl(ctx, "wait",
		"-n", r.cfg.NamespaceInfra,
		"deployment/kong",
		"--for=condition=Available",
		"--timeout=120s",
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
		logger.Step("      vault-init.json already exists — skipping init")
		return nil
	}

	if r.dryRun {
		logger.Step("      kubectl exec -n %s %s -- vault operator init -format=json > %s",
			r.cfg.NamespaceInfra, r.cfg.VaultPod, outPath)
		return nil
	}

	out, err := util.RunCmdOutput(ctx,
		"kubectl", "exec", "-n", r.cfg.NamespaceInfra, r.cfg.VaultPod, "--",
		"vault", "operator", "init", "-format=json",
	)
	if err != nil {
		return fmt.Errorf("vault init: %w", err)
	}

	if err := os.WriteFile(outPath, out, 0600); err != nil {
		return fmt.Errorf("write vault-init.json: %w", err)
	}

	logger.Step("      saved to %s — keep this file safe!", outPath)
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

	yamlBody := fmt.Sprintf(`apiVersion: v1
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
		logger.Step("      kubectl apply -f - (vault-unseal-keys secret in %s)", r.cfg.NamespaceInfra)
		return nil
	}

	return util.RunCmdStdin(ctx, strings.NewReader(yamlBody), "kubectl", "apply", "-f", "-")
}

func (r *Runner) unsealVault(ctx context.Context) error {
	init, err := r.loadVaultInit()
	if err != nil {
		return err
	}

	for i, key := range init.UnsealKeysB64[:3] {
		logger.Step("      unsealing with key %d/3...", i+1)
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
		logger.Step("      [dry-run] would get nginx ingress ClusterIP and apply coredns-patch.yaml")
		return nil
	}

	// Get nginx ingress controller ClusterIP
	out, err := util.RunCmdOutput(ctx,
		"kubectl", "get", "svc",
		"-n", "ingress-nginx",
		"-l", "app.kubernetes.io/component=controller",
		"-o", "jsonpath={.items[0].spec.clusterIP}",
	)

	ip := strings.TrimSpace(string(out))
	if err != nil || ip == "" {
		logger.Warn("could not auto-detect nginx ingress ClusterIP (%v)", err)
		logger.StepMsg("      Update the IP in k8s/coredns-patch.yaml manually, then:")
		logger.StepMsg("      kubectl apply -f k8s/coredns-patch.yaml")
		return nil
	}

	logger.Step("      nginx ingress ClusterIP: %s", ip)

	patchYAML, err := os.ReadFile(coreDNSPatch)
	if err != nil {
		logger.Warn("could not read %s: %v", coreDNSPatch, err)
		return nil
	}

	updated := substituteCoreDNSIP(string(patchYAML), ip)

	if err := util.RunCmdStdin(ctx, bytes.NewReader([]byte(updated)), "kubectl", "apply", "-f", "-"); err != nil {
		return fmt.Errorf("kubectl apply coredns patch: %w", err)
	}
	return nil
}

func (r *Runner) waitForCoreDNS(ctx context.Context) error {
	if err := r.kubectl(ctx, "rollout", "restart", "deployment/coredns", "-n", "kube-system"); err != nil {
		return fmt.Errorf("restart coredns: %w", err)
	}
	return r.kubectl(ctx, "rollout", "status", "deployment/coredns", "-n", "kube-system", "--timeout=60s")
}

func (r *Runner) verifyAllDNS(ctx context.Context) error {
	if r.dryRun {
		logger.Step("      [dry-run] would verify .local DNS resolution inside cluster")
		return nil
	}

	domains := []string{
		"jeeb-dev.local", "api.jeeb-dev.local", "auth.jeeb-dev.local", "learning.jeeb-dev.local",
		"jenkins.jeeb.local", "nexus.jeeb.local", "sonarqube.jeeb.local", "vault.jeeb.local",
		"grafana.jeeb.local", "rancher.jeeb-infra.local",
	}

	pass, fail := 0, 0
	for i, domain := range domains {
		podName := fmt.Sprintf("dns-setup-%d", i+1)
		err := util.RunCmd(ctx, "kubectl", "run", podName,
			"--image=busybox", "--restart=Never", "--rm", "--attach",
			"--", "nslookup", domain,
		)
		if err != nil {
			logger.Warn("%s — not resolved", domain)
			fail++
		} else {
			logger.Step("      [OK]   %s", domain)
			pass++
		}
	}

	logger.Step("      DNS: %d/%d domains resolved", pass, pass+fail)
	if fail > 0 {
		logger.StepMsg("      Some domains failed. Run 'k8s-manager maintain' to diagnose.")
	}
	return nil // non-fatal
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
	return util.DryRunOrExec(ctx, r.dryRun, "kubectl", args...)
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
		logger.Step("      vault policy write %s -", name)
		return nil
	}
	return util.RunCmdStdin(ctx, strings.NewReader(policy),
		"kubectl", "exec", "-i", "-n", r.cfg.NamespaceInfra, r.cfg.VaultPod, "--",
		"env",
		fmt.Sprintf("VAULT_ADDR=%s", r.cfg.VaultAddr),
		"VAULT_TOKEN="+token,
		"vault", "policy", "write", name, "-",
	)
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
	repoRoot := filepath.Dir(filepath.Dir(r.chartsDir))
	groovyPath := filepath.Join(repoRoot, "jenkins", "jobs", "seed.groovy")

	if _, err := os.Stat(groovyPath); err != nil {
		logger.Warn("seed.groovy not found at %s", groovyPath)
		logger.StepMsg("      Run 'k8s-manager seed --groovy-path <path>' after Jenkins is ready")
		return nil
	}

	return jenkins.NewSeeder(r.cfg, r.creds, groovyPath, "", r.dryRun).Run(ctx)
}

func (r *Runner) printAccessTable() {
	logger.Step(`
  %s (infra):
    Jenkins    http://localhost:%d
    Nexus      http://localhost:%d
    SonarQube  http://localhost:%d
    Kong       http://localhost:%d
    Vault      http://localhost:%d
    Rancher    https://localhost:%d

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
		r.cfg.RancherNodePort,
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
