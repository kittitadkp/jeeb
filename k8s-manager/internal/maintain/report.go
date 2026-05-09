package maintain

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"k8s-manager/internal/config"
	"k8s-manager/internal/healthcheck"
)

// RunReport runs all health checks and prints a diagnosis block for each failure.
// No commands are executed — this is a read-only report.
func RunReport(ctx context.Context, cfg *config.ClusterConfig, outputDir string) error {
	results := healthcheck.RunAll(ctx, cfg, outputDir)

	fmt.Println("=== Cluster Maintenance Report ===")
	fmt.Println()

	failures := 0
	for _, r := range results {
		if r.Passed {
			fmt.Printf("[PASS] %s\n", r.Name)
		} else {
			failures++
			fmt.Printf("[FAIL] %s — %s\n", r.Name, r.Detail)
			printDiagnosis(r.Name, cfg, outputDir)
			fmt.Println()
		}
	}

	if failures == 0 {
		fmt.Println()
		fmt.Println("All checks passed. Cluster is healthy.")
		return nil
	}

	fmt.Printf("%d check(s) failed.\n", failures)
	return nil
}

// ── diagnosis printer ─────────────────────────────────────────────────────────

func printDiagnosis(name string, cfg *config.ClusterConfig, outputDir string) {
	if strings.HasPrefix(name, "DNS: ") {
		domain := strings.TrimPrefix(name, "DNS: ")
		fmt.Println("  Diagnosis : CoreDNS patch not applied or nginx ingress IP changed")
		fmt.Println("  Fix       :")
		fmt.Println("    kubectl apply -f k8s/coredns-patch.yaml")
		fmt.Println("    kubectl rollout restart deployment/coredns -n kube-system")
		fmt.Printf("    kubectl run -it --rm dns-test --image=busybox --restart=Never -- nslookup %s\n", domain)
		return
	}

	vaultInitPath := filepath.Join(outputDir, "vault-init.json")
	secretPath := config.KVCLIPath(cfg.VaultPathBackend)

	type entry struct{ diagnosis, fix string }
	entries := map[string]entry{
		"Pod health": {
			"One or more pods are in a failed state",
			"kubectl get pods -A\n" +
				"kubectl describe pod <pod-name> -n <namespace>\n" +
				"kubectl logs <pod-name> -n <namespace> --previous",
		},
		"Keycloak": {
			fmt.Sprintf("Keycloak is not reachable on port %d", cfg.KeycloakNodePort),
			fmt.Sprintf("kubectl get pods -n %s -l app=keycloak\n", cfg.NamespaceDev) +
				fmt.Sprintf("kubectl logs -n %s deployment/keycloak --tail=50", cfg.NamespaceDev),
		},
		"Vault": {
			"Vault may be sealed or unreachable",
			fmt.Sprintf("# Check status:\nkubectl exec -n %s %s -- vault status\n\n", cfg.NamespaceInfra, cfg.VaultPod) +
				fmt.Sprintf("# Unseal (keys in %s):\n", vaultInitPath) +
				fmt.Sprintf("kubectl exec -n %s %s -- vault operator unseal <key1>\n", cfg.NamespaceInfra, cfg.VaultPod) +
				fmt.Sprintf("kubectl exec -n %s %s -- vault operator unseal <key2>\n", cfg.NamespaceInfra, cfg.VaultPod) +
				fmt.Sprintf("kubectl exec -n %s %s -- vault operator unseal <key3>", cfg.NamespaceInfra, cfg.VaultPod),
		},
		"Kong": {
			"Kong may be missing the Keycloak public key or not running",
			fmt.Sprintf("kubectl get pods -n %s -l app=kong\n", cfg.NamespaceInfra) +
				"k8s-manager kong-key   # refresh Keycloak public key",
		},
		"Jenkins": {
			fmt.Sprintf("Jenkins is not reachable on port %d", cfg.JenkinsNodePort),
			fmt.Sprintf("kubectl get pods -n %s -l app=jenkins\n", cfg.NamespaceInfra) +
				fmt.Sprintf("kubectl logs -n %s deployment/jenkins --tail=50", cfg.NamespaceInfra),
		},
		"Vault secrets": {
			"Vault secrets not readable — Vault may need re-configuration",
			fmt.Sprintf("# Read using root token from %s:\n", vaultInitPath) +
				fmt.Sprintf("kubectl exec -n %s %s -- env VAULT_ADDR=%s VAULT_TOKEN=<root-token> vault kv get %s",
					cfg.NamespaceInfra, cfg.VaultPod, cfg.VaultAddr, secretPath),
		},
	}

	d, ok := entries[name]
	if !ok {
		return
	}
	fmt.Printf("  Diagnosis : %s\n", d.diagnosis)
	fmt.Println("  Fix       :")
	for _, line := range strings.Split(d.fix, "\n") {
		fmt.Printf("    %s\n", line)
	}
}

