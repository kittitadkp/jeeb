# Plan: Rewrite k8s-manager for Correct Fresh Cluster Setup

## Context

The cluster setup currently fails in three ways on a fresh cluster:

1. **Keycloak realm is never imported** — `jeeb-data/templates/keycloak/configmap.yaml` calls
   `.Files.Get "files/realm-jeeb.json"` but there is no `files/` directory under `jeeb-data/`.
   The file lives under `jeeb-app/files/`. Result: empty ConfigMap, Keycloak starts with no `jeeb`
   realm, so Vault policies fail and the frontend can't authenticate.

2. **Kong crashes permanently** — `jeeb-infra` deploys Kong (needs RSA public key) before
   `jeeb-data` deploys Keycloak (the key source). The `kong-key` command exists as a standalone fix
   but is never called during `setup`. Kong stays in CrashLoopBackOff forever.

3. **CoreDNS patch fails silently** — `patchCoreDNS` looks for a service labelled
   `app.kubernetes.io/name=ingress-nginx` in `jeeb-dev` namespace. Docker Desktop has no nginx
   ingress controller, so the IP is never found and `.local` DNS never works.

After the cluster reset this plan fixes all three permanently so `k8s-manager setup` completes in
one shot.

---

## Files to Change

| File | Change |
|------|--------|
| `k8s/charts/jeeb-data/files/realm-jeeb.json` | **NEW** — copy of `jeeb-app/files/realm-jeeb.json` |
| `k8s-manager/internal/setup/setup.go` | Revised `Run()` order + new steps |
| `k8s-manager/internal/kongkey/update.go` | Export `FetchPublicKey()` + `UpdateCredsFile()` |
| `k8s-manager/cmd/k8s-manager/main.go` | Add `check` + `maintain` commands; pass `credsPath` to `NewRunner` |
| `k8s-manager/scripts/health_check.py` | **NEW** — Python health check script |
| `k8s-manager/internal/maintain/report.go` | **NEW** — maintain diagnosis package |
| `.claude/skills/check/skill.md` | **NEW** — `/check` skill |

---

## Implementation Steps

### Step 1 — Fix Keycloak realm ConfigMap (Helm chart)

Create `k8s/charts/jeeb-data/files/realm-jeeb.json` as a copy of
`k8s/charts/jeeb-app/files/realm-jeeb.json`.

The existing template in `jeeb-data/templates/keycloak/configmap.yaml` is already correct:
```yaml
data:
  realm-jeeb.json: |
{{ .Files.Get "files/realm-jeeb.json" | indent 4 }}
```
It just needs the file to exist in the right chart's `files/` directory. No template change needed.

---

### Step 2 — Export `FetchPublicKey` from kongkey package

**File:** `k8s-manager/internal/kongkey/update.go`

Add a public function so `setup.go` can call it without duplicating the HTTP + JWK-to-PEM logic:

```go
// FetchPublicKey fetches the RS256 public key from Keycloak's JWKS endpoint.
// It is exported so the setup runner can integrate key fetch into the setup flow.
func FetchPublicKey(keycloakNodePort int, realm string) (string, error) {
    u := &Updater{cfg: &config.ClusterConfig{KeycloakNodePort: keycloakNodePort, KeycloakRealm: realm}}
    return u.fetchFromKeycloak()
}
```

---

### Step 3 — Rewrite `setup.go` `Run()` with correct step order

**File:** `k8s-manager/internal/setup/setup.go`

#### 3a. Pre-flight check (new first step, before [0])

Before writing any files, verify the environment is sane:
```go
func (r *Runner) preflight(ctx context.Context) error {
    // 1. kubectl cluster-info — confirm API server reachable
    // 2. Check context is docker-desktop (warn if not, don't block)
    // 3. Warn if vault-init.json already exists (stale from old cluster)
    //    — print: "vault-init.json found from previous run. Delete it or Vault init will be skipped."
    // 4. Verify helm is on PATH
}
```

#### 3a. New step order in `Run()`

`jeeb-app` and `jeeb-learning` are **removed from the automated setup** — they require images
built and pushed to Nexus by Jenkins pipelines before they can be deployed. The operator runs
`k8s-manager deploy app learning` manually after pipelines succeed.

```
[-1] Pre-flight           (cluster reachable, helm on PATH, warn if vault-init.json exists)
[0]  Generate values-secrets.yaml (Kong key empty — OK for now)
[1]  Install nginx ingress controller (if not already present)
[2]  Install Rancher + cert-manager (reuse existing rancher.Deployer)
[3]  Deploy jeeb-infra    (Kong crashes — acceptable; Vault/Jenkins/Nexus start fine)
[4]  Deploy jeeb-data     (MongoDB + Keycloak; realm auto-imported from ConfigMap)
[5]  Deploy jeeb-obs
[6]  Initialize Nexus Docker registry (REST API — create hosted repo port 5000)
[7]  Wait for Keycloak ready  (HTTP poll localhost:<nodePort>/realms/<realm>)
[8]  Fetch Kong RS256 key from Keycloak + update credentials.yaml + values-secrets.yaml
[9]  Re-deploy jeeb-infra    (Kong gets the key, stops crashing)
[10] Wait for Kong ready     (kubectl wait deployment/kong --for=condition=Available)
[11] Wait for Vault pod ready
[12] Initialize Vault
[13] Store unseal keys in K8s secret
[14] Unseal Vault
[15] Configure Vault (KV engine, secrets, policies, K8s auth)
[16] Patch CoreDNS for .local DNS  (nginx ingress ClusterIP)
[17] Wait for CoreDNS rollout      (kubectl rollout status deployment/coredns -n kube-system)
[18] Verify DNS for all .local domains
[19] Seed Jenkins
```

#### 3g. New function: `initNexusDockerRepo(ctx)`

Nexus 3 exposes a REST API. After the Nexus pod is ready, this step:
1. Waits for Nexus to be available (HTTP poll `http://localhost:30083/service/rest/v1/status`)
2. Retrieves the initial admin password from the Nexus pod:
   ```
   kubectl exec -n jeeb-infra <nexus-pod> -- cat /nexus-data/admin.password
   ```
3. Changes the admin password to `r.creds.NexusAdminPassword` via REST API
4. Creates the Docker hosted repository (port 5000) via REST API if not already present
5. This unblocks Jenkins pipelines from pushing images

```go
func (r *Runner) initNexusDockerRepo(ctx context.Context) error {
    // Wait for Nexus readiness
    // GET http://localhost:30083/service/rest/v1/status → 200
    // Read initial password from pod
    // Change password via POST /service/rest/v1/security/users/admin/change-password
    // Create docker hosted repo via POST /service/rest/v1/repositories/docker/hosted
    //   body: { name: "jeeb", online: true, storage: {...}, docker: { httpPort: 5000 } }
}
```

> **Rancher note:** The existing `rancher.Deployer` in `internal/rancher/deploy.go` already handles
> cert-manager + Rancher install and patch to NodePort 30443. It is reused directly in `Run()`
> — no new code needed, just wire it in as step [2].

#### 3f. New function: `verifyAllDNS(ctx)`

After CoreDNS restarts, launch a busybox pod and verify that all `.local` hostnames defined in
`coredns-patch.yaml` resolve correctly. Print pass/fail per domain. Non-fatal (warns, doesn't
fail setup) so a single missing entry doesn't block Jenkins seeding.

```go
func (r *Runner) verifyAllDNS(ctx context.Context) error {
    domains := []string{
        "jeeb-dev.local", "api.jeeb-dev.local", "auth.jeeb-dev.local", "learning.jeeb-dev.local",
        "jenkins.jeeb.local", "nexus.jeeb.local", "sonarqube.jeeb.local", "vault.jeeb.local",
        "grafana.jeeb.local", "rancher.jeeb-infra.local",
    }
    // kubectl run busybox --image=busybox --restart=Never --rm -it -- sh -c "nslookup <domain>"
    // for each domain — check exit code 0 means resolved
    // print [PASS]/[WARN] per domain
    // return nil always (non-fatal)
}
```

After setup, the operator:
1. Runs Jenkins pipelines (backend, frontend, learning-backend, learning-frontend)
2. Pipelines build images and push to Nexus (`localhost:30050`)
3. Then runs: `k8s-manager deploy app learning`

#### 3b. New function: `ensureNginxIngress(ctx)`

```go
func (r *Runner) ensureNginxIngress(ctx context.Context) error {
    // Check if already installed
    out, _ := exec.CommandContext(ctx, "kubectl", "get", "ns", "ingress-nginx").Output()
    if strings.Contains(string(out), "ingress-nginx") {
        fmt.Println("      nginx ingress already installed — skipping")
        return nil
    }

    const manifest = "https://raw.githubusercontent.com/kubernetes/ingress-nginx/" +
        "controller-v1.12.2/deploy/static/provider/cloud/deploy.yaml"

    fmt.Printf("      installing nginx ingress from %s\n", manifest)
    if err := r.kubectl(ctx, "apply", "-f", manifest); err != nil {
        return fmt.Errorf("install nginx ingress: %w", err)
    }
    // Wait for controller pod ready
    return r.kubectl(ctx, "wait",
        "-n", "ingress-nginx",
        "deployment/ingress-nginx-controller",
        "--for=condition=Available",
        "--timeout=120s",
    )
}
```

#### 3c. New function: `waitForKeycloak(ctx)`

HTTP-poll `http://localhost:<KeycloakNodePort>/realms/<KeycloakRealm>` with 5-second intervals
up to a 5-minute timeout. Return error on timeout.

```go
func (r *Runner) waitForKeycloak(ctx context.Context) error {
    url := fmt.Sprintf("http://localhost:%d/realms/%s", r.cfg.KeycloakNodePort, r.cfg.KeycloakRealm)
    fmt.Printf("      polling %s (up to 5 min)...\n", url)
    deadline := time.Now().Add(5 * time.Minute)
    for time.Now().Before(deadline) {
        resp, err := http.Get(url) //nolint:noctx
        if err == nil && resp.StatusCode == 200 {
            resp.Body.Close()
            fmt.Println("      Keycloak ready")
            return nil
        }
        if resp != nil { resp.Body.Close() }
        time.Sleep(5 * time.Second)
    }
    return fmt.Errorf("timed out waiting for Keycloak at %s", url)
}
```

#### 3d. New function: `fetchAndApplyKongKey(ctx)`

```go
func (r *Runner) fetchAndApplyKongKey(ctx context.Context) error {
    if r.dryRun {
        fmt.Println("      [dry-run] would fetch Kong key from Keycloak")
        return nil
    }
    pemKey, err := kongkey.FetchPublicKey(r.cfg.KeycloakNodePort, r.cfg.KeycloakRealm)
    if err != nil {
        return fmt.Errorf("fetch Kong key: %w", err)
    }
    fmt.Println("      fetched RS256 public key from Keycloak")

    // Update credentials + values-secrets.yaml via the existing Updater
    updater := kongkey.NewUpdater(r.cfg, r.credsPath, r.chartsDir, r.outputDir, r.dryRun)
    // Only update the creds file + regenerate secrets; we'll call deployInfra ourselves
    creds, err := updater.UpdateCredsFile(pemKey)  // make UpdateCredsFile exported
    if err != nil {
        return err
    }
    path, err := WriteSecretsValuesFile(r.outputDir, creds, r.cfg.NexusRegistry)
    if err != nil {
        return err
    }
    r.secretsFile = path
    r.creds = creds
    fmt.Println("      updated credentials.yaml and values-secrets.yaml")
    return nil
}
```

> Note: `updateCredsFile` in `kongkey/update.go` becomes `UpdateCredsFile` (exported).

#### 3e. Fix `patchCoreDNS` — use nginx ingress namespace

Change the service lookup from `r.cfg.NamespaceDev` + `r.cfg.IngressLabel` to:

```go
out, err := exec.CommandContext(ctx,
    "kubectl", "get", "svc",
    "-n", "ingress-nginx",
    "-l", "app.kubernetes.io/component=controller",
    "-o", "jsonpath={.items[0].spec.clusterIP}",
).Output()
```

Also add a `restartCoreDNS` step after patching:

```go
func (r *Runner) restartCoreDNS(ctx context.Context) error {
    return r.kubectl(ctx, "rollout", "restart", "deployment/coredns", "-n", "kube-system")
}
```

#### 3f. Runner needs `credsPath` field

Add `credsPath string` to the `Runner` struct and `NewRunner()` signature so `fetchAndApplyKongKey`
can pass it to the kongkey updater. Update the call site in `cmd/k8s-manager/main.go`.

---

### Step 4 — Export `updateCredsFile` in kongkey package

**File:** `k8s-manager/internal/kongkey/update.go`

Rename `updateCredsFile` → `UpdateCredsFile` and update all call sites.

---

## Setup Flow Summary (visual)

```
Pre-flight checks (cluster reachable, stale vault-init.json warning)
                          │
credentials.yaml  →  values-secrets.yaml (kong key empty)
                          │
       ┌──────────────────┼──────────────────────┐
  nginx ingress      Rancher+cert-mgr        jeeb-infra  (Kong crashes - OK)
                                             jeeb-data   (Keycloak + realm ✓)
                                             jeeb-obs
                                                  │
                                        Init Nexus Docker repo (REST API)
                                                  │
                                        Keycloak ready (HTTP poll)
                                                  │
                                        fetch RS256 key → update creds
                                                  │
                                        re-deploy jeeb-infra → Kong ✓
                                        wait Kong Available
                                                  │
                                   Vault init → unseal → configure
                                                  │
                              CoreDNS patch (nginx ClusterIP) → rollout wait
                                                  │
                                   Verify DNS — all 10 .local domains
                                                  │
                                        Jenkins seed job
                                                  │
                               ╔══════ MANUAL STEPS ══════╗
                               ║  Run Jenkins pipelines   ║
                               ║  (images → Nexus)        ║
                               ║  k8s-manager deploy      ║
                               ║    app learning          ║
                               ╚══════════════════════════╝
```

---

## Verification Steps (after `k8s-manager setup`)

```bash
# 1. All pods running (no CrashLoopBackOff)
kubectl get pods -A

# 2. Keycloak realm exists
curl http://localhost:30081/realms/jeeb | jq .realm

# 3. Kong serves traffic (JWT rejected — correct, no token yet)
curl http://localhost:30088/api/health   # expect 401, not 502

# 4. Vault is unsealed
kubectl exec -n jeeb-infra vault-0 -- vault status

# 5. DNS resolves inside cluster
kubectl run -it --rm dns-test --image=busybox --restart=Never -- \
  nslookup auth.jeeb-dev.local

# 6. Jenkins seed job ran
curl -u admin:<password> http://localhost:30082/job/seed/lastBuild/consoleText
```

---

## What Is NOT Changed

- Jenkins seeding logic (`jenkins/seed.go`) — works correctly as-is
- Vault configuration steps — work correctly as-is
- All Helm chart templates except the Keycloak ConfigMap fix
- `deploy`, `kong-key`, `validate` commands — unchanged
- `values-secrets.yaml` format and content — unchanged

---

## Step 5 — Post-setup Health Check Script (Python)

**New file:** `k8s-manager/scripts/health_check.py`

A standalone Python 3 script (no external deps beyond `requests`) that runs after setup and
produces a pass/fail report. Also callable as `k8s-manager check` command (new Cobra command
that `exec`s the script via `python3`).

### Checks to implement

```
[1] Pod health         — all pods in jeeb-dev / jeeb-infra / jeeb-obs are Running/Completed,
                         no CrashLoopBackOff or ImagePullBackOff
[2] Keycloak endpoint  — GET http://localhost:30081/realms/jeeb → 200, realm == "jeeb"
[3] Vault endpoint     — GET http://localhost:30091/v1/sys/health → 200, sealed == false
[4] Kong endpoint      — GET http://localhost:30088/health → 200 (backend health proxied)
[5] Jenkins endpoint   — GET http://localhost:30082/login → 200
[6] DNS resolution     — kubectl run busybox, nslookup auth.jeeb-dev.local,
                         jenkins.jeeb.local, grafana.jeeb.local; verify each resolves
[7] Vault secrets      — kubectl exec vault-0 vault kv get secret/jeeb/backend/develop → exit 0
```

### Output format

```
=== Jeeb Cluster Health Check ===

[PASS] Pod health         all 14 pods healthy across 3 namespaces
[PASS] Keycloak           http://localhost:30081/realms/jeeb → 200
[FAIL] Vault              http://localhost:30091/v1/sys/health → sealed=true
[PASS] Kong               http://localhost:30088/health → 200
[PASS] Jenkins            http://localhost:30082/login → 200
[FAIL] DNS: auth          nslookup auth.jeeb-dev.local — no address found
[PASS] Vault secrets      backend secrets readable

2 checks failed. Run `k8s-manager maintain` for diagnosis.
Exit code: 1 (if any FAIL), 0 (all pass)
```

### Script structure

```python
# health_check.py
import subprocess, sys, json
import urllib.request   # stdlib only — no pip install needed

checks = []   # list of (name, fn) tuples
results = []  # list of (name, ok, detail)

def check_pods(): ...
def check_http(name, url, assert_fn): ...
def check_dns(hostname): ...
def check_vault_secrets(): ...

if __name__ == "__main__":
    run_all()
    sys.exit(0 if all_passed else 1)
```

---

## Step 6 — `maintain` command (diagnosis report)

**New Cobra command** `k8s-manager maintain` in `cmd/k8s-manager/main.go`.

It calls the health check, collects failures, then for each failing check prints a diagnosis
block with the exact kubectl/helm commands to resolve it — **report only, no auto-execution**.

### Example output for failing checks

```
=== Cluster Maintenance Report ===

[FAIL] Vault is sealed
  Diagnosis : Vault pod restarted and needs unsealing
  Fix       :
    kubectl exec -n jeeb-infra vault-0 -- vault operator unseal <key1>
    kubectl exec -n jeeb-infra vault-0 -- vault operator unseal <key2>
    kubectl exec -n jeeb-infra vault-0 -- vault operator unseal <key3>
  Keys stored in: k8s-manager/vault-init.json

[FAIL] DNS: auth.jeeb-dev.local not resolving
  Diagnosis : CoreDNS patch not applied or nginx ingress IP changed
  Fix       :
    k8s-manager setup --steps coredns   # re-run just that step
    # or manually:
    kubectl apply -f k8s/coredns-patch.yaml
    kubectl rollout restart deployment/coredns -n kube-system
```

The maintain command is implemented in a new package:
**`k8s-manager/internal/maintain/report.go`**

It maps each check name → a static diagnosis + fix template, substitutes known values
(vault-init.json path, namespace names), and prints the report.

---

## Step 7 — Claude skill: `/check`

**New file:** `.claude/skills/check/`

Skill that runs `k8s-manager check` (or `python3 k8s-manager/scripts/health_check.py`)
and formats the output for the conversation. Triggers on: "check cluster", "health check",
"is cluster ok", `/check`.

---

## Updated Verification Steps

```bash
# After setup completes, run the health check:
python3 k8s-manager/scripts/health_check.py
# or via CLI:
k8s-manager check

# If failures found, view diagnosis:
k8s-manager maintain
```
