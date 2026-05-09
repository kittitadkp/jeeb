package healthcheck

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"k8s-manager/internal/config"
)

// Result holds the outcome of a single health check.
type Result struct {
	Name   string
	Passed bool
	Detail string
}

// ErrChecksFailed is returned by commands when one or more checks fail.
// It carries no message because results are already printed to stdout.
var ErrChecksFailed = errors.New("")

// RunAll executes every health check and returns results in order.
func RunAll(ctx context.Context, cfg *config.ClusterConfig, outputDir string) []Result {
	return []Result{
		checkPods(ctx),
		checkHTTP(ctx, "Keycloak",
			fmt.Sprintf("http://localhost:%d/realms/%s", cfg.KeycloakNodePort, cfg.KeycloakRealm),
			assertKeycloak(cfg.KeycloakRealm)),
		checkHTTP(ctx, "Vault",
			fmt.Sprintf("http://localhost:%d/v1/sys/health", cfg.VaultNodePort),
			assertVault),
		checkHTTP(ctx, "Kong",
			fmt.Sprintf("http://localhost:%d/health", cfg.KongNodePort),
			nil),
		checkHTTP(ctx, "Jenkins",
			fmt.Sprintf("http://localhost:%d/login", cfg.JenkinsNodePort),
			nil),
		checkDNS(ctx, "auth.jeeb-dev.local", 1),
		checkDNS(ctx, "jenkins.jeeb.local", 2),
		checkDNS(ctx, "grafana.jeeb.local", 3),
		checkVaultSecrets(ctx, cfg, outputDir),
	}
}

// Print writes a pass/fail table to stdout.
func Print(results []Result) {
	fmt.Println("=== Jeeb Cluster Health Check ===")
	fmt.Println()

	nameW := 0
	for _, r := range results {
		if len(r.Name) > nameW {
			nameW = len(r.Name)
		}
	}
	nameW += 2

	for _, r := range results {
		tag := "[PASS]"
		if !r.Passed {
			tag = "[FAIL]"
		}
		fmt.Printf("%s %-*s %s\n", tag, nameW, r.Name, r.Detail)
	}

	failed := FailCount(results)
	fmt.Println()
	if failed > 0 {
		fmt.Printf("%d check(s) failed. Run `k8s-manager maintain` for diagnosis.\n", failed)
	} else {
		fmt.Println("All checks passed.")
	}
}

// AllPassed returns true only when every result passed.
func AllPassed(results []Result) bool {
	return FailCount(results) == 0
}

// FailCount returns the number of failed results.
func FailCount(results []Result) int {
	n := 0
	for _, r := range results {
		if !r.Passed {
			n++
		}
	}
	return n
}

// ── individual checks ─────────────────────────────────────────────────────────

func checkPods(ctx context.Context) Result {
	out, err := exec.CommandContext(ctx, "kubectl", "get", "pods", "-A", "--no-headers").Output()
	if err != nil {
		return Result{"Pod health", false, "kubectl get pods -A failed"}
	}
	var bad []string
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}
		if strings.Contains(line, "CrashLoopBackOff") ||
			strings.Contains(line, "ImagePullBackOff") ||
			strings.Contains(line, "OOMKilled") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				bad = append(bad, fields[0]+"/"+fields[1])
			}
		}
	}
	if len(bad) > 0 {
		return Result{"Pod health", false, "unhealthy: " + strings.Join(bad, ", ")}
	}
	return Result{"Pod health", true, "all pods healthy"}
}

// assertFn validates an HTTP response body; returns an error message or "".
type assertFn func(body string) string

func assertKeycloak(realm string) assertFn {
	return func(body string) string {
		var data map[string]json.RawMessage
		if err := json.Unmarshal([]byte(body), &data); err != nil {
			return "could not parse Keycloak response"
		}
		var got string
		_ = json.Unmarshal(data["realm"], &got)
		if got != realm {
			return fmt.Sprintf("realm=%q, expected %q", got, realm)
		}
		return ""
	}
}

func assertVault(body string) string {
	var data struct {
		Sealed bool `json:"sealed"`
	}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return "could not parse Vault response"
	}
	if data.Sealed {
		return "Vault is sealed"
	}
	return ""
}

func checkHTTP(ctx context.Context, name, url string, assert assertFn) Result {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{name, false, fmt.Sprintf("build request: %v", err)}
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Result{name, false, fmt.Sprintf("%s — %v", url, err)}
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return Result{name, false, fmt.Sprintf("%s → HTTP %d", url, resp.StatusCode)}
	}

	if assert != nil {
		body, _ := io.ReadAll(resp.Body)
		if msg := assert(string(body)); msg != "" {
			return Result{name, false, msg}
		}
	}

	return Result{name, true, fmt.Sprintf("HTTP %d", resp.StatusCode)}
}

// checkDNS launches a unique busybox pod to run nslookup inside the cluster.
// idx makes the pod name deterministic within one RunAll call; a nanosecond
// suffix prevents collision across concurrent or rapid sequential calls.
func checkDNS(ctx context.Context, domain string, idx int) Result {
	podName := fmt.Sprintf("dns-check-%d-%d", idx, time.Now().UnixNano()%1_000_000_000)
	cmd := exec.CommandContext(ctx,
		"kubectl", "run", podName,
		"--image=busybox", "--restart=Never", "--rm", "--attach",
		"--", "nslookup", domain,
	)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		return Result{"DNS: " + domain, false, "not resolved inside cluster"}
	}
	return Result{"DNS: " + domain, true, "resolved"}
}

// checkVaultSecrets reads the root token from vault-init.json and calls
// `vault kv get` via kubectl exec — no shell, no token injection risk.
func checkVaultSecrets(ctx context.Context, cfg *config.ClusterConfig, outputDir string) Result {
	vaultInitPath := filepath.Join(outputDir, "vault-init.json")
	data, err := os.ReadFile(vaultInitPath)
	if err != nil {
		return Result{"Vault secrets", false, "vault-init.json not found — run setup first"}
	}
	var v struct {
		RootToken string `json:"root_token"`
	}
	if err := json.Unmarshal(data, &v); err != nil || v.RootToken == "" {
		return Result{"Vault secrets", false, "vault-init.json missing root_token"}
	}

	secretPath := config.KVCLIPath(cfg.VaultPathBackend)
	cmd := exec.CommandContext(ctx,
		"kubectl", "exec", "-n", cfg.NamespaceInfra, cfg.VaultPod, "--",
		"env",
		"VAULT_ADDR="+cfg.VaultAddr,
		"VAULT_TOKEN="+v.RootToken,
		"vault", "kv", "get", secretPath,
	)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		return Result{"Vault secrets", false, "backend secrets not readable"}
	}
	return Result{"Vault secrets", true, "backend secrets readable"}
}
