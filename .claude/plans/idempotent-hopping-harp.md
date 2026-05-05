# Plan: k8s-manager — env-driven config + function deduplication

## Context
`setup.go` and `kongkey/update.go` have ~20 hardcoded cluster values (namespaces, NodePorts, hostnames, Vault paths, Helm release names). Changing cluster topology requires grep-and-replace across multiple files. Several flag definitions and the helm invocation are duplicated across commands. A `jeeb-data` chart also exists in k8s/charts but the k8s-manager doesn't deploy it. This plan moves all cluster constants to a single env-configurable struct, eliminates duplication, and adds the missing `data` deploy target.

---

## Current feature inventory

| Command | What it does |
|---------|-------------|
| `status` | List pods across jeeb namespaces |
| `restart <deployment>` | Rollout restart a deployment |
| `logs <deployment>` | Stream logs from a deployment |
| `setup` | Full bootstrap (charts → Vault init/unseal/configure → CoreDNS) |
| `deploy [infra\|app\|learning\|obs]` | Re-deploy one or more Helm charts |
| `rancher` | Install cert-manager + Rancher (optional) |
| `validate` | Check credentials.yaml completeness |
| `kong-key` | Fetch Keycloak RS256 key and redeploy infra |

### Chart inventory (from k8s/charts/)
| Chart | Namespace | Services |
|-------|-----------|---------|
| `jeeb-data` | jeeb-dev | MongoDB, Keycloak — **new chart, not yet in k8s-manager** |
| `jeeb-app` | jeeb-dev | backend, frontend, learning-frontend |
| `jeeb-learning` | jeeb-dev | learning-backend, learning-frontend |
| `jeeb-infra` | jeeb-infra | Vault, Jenkins, Nexus, SonarQube, Kong |
| `jeeb-obs` | jeeb-obs | Prometheus, Loki, Tempo, Promtail, Grafana |

---

## Problems to fix

### 1. Hardcoded values
| Location | Hardcoded value |
|----------|----------------|
| `setup.go:deployInfra/App/Learning/Obs` | namespaces `jeeb-infra/jeeb-dev/jeeb-obs`; release names |
| `setup.go:waitForVault` | `vault-0`, `jeeb-infra` |
| `setup.go:vaultExec` | `jeeb-infra`, `vault-0`, `http://127.0.0.1:8200` |
| `setup.go:storeUnsealKeysApply` | `jeeb-infra` in inline YAML |
| `setup.go:configureVault` | Vault KV paths, MongoDB URI, Keycloak URLs, ports 30080/30081/30086, realm/client ID |
| `setup.go:patchCoreDNS` | ingress label `app.kubernetes.io/name=ingress-nginx`, `-n jeeb-dev` |
| `kongkey/update.go:fetchFromKeycloak` | `http://localhost:30081`, realm `jeeb` |
| `kongkey/update.go:redeployInfra` | release `jeeb-infra`, namespace `jeeb-infra` |
| `secrets_values.go:nexusDockerConfigJSON` | `localhost:30050` registry |

### 2. Duplicated code
- `--credentials`, `--charts-dir`, `--output-dir`, `--dry-run` flags declared independently in three commands
- `helm()` exec exists in `setup.go` and duplicated in `kongkey/update.go:redeployInfra`
- `kongkey.Updater.redeployInfra` duplicates `Runner.deployInfra` entirely
- Dead code: `setup.go:storeUnsealKeys` (never called — only `storeUnsealKeysApply` runs)

### 3. Missing `jeeb-data` deploy target
`k8s/charts/jeeb-data` (MongoDB + Keycloak) exists but `deploy` and `setup` don't know about it.

---

## Implementation plan

### Step 1 — `internal/config/config.go` (new file)
Define `ClusterConfig` with env-overridable defaults (prefix `K8SM_`, no new library — plain `getEnv(key, default)` helper):

```go
type ClusterConfig struct {
    NamespaceDev    string // K8SM_NAMESPACE_DEV    = jeeb-dev
    NamespaceInfra  string // K8SM_NAMESPACE_INFRA  = jeeb-infra
    NamespaceObs    string // K8SM_NAMESPACE_OBS    = jeeb-obs

    ReleaseData     string // K8SM_RELEASE_DATA     = jeeb-data
    ReleaseDev      string // K8SM_RELEASE_DEV      = jeeb-dev
    ReleaseInfra    string // K8SM_RELEASE_INFRA    = jeeb-infra
    ReleaseLearning string // K8SM_RELEASE_LEARNING = jeeb-learning
    ReleaseObs      string // K8SM_RELEASE_OBS      = jeeb-obs

    VaultPod        string // K8SM_VAULT_POD        = vault-0
    VaultAddr       string // K8SM_VAULT_ADDR       = http://127.0.0.1:8200

    KeycloakRealm    string // K8SM_KEYCLOAK_REALM     = jeeb
    KeycloakClientID string // K8SM_KEYCLOAK_CLIENT_ID = jeeb-app
    KongIssuer       string // K8SM_KONG_ISSUER        = http://auth.jeeb-dev.local/realms/jeeb

    MongoHost        string // K8SM_MONGO_HOST        = mongodb.jeeb-dev.svc.cluster.local:27017
    KeycloakHost     string // K8SM_KEYCLOAK_HOST     = keycloak.jeeb-dev.svc.cluster.local:8080

    KeycloakNodePort int    // K8SM_KEYCLOAK_NODEPORT = 30081
    BackendNodePort  int    // K8SM_BACKEND_NODEPORT  = 30080
    LearningNodePort int    // K8SM_LEARNING_NODEPORT = 30086
    NexusRegistry    string // K8SM_NEXUS_REGISTRY    = localhost:30050

    IngressLabel string // K8SM_INGRESS_LABEL = app.kubernetes.io/name=ingress-nginx

    VaultPathBackend          string // K8SM_VAULT_PATH_BACKEND           = secret/jeeb/backend/develop
    VaultPathFrontend         string // K8SM_VAULT_PATH_FRONTEND          = secret/jeeb/frontend/develop
    VaultPathLearningBackend  string // K8SM_VAULT_PATH_LEARNING_BACKEND  = secret/jeeb/learning/backend/develop
    VaultPathLearningFrontend string // K8SM_VAULT_PATH_LEARNING_FRONTEND = secret/jeeb/learning/frontend/develop
}

func Load() *ClusterConfig
```

**New file:** `k8s-manager/internal/config/config.go`

---

### Step 2 — Thread `ClusterConfig` into `Runner`
- Add `cfg *config.ClusterConfig` to `setup.Runner`
- `setup.NewRunner` gains a `cfg` param
- Replace every hardcoded string in `setup.go` with `r.cfg.*`
- `config.Load()` called once in `main.go`, passed to all command runners

**Modified:** `k8s-manager/internal/setup/setup.go`

---

### Step 3 — Shared helm runner in `internal/helm/runner.go` (new file)
Extract `helm()` into a shared package:
```go
package helm

func Run(ctx context.Context, dryRun bool, args ...string) error
```
- `setup.Runner.helm()` becomes `helm.Run(ctx, r.dryRun, args...)`
- `kongkey.Updater` uses `helm.Run` directly (Step 4 removes its custom reimplementation)

**New file:** `k8s-manager/internal/helm/runner.go`

---

### Step 4 — `kongkey` reuses `setup.Runner.deployInfra`
- Delete `Updater.redeployInfra`
- After updating credentials, construct a `setup.Runner` and call `runner.deployInfra(ctx)`
- Replace hardcoded `http://localhost:30081` with `fmt.Sprintf("http://localhost:%d", u.cfg.KeycloakNodePort)`
- `Updater` gains a `cfg *config.ClusterConfig` field

**Modified:** `k8s-manager/internal/kongkey/update.go`

---

### Step 5 — Add `deploy data` target + setup step
Add `deployData(ctx)` to `Runner` for the `jeeb-data` chart:
```go
func (r *Runner) deployData(ctx context.Context) error {
    return r.helm(ctx, "upgrade", "--install", r.cfg.ReleaseData,
        filepath.Join(r.chartsDir, "jeeb-data"),
        "--namespace", r.cfg.NamespaceDev,
        "--create-namespace",
        "-f", filepath.Join(r.chartsDir, "jeeb-data", "values-dev.yaml"),
        "-f", r.secretsFile,
    )
}
```
- Add `"data"` as a valid target in `newDeployCmd` argument validation
- Insert `deployData` as step 2 in `Runner.Run` (before `deployApp`)

**Modified:** `k8s-manager/internal/setup/setup.go`, `k8s-manager/cmd/k8s-manager/main.go`

---

### Step 6 — Hoist shared flags to root command in `main.go`
Move `--credentials`, `--charts-dir`, `--output-dir`, `--dry-run` to `root.PersistentFlags()`. Remove the three duplicate copies from `newSetupCmd`, `newDeployCmd`, `newKongKeyCmd`.

**Modified:** `k8s-manager/cmd/k8s-manager/main.go`

---

### Step 7 — Remove dead code + typed vault secrets
- Delete `setup.go:storeUnsealKeys` (the dry-run-client version — never called)
- Replace the anonymous `struct{ path, key, val string }` in `configureVault` with a named `vaultKV` type

**Modified:** `k8s-manager/internal/setup/setup.go`

---

## Updated setup flow (after changes)
| Step | Action |
|------|--------|
| 0 | Generate values-secrets.yaml |
| 1 | Deploy `jeeb-infra` (Vault, Jenkins, Nexus, SonarQube, Kong) |
| 2 | Deploy `jeeb-data` (MongoDB, Keycloak) ← **new** |
| 3 | Deploy `jeeb-app` (backend, frontend) |
| 4 | Deploy `jeeb-learning` |
| 5 | Deploy `jeeb-obs` |
| 6 | Wait for Vault pod |
| 7 | Init Vault → vault-init.json |
| 8 | Store unseal keys in K8s secret |
| 9 | Unseal Vault |
| 10 | Configure Vault (KV, policies, K8s auth) |
| 11 | Patch CoreDNS |

---

## Files changed

| File | Change |
|------|--------|
| `internal/config/config.go` | **new** — `ClusterConfig` + `Load()` |
| `internal/helm/runner.go` | **new** — shared `helm.Run()` |
| `internal/setup/setup.go` | replace hardcoded strings; add `cfg`; add `deployData`; remove `storeUnsealKeys`; typed `vaultKV` |
| `internal/kongkey/update.go` | remove `redeployInfra`; use `cfg` for URLs; use `setup.Runner` |
| `internal/setup/secrets_values.go` | accept `nexusRegistry` param instead of hardcoded `localhost:30050` |
| `cmd/k8s-manager/main.go` | hoist shared flags to root; call `config.Load()` once; add `"data"` deploy target |

---

## Verification
```bash
# dry-run full setup — should print commands using configured values
K8SM_NAMESPACE_DEV=custom-dev go run ./cmd/k8s-manager setup --dry-run

# override Keycloak port for kong-key
K8SM_KEYCLOAK_NODEPORT=30099 go run ./cmd/k8s-manager kong-key --dry-run

# deploy only data chart
go run ./cmd/k8s-manager deploy data --dry-run

# existing tests
go test ./...
```
