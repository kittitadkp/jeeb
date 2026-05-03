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

	"k8s-manager/internal/credentials"
)

type Runner struct {
	creds       *credentials.Credentials
	chartsDir   string
	outputDir   string
	dryRun      bool
	secretsFile string // path to generated values-secrets.yaml
}

func NewRunner(creds *credentials.Credentials, chartsDir, outputDir string, dryRun bool) *Runner {
	return &Runner{
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
	path, err := WriteSecretsValuesFile(r.outputDir, r.creds)
	if err != nil {
		return err
	}
	r.secretsFile = path
	fmt.Printf("      wrote %s\n", path)
	return nil
}

// Deploy runs helm upgrade --install for the given targets (infra, app, obs).
// If targets is empty all three charts are deployed.
// It generates values-secrets.yaml first, same as the full setup flow.
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
		{"Deploy jeeb-app      (MongoDB, Keycloak, backend, frontend)", "app", r.deployApp},
		{"Deploy jeeb-learning", "learning", r.deployLearning},
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
		{"Deploy jeeb-app (MongoDB, Keycloak, backend, frontend)", r.deployApp},
		{"Deploy jeeb-learning", r.deployLearning},
		{"Deploy jeeb-obs (Prometheus, Loki, Grafana)", r.deployObs},
		{"Wait for Vault pod ready", r.waitForVault},
		{"Initialize Vault", r.initVault},
		{"Store unseal keys in Kubernetes secret", r.storeUnsealKeys},
		{"Unseal Vault", r.unsealVault},
		{"Configure Vault (KV engine, policies, K8s auth roles)", r.configureVault},
		{"Patch CoreDNS for .local DNS", r.patchCoreDNS},
	}

	fmt.Println("=== New Cluster Setup ===")
	fmt.Printf("Charts directory : %s\n", r.chartsDir)
	fmt.Printf("Output directory : %s\n", r.outputDir)
	if r.dryRun {
		fmt.Println("DRY RUN — commands will be printed, not executed")
	}
	fmt.Println()

	fmt.Println("[0/10] Generate values-secrets.yaml from credentials")
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

// ── Step 1 ───────────────────────────────────────────────────────────────────

func (r *Runner) deployInfra(ctx context.Context) error {
	return r.helm(ctx, "upgrade", "--install", "jeeb-infra",
		filepath.Join(r.chartsDir, "jeeb-infra"),
		"--namespace", "jeeb-infra",
		"--create-namespace",
		"-f", r.secretsFile,
	)
}

// ── Step 2 ───────────────────────────────────────────────────────────────────

func (r *Runner) deployApp(ctx context.Context) error {
	return r.helm(ctx, "upgrade", "--install", "jeeb-dev",
		filepath.Join(r.chartsDir, "jeeb-app"),
		"--namespace", "jeeb-dev",
		"--create-namespace",
		"-f", filepath.Join(r.chartsDir, "jeeb-app", "values-dev.yaml"),
		"-f", r.secretsFile,
	)
}

// ── Step 3 ───────────────────────────────────────────────────────────────────

func (r *Runner) deployLearning(ctx context.Context) error {
	return r.helm(ctx, "upgrade", "--install", "jeeb-learning",
		filepath.Join(r.chartsDir, "jeeb-learning"),
		"--namespace", "jeeb-dev",
		"--create-namespace",
		"-f", filepath.Join(r.chartsDir, "jeeb-learning", "values-dev.yaml"),
		"-f", r.secretsFile,
	)
}

// ── Step 4 ───────────────────────────────────────────────────────────────────

func (r *Runner) deployObs(ctx context.Context) error {
	return r.helm(ctx, "upgrade", "--install", "jeeb-obs",
		filepath.Join(r.chartsDir, "jeeb-obs"),
		"--namespace", "jeeb-obs",
		"--create-namespace",
		"-f", r.secretsFile,
	)
}

// ── Step 5 ───────────────────────────────────────────────────────────────────

func (r *Runner) waitForVault(ctx context.Context) error {
	return r.kubectl(ctx, "wait", "pod/vault-0",
		"-n", "jeeb-infra",
		"--for=condition=Ready",
		"--timeout=120s",
	)
}

// ── Step 6 ───────────────────────────────────────────────────────────────────

type vaultInitOutput struct {
	UnsealKeysB64 []string `json:"unseal_keys_b64"`
	RootToken     string   `json:"root_token"`
}

func (r *Runner) initVault(ctx context.Context) error {
	outPath := filepath.Join(r.outputDir, "vault-init.json")

	// skip if already initialised
	if _, err := os.Stat(outPath); err == nil {
		fmt.Printf("      vault-init.json already exists — skipping init\n")
		return nil
	}

	if r.dryRun {
		fmt.Printf("      kubectl exec -n jeeb-infra vault-0 -- vault operator init -format=json > %s\n", outPath)
		return nil
	}

	out, err := exec.CommandContext(ctx,
		"kubectl", "exec", "-n", "jeeb-infra", "vault-0", "--",
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

// ── Step 6 ───────────────────────────────────────────────────────────────────

func (r *Runner) storeUnsealKeys(ctx context.Context) error {
	init, err := r.loadVaultInit()
	if err != nil {
		return err
	}

	if len(init.UnsealKeysB64) < 3 {
		return fmt.Errorf("expected at least 3 unseal keys, got %d", len(init.UnsealKeysB64))
	}

	return r.kubectl(ctx,
		"create", "secret", "generic", "vault-unseal-keys",
		"-n", "jeeb-infra",
		"--from-literal=key1="+init.UnsealKeysB64[0],
		"--from-literal=key2="+init.UnsealKeysB64[1],
		"--from-literal=key3="+init.UnsealKeysB64[2],
		"--dry-run=client", "-o", "yaml",
	)
	// pipe to kubectl apply is handled separately below
}

func (r *Runner) storeUnsealKeysApply(ctx context.Context) error {
	init, err := r.loadVaultInit()
	if err != nil {
		return err
	}

	yaml := fmt.Sprintf(`apiVersion: v1
kind: Secret
metadata:
  name: vault-unseal-keys
  namespace: jeeb-infra
stringData:
  key1: %s
  key2: %s
  key3: %s
`, init.UnsealKeysB64[0], init.UnsealKeysB64[1], init.UnsealKeysB64[2])

	if r.dryRun {
		fmt.Println("      kubectl apply -f - (vault-unseal-keys secret)")
		return nil
	}

	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(yaml)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ── Step 7 ───────────────────────────────────────────────────────────────────

func (r *Runner) unsealVault(ctx context.Context) error {
	// override storeUnsealKeys step to use apply approach
	if err := r.storeUnsealKeysApply(ctx); err != nil {
		return fmt.Errorf("store unseal keys: %w", err)
	}

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

// ── Step 8 ───────────────────────────────────────────────────────────────────

func (r *Runner) configureVault(ctx context.Context) error {
	init, err := r.loadVaultInit()
	if err != nil {
		return err
	}
	token := init.RootToken

	type kv struct{ path, key, val string }
	secrets := []kv{
		// backend
		{"secret/jeeb/backend/develop", "PORT", "8080"},
		{"secret/jeeb/backend/develop", "LOG_LEVEL", "INFO"},
		{"secret/jeeb/backend/develop", "MONGO_DATABASE", "jeeb"},
		{"secret/jeeb/backend/develop", "MONGO_URI", fmt.Sprintf("mongodb://%s:%s@mongodb.jeeb-dev.svc.cluster.local:27017/jeeb?authSource=admin", r.creds.MongoDBUsername, r.creds.MongoDBPassword)},
		{"secret/jeeb/backend/develop", "KEYCLOAK_URL", "http://keycloak.jeeb-dev.svc.cluster.local:8080"},
		{"secret/jeeb/backend/develop", "KEYCLOAK_REALM", "jeeb"},
		{"secret/jeeb/backend/develop", "KEYCLOAK_CLIENT_ID", "jeeb-app"},
		// frontend
		{"secret/jeeb/frontend/develop", "VITE_KEYCLOAK_URL", "http://localhost:30081"},
		{"secret/jeeb/frontend/develop", "VITE_KEYCLOAK_REALM", "jeeb"},
		{"secret/jeeb/frontend/develop", "VITE_KEYCLOAK_CLIENT_ID", "jeeb-app"},
		{"secret/jeeb/frontend/develop", "VITE_API_URL", "http://localhost:30080"},
		// learning backend
		{"secret/jeeb/learning/backend/develop", "PORT", "8080"},
		{"secret/jeeb/learning/backend/develop", "LOG_LEVEL", "INFO"},
		{"secret/jeeb/learning/backend/develop", "MONGO_DATABASE", "jeeb_learning"},
		{"secret/jeeb/learning/backend/develop", "MONGO_URI", fmt.Sprintf("mongodb://%s:%s@mongodb.jeeb-dev.svc.cluster.local:27017/jeeb_learning?authSource=admin", r.creds.MongoDBUsername, r.creds.MongoDBPassword)},
		{"secret/jeeb/learning/backend/develop", "KEYCLOAK_URL", "http://keycloak.jeeb-dev.svc.cluster.local:8080"},
		{"secret/jeeb/learning/backend/develop", "KEYCLOAK_REALM", "jeeb"},
		{"secret/jeeb/learning/backend/develop", "KEYCLOAK_CLIENT_ID", "jeeb-app"},
		// learning frontend
		{"secret/jeeb/learning/frontend/develop", "VITE_KEYCLOAK_URL", "http://localhost:30081"},
		{"secret/jeeb/learning/frontend/develop", "VITE_KEYCLOAK_REALM", "jeeb"},
		{"secret/jeeb/learning/frontend/develop", "VITE_KEYCLOAK_CLIENT_ID", "jeeb-app"},
		{"secret/jeeb/learning/frontend/develop", "VITE_API_URL", "http://localhost:30086"},
	}

	steps := []func() error{
		func() error {
			return r.vaultExec(ctx, token, "secrets", "enable", "-path=secret", "kv-v2")
		},
		func() error {
			// write secrets per path (group by path)
			written := map[string]bool{}
			for _, s := range secrets {
				if written[s.path] {
					continue
				}
				// collect all kv pairs for this path
				var args []string
				args = append(args, "kv", "put", s.path)
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
		func() error {
			return r.vaultWritePolicy(ctx, token, "backend-policy",
				`path "secret/data/jeeb/backend/develop" { capabilities = ["read"] }`)
		},
		func() error {
			return r.vaultExec(ctx, token, "write", "auth/kubernetes/role/backend",
				"bound_service_account_names=backend",
				"bound_service_account_namespaces=jeeb-dev",
				"policies=backend-policy",
				"ttl=1h")
		},
		func() error {
			return r.vaultWritePolicy(ctx, token, "frontend-policy",
				`path "secret/data/jeeb/frontend/develop" { capabilities = ["read"] }`)
		},
		func() error {
			return r.vaultExec(ctx, token, "write", "auth/kubernetes/role/frontend",
				"bound_service_account_names=frontend",
				"bound_service_account_namespaces=jeeb-dev",
				"policies=frontend-policy",
				"ttl=1h")
		},
		func() error {
			return r.vaultWritePolicy(ctx, token, "learning-backend-policy",
				`path "secret/data/jeeb/learning/backend/develop" { capabilities = ["read"] }`)
		},
		func() error {
			return r.vaultExec(ctx, token, "write", "auth/kubernetes/role/learning-backend",
				"bound_service_account_names=learning-backend",
				"bound_service_account_namespaces=jeeb-dev",
				"policies=learning-backend-policy",
				"ttl=1h")
		},
		func() error {
			return r.vaultWritePolicy(ctx, token, "learning-frontend-policy",
				`path "secret/data/jeeb/learning/frontend/develop" { capabilities = ["read"] }`)
		},
		func() error {
			return r.vaultExec(ctx, token, "write", "auth/kubernetes/role/learning-frontend",
				"bound_service_account_names=learning-frontend",
				"bound_service_account_namespaces=jeeb-dev",
				"policies=learning-frontend-policy",
				"ttl=1h")
		},
	}

	for _, s := range steps {
		if err := s(); err != nil {
			// "already enabled" errors are non-fatal
			if !strings.Contains(err.Error(), "already enabled") &&
				!strings.Contains(err.Error(), "path is already in use") {
				return err
			}
		}
	}
	return nil
}

// ── Step 9 ───────────────────────────────────────────────────────────────────

func (r *Runner) patchCoreDNS(ctx context.Context) error {
	coreDNSPatch := filepath.Join(filepath.Dir(r.chartsDir), "coredns-patch.yaml")

	if r.dryRun {
		fmt.Println("      kubectl get svc -n jeeb-dev -l app.kubernetes.io/name=ingress-nginx -o jsonpath='{.items[0].spec.clusterIP}'")
		fmt.Printf("      [substitute IP in %s and pipe to] kubectl apply -f -\n", coreDNSPatch)
		return nil
	}

	out, err := exec.CommandContext(ctx,
		"kubectl", "get", "svc",
		"-n", "jeeb-dev",
		"-l", "app.kubernetes.io/name=ingress-nginx",
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

	// Replace the placeholder IP with the detected one. The file contains a
	// single IP address on each hosts line — replace any 4-octet pattern.
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

// substituteCoreDNSIP replaces the existing IP on each 'hosts' entry line
// with newIP. It matches the first IPv4 address at the start of any data line.
func substituteCoreDNSIP(yaml, newIP string) string {
	lines := strings.Split(yaml, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Lines inside the Corefile hosts block start with an IP
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

func (r *Runner) helm(ctx context.Context, args ...string) error {
	if r.dryRun {
		fmt.Printf("      helm %s\n", strings.Join(args, " "))
		return nil
	}
	cmd := exec.CommandContext(ctx, "helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

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
		"exec", "-i", "-n", "jeeb-infra", "vault-0", "--",
		"env",
		"VAULT_ADDR=http://127.0.0.1:8200",
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
		"kubectl", "exec", "-i", "-n", "jeeb-infra", "vault-0", "--",
		"env",
		"VAULT_ADDR=http://127.0.0.1:8200",
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

func (r *Runner) printAccessTable() {
	fmt.Print(`
  jeeb-infra:
    Jenkins    http://localhost:30082
    SonarQube  http://localhost:30090
    Nexus      http://localhost:30083
    Vault      http://localhost:30091

  jeeb-dev:
    Frontend          http://localhost:30000
    Backend           http://localhost:30080
    Keycloak          http://localhost:30081
    Learning frontend http://localhost:30087
    Learning backend  http://localhost:30086

  jeeb-obs:
    Grafana    http://localhost:30092
    Prometheus http://localhost:30093
`)
}
