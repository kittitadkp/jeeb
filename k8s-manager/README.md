# k8s-manager

CLI for bootstrapping and managing the jeeb Kubernetes cluster on Docker Desktop.

## Prerequisites

- Docker Desktop with Kubernetes enabled
- `kubectl` and `helm` on PATH
- Go 1.23+

## Setup

```powershell
# 1. Copy and fill credentials
cp env/secrets.yaml.example env/secrets.yaml
# edit env/secrets.yaml — fill all required fields

# 2. (Optional) Override cluster topology
cp env/config.yaml.example env/config.yaml
# edit env/config.yaml — only needed if ports/namespaces differ from defaults

# 3. Validate before running
go run ./cmd/k8s-manager validate
```

## First-time cluster setup

```powershell
# Reset cluster first (Docker Desktop → Settings → Kubernetes → Reset Kubernetes Cluster)

# Bootstrap everything (~15-20 min)
go run ./cmd/k8s-manager setup
```

The `setup` command runs 12 steps in order:

| # | Step |
|---|------|
| 1 | Deploy `jeeb-infra` — Vault, Jenkins, Nexus, SonarQube, Kong |
| 2 | Deploy `jeeb-data` — MongoDB, Keycloak |
| 3 | Deploy `jeeb-app` — backend, frontend |
| 4 | Deploy `jeeb-learning` |
| 5 | Deploy `jeeb-obs` — Prometheus, Loki, Grafana |
| 6 | Wait for Vault pod ready |
| 7 | Initialize Vault → saves `vault-init.json` to `--output-dir` |
| 8 | Store unseal keys → Kubernetes secret `vault-unseal-keys` |
| 9 | Unseal Vault |
| 10 | Configure Vault (KV engine, policies, K8s auth roles) |
| 11 | Patch CoreDNS (detects ingress ClusterIP automatically) |
| 12 | Seed Jenkins (create seed job, generate all pipeline jobs) |

## After setup

```powershell
# Trigger Jenkins pipelines via UI at http://localhost:30082
# Run in order: backend → frontend → learning-backend → learning-frontend
# Each pipeline: test → SonarQube → build image → push to Nexus → kubectl set image

# Deploy app charts once images are in Nexus
go run ./cmd/k8s-manager deploy app learning

# Verify cluster health
go run ./cmd/k8s-manager check
```

## Commands

```
setup                                      Bootstrap a new cluster end-to-end
deploy [infra|data|app|learning|obs]       Re-deploy one or more charts (no Vault init)
seed                                       Create/run Jenkins seed job
kong-key                                   Fetch Keycloak RS256 key and update Kong
check                                      Health check — pass/fail table for all services
maintain                                   Diagnose failures and print fix commands
validate                                   Check all secrets.yaml fields are filled
namespace                                  Create jeeb namespaces
status [-n <ns>]                           Show pod status
restart <deployment>                       Restart a deployment
logs <deployment>                          Stream logs from a deployment
rancher                                    Install cert-manager + Rancher (optional)
patch-jenkins-creds                        Patch jenkins-secret from secrets.yaml and restart Jenkins
redeploy-jenkins                           Rollout-restart Jenkins and wait until healthy
```

### Command details

**`deploy`** — re-run `helm upgrade --install` for selected charts; defaults to all five if no target given.

```powershell
go run ./cmd/k8s-manager deploy              # all charts
go run ./cmd/k8s-manager deploy app          # jeeb-app only
go run ./cmd/k8s-manager deploy infra obs    # jeeb-infra + jeeb-obs
```

**`logs`** — stream logs from a deployment.

```powershell
go run ./cmd/k8s-manager logs backend               # last 100 lines
go run ./cmd/k8s-manager logs backend -f            # follow (stream continuously)
go run ./cmd/k8s-manager logs backend --tail 500    # show last 500 lines
go run ./cmd/k8s-manager logs backend -n jeeb-infra # different namespace
```

**`seed`** — create and run the Job DSL seed job in Jenkins.

```powershell
go run ./cmd/k8s-manager seed
go run ./cmd/k8s-manager seed --groovy-path path/to/seed.groovy
go run ./cmd/k8s-manager seed --jenkins-url http://localhost:30082
```

**`patch-jenkins-creds`** — updates `jenkins-secret` with current values from `secrets.yaml` and rolls Jenkins to pick them up. Patches all six keys: `admin-password`, `github-user`, `github-pat`, `nexus-user`, `nexus-password`, `sonar-token`.

**`redeploy-jenkins`** — restarts the Jenkins deployment and waits for it to be ready (rollout status → `/login` poll).

```powershell
go run ./cmd/k8s-manager redeploy-jenkins
go run ./cmd/k8s-manager redeploy-jenkins --timeout 10m
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--secrets` | `env/secrets.yaml` | Path to secrets file |
| `--config` | `env/config.yaml` | Path to cluster topology overrides (optional) |
| `--charts-dir` | `../k8s/charts` | Path to Helm charts directory |
| `--output-dir` | `.` (cwd) | Directory for `vault-init.json` and `values-secrets.yaml` |
| `--dry-run` | false | Print commands without executing |
| `--log-level` | `info` | Verbosity: `debug`, `info`, `warn`, `error` |
| `--kubeconfig` | `~/.kube/config` | Path to kubeconfig |

### Per-command flags

| Command | Flag | Default | Description |
|---------|------|---------|-------------|
| `status`, `restart`, `logs` | `-n`, `--namespace` | `jeeb-dev` | Target namespace |
| `logs` | `-f`, `--follow` | false | Stream logs continuously |
| `logs` | `--tail` | 100 | Number of recent lines to show |
| `seed` | `--groovy-path` | auto-detected | Path to `seed.groovy` |
| `seed` | `--jenkins-url` | `http://localhost:<nodeport>` | Jenkins URL |
| `patch-jenkins-creds` | `-n`, `--namespace` | `jeeb-infra` | Namespace where Jenkins is deployed |
| `redeploy-jenkins` | `-n`, `--namespace` | `jeeb-infra` | Namespace where Jenkins is deployed |
| `redeploy-jenkins` | `--timeout` | `5m` | Rollout + health-check timeout |

## Service endpoints

| Service | NodePort | Namespace |
|---------|----------|-----------|
| Frontend | 30000 | jeeb-dev |
| Backend | 30080 | jeeb-dev |
| Keycloak | 30081 | jeeb-dev |
| Jenkins | 30082 | jeeb-infra |
| Nexus (UI) | 30083 | jeeb-infra |
| Nexus (registry) | 30050 | jeeb-infra |
| Learning backend | 30086 | jeeb-dev |
| Learning frontend | 30087 | jeeb-dev |
| Kong | 30088 | jeeb-infra |
| SonarQube | 30090 | jeeb-infra |
| Vault | 30091 | jeeb-infra |
| Grafana | 30092 | jeeb-obs |
| Prometheus | 30093 | jeeb-obs |
| Rancher | 30443 | cattle-system |

## File layout

```
k8s-manager/
  cmd/k8s-manager/main.go     CLI entry point, flags, commands
  internal/
    config/                   Cluster topology (namespaces, ports, hosts)
    credentials/              Secrets loader (secrets.yaml parser)
    helm/                     helm upgrade --install wrapper
    healthcheck/              HTTP, DNS, Vault, pod health checks
    jenkins/                  Jenkins seed job orchestrator + credentials patcher
    kongkey/                  Keycloak RS256 key fetcher/updater
    kube/                     Kubernetes client (namespaces, pods, logs)
    logger/                   zerolog wrapper (Step, Info, Debug, Warn, Error)
    maintain/                 Diagnosis report with fix commands
    printer/                  Pod status table formatter
    rancher/                  cert-manager + Rancher deployer
    redeploy/                 Jenkins rollout-restart + health-check waiter
    setup/                    12-step bootstrap orchestrator
    util/                     Shared: exec, HTTP polling, JWK→PEM
    validate/                 Credential completeness checker
  env/
    secrets.yaml              Operator credentials (gitignored)
    secrets.yaml.example      Template — copy and fill
    config.yaml               Cluster topology overrides (gitignored)
    config.yaml.example       Template — all fields optional
  vault-init.json             Vault unseal keys + root token (gitignored — keep safe)
  values-secrets.yaml         Generated Helm values (gitignored)
```

## Troubleshooting

```powershell
# See what's failing
go run ./cmd/k8s-manager check

# Get diagnosis + fix commands for each failure
go run ./cmd/k8s-manager maintain

# Re-run setup after fixing (safe to re-run — idempotent)
go run ./cmd/k8s-manager setup

# Verbose output
go run ./cmd/k8s-manager --log-level debug setup

# Refresh Kong key after Keycloak restart
go run ./cmd/k8s-manager kong-key

# Fix Jenkins credentials after secrets.yaml change
go run ./cmd/k8s-manager patch-jenkins-creds
```

> **vault-init.json** contains unseal keys and root token. Keep it safe — losing it means Vault cannot be unsealed after a pod restart.
