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

The `setup` command runs 20 steps in order:

| # | Step |
|---|------|
| 1 | Pre-flight checks (cluster reachable, helm on PATH, clear stale Kong key) |
| 2 | Remove stale files from previous run (`vault-init.json`, `values-secrets.yaml`) |
| 3 | Generate `values-secrets.yaml` from `env/secrets.yaml` |
| 4 | Install nginx ingress controller |
| 5 | Install Rancher + cert-manager (skipped if already installed) |
| 6 | Deploy `jeeb-data` — MongoDB, Keycloak |
| 7 | Wait for Keycloak ready |
| 8 | Fetch Kong RS256 key from Keycloak JWKS |
| 9 | Deploy `jeeb-infra` — Vault, Jenkins, Nexus, SonarQube, Kong |
| 10 | Wait for Kong ready |
| 11 | Wait for Vault pod ready |
| 12 | Initialize Vault → saves `vault-init.json` |
| 13 | Store unseal keys in Kubernetes secret |
| 14 | Unseal Vault |
| 15 | Configure Vault (KV engine, policies, K8s auth roles) |
| 16 | Initialize Nexus Docker registry |
| 17 | Patch CoreDNS for `.local` DNS |
| 18 | Wait for CoreDNS rollout |
| 19 | Verify DNS for all `.local` domains |
| 20 | Seed Jenkins (create seed job, generate pipeline jobs) |

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
setup                     Bootstrap a new cluster end-to-end
deploy [infra|data|app|learning|obs]  Re-deploy one or more charts
seed                      Create/run Jenkins seed job
kong-key                  Fetch Keycloak RS256 key and update Kong
check                     Health check — pass/fail table for all services
maintain                  Diagnose failures and print fix commands
validate                  Check all secrets.yaml fields are filled
namespace                 Create jeeb namespaces
status [-n <ns>]          Show pod status
restart <deployment>      Restart a deployment
logs <deployment>         Stream logs from a deployment
rancher                   Install cert-manager + Rancher (optional)
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
    jenkins/                  Jenkins seed job orchestrator
    kongkey/                  Keycloak RS256 key fetcher/updater
    kube/                     Kubernetes client (namespaces, pods, logs)
    logger/                   zerolog wrapper (Step, Info, Debug, Warn, Error)
    maintain/                 Diagnosis report with fix commands
    printer/                  Pod status table formatter
    rancher/                  cert-manager + Rancher deployer
    setup/                    20-step bootstrap orchestrator
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
```

> **vault-init.json** contains unseal keys and root token. Keep it safe — losing it means Vault cannot be unsealed after a pod restart.
