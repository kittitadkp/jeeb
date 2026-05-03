# New Cluster Setup

Full procedure for bootstrapping the jeeb stack on a fresh Docker Desktop Kubernetes cluster.
Use either the automated CLI or the manual steps below.

---

## Prerequisites

- Docker Desktop with Kubernetes enabled
- `kubectl` configured and pointing at the cluster (`kubectl cluster-info`)
- `helm` v3 installed
- Go 1.22+ (only for `k8s-manager`)

---

## Option A — Automated (k8s-manager CLI)

### 1. Fill credentials

```bash
cd k8s-manager
cp env/credentials.env.example env/credentials.env
# open env/credentials.env and fill every value
```

| Variable | Description |
|---|---|
| `JENKINS_ADMIN_PASSWORD` | Jenkins admin UI password |
| `JENKINS_GITHUB_USER` | GitHub username for SCM polling |
| `JENKINS_GITHUB_PAT` | GitHub personal access token |
| `JENKINS_NEXUS_USER` | Nexus user Jenkins pushes images as |
| `JENKINS_NEXUS_PAT` | Nexus password for that user |
| `JENKINS_SONAR_TOKEN` | SonarQube analysis token |
| `KEYCLOAK_ADMIN_USER` | Keycloak admin username (e.g. `admin`) |
| `KEYCLOAK_ADMIN_PASSWORD` | Keycloak admin password |
| `MONGODB_USERNAME` | MongoDB app user (e.g. `jeeb`) |
| `MONGODB_PASSWORD` | MongoDB app user password |
| `NEXUS_ADMIN_PASSWORD` | Nexus admin password |
| `SONARQUBE_ADMIN_PASSWORD` | SonarQube admin password |
| `KONG_KEYCLOAK_PUBLIC_KEY` | RS256 public key from Keycloak (fill after Step 8) |

### 2. Dry-run to preview all commands

```bash
go run ./cmd/k8s-manager setup \
  --charts-dir ../k8s/charts \
  --output-dir ./env \
  --dry-run
```

### 3. Run setup

```bash
go run ./cmd/k8s-manager setup \
  --charts-dir ../k8s/charts \
  --output-dir ./env
```

`vault-init.json` is saved to `--output-dir`. It contains the root token and unseal keys — keep it safe and never commit it (`env/.gitignore` covers it).

### 4. Update Kong public key (after Keycloak is running)

Export the realm's RS256 public key:

```
http://localhost:30081/realms/jeeb/protocol/openid-connect/certs
```

Copy the `n` value, convert to PEM format, and set it in `env/credentials.env`:

```
KONG_KEYCLOAK_PUBLIC_KEY=-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----
```

Then update the Helm chart and redeploy:

```bash
go run ./cmd/k8s-manager setup --charts-dir ../k8s/charts --output-dir ./env
```

---

## Option B — Manual Steps

### Step 1 — Deploy all Helm charts

```bash
# jeeb-infra — pass GitHub PAT at install time
helm upgrade --install jeeb-infra k8s/charts/jeeb-infra \
  --namespace jeeb-infra --create-namespace \
  --set jenkins.credentials.githubPat=<token>

# jeeb-app
helm upgrade --install jeeb-dev k8s/charts/jeeb-app \
  --namespace jeeb-dev --create-namespace \
  -f k8s/charts/jeeb-app/values-dev.yaml

# jeeb-obs
helm upgrade --install jeeb-obs k8s/charts/jeeb-obs \
  --namespace jeeb-obs --create-namespace
```

Or run all three at once:

```bash
bash k8s/apply.sh
```

### Step 2 — Initialize Vault

```bash
kubectl exec -n jeeb-infra vault-0 -- vault operator init -format=json > vault-init.json
```

Keep `vault-init.json` — it holds the root token and unseal keys.

### Step 3 — Store unseal keys (enables auto-unseal on pod restart)

```bash
bash k8s/vault/store-unseal-keys.sh vault-init.json
```

### Step 4 — Unseal Vault (first time only)

```bash
# run three times with different keys from vault-init.json
kubectl exec -n jeeb-infra vault-0 -- vault operator unseal <key1>
kubectl exec -n jeeb-infra vault-0 -- vault operator unseal <key2>
kubectl exec -n jeeb-infra vault-0 -- vault operator unseal <key3>
```

After storing the unseal keys secret (Step 3) and restarting the pod, Vault unseals itself automatically on every subsequent restart.

### Step 5 — Configure Vault

```bash
VAULT_TOKEN=<root_token_from_vault-init.json> bash k8s/vault/setup-vault.sh
```

This creates:
- KV v2 secrets engine at `secret/`
- Secrets for `backend`, `frontend`, and `learning` under `secret/jeeb/<service>/develop`
- Kubernetes auth method + roles for each service

### Step 6 — Patch CoreDNS for `.local` DNS

The IP in `k8s/coredns-patch.yaml` is the ClusterIP of the Nginx ingress controller and **changes on every new cluster**.

```bash
# get the new ClusterIP
kubectl get svc -n jeeb-dev -l app.kubernetes.io/name=ingress-nginx \
  -o jsonpath='{.items[0].spec.clusterIP}'
```

Update that IP in `k8s/coredns-patch.yaml`, then apply:

```bash
kubectl apply -f k8s/coredns-patch.yaml
```

### Step 7 — Set Kong public key

After Keycloak is running, get the RS256 public key from:

```
http://localhost:30081/realms/jeeb/protocol/openid-connect/certs
```

Update `kong.keycloakPublicKey` in `k8s/charts/jeeb-infra/values.yaml`, then re-run `apply.sh`.

---

## Step order summary

| # | Step | One-time? |
|---|------|-----------|
| 1 | Deploy Helm charts | No — run on every deploy |
| 2 | `vault operator init` | **Yes** |
| 3 | Store unseal keys secret | **Yes** |
| 4 | Unseal Vault manually | **Yes** (auto after this) |
| 5 | `setup-vault.sh` | **Yes** |
| 6 | CoreDNS patch (update IP first) | **Yes** |
| 7 | Kong public key from Keycloak | **Yes** |

---

## Access table

| Service | NodePort | Namespace |
|---|---|---|
| Frontend | http://localhost:30000 | jeeb-dev |
| Backend | http://localhost:30080 | jeeb-dev |
| Keycloak | http://localhost:30081 | jeeb-dev |
| MongoDB | localhost:30017 | jeeb-dev |
| Learning | http://localhost:30086 | jeeb-dev |
| Jenkins | http://localhost:30082 | jeeb-infra |
| Nexus (UI) | http://localhost:30083 | jeeb-infra |
| Nexus (registry) | localhost:30050 | jeeb-infra |
| Kong | http://localhost:30088 | jeeb-infra |
| SonarQube | http://localhost:30090 | jeeb-infra |
| Vault | http://localhost:30091 | jeeb-infra |
| Grafana | http://localhost:30092 | jeeb-obs |
| Prometheus | http://localhost:30093 | jeeb-obs |

---

## k8s-manager CLI reference

```bash
# check pod status across all jeeb namespaces
go run ./cmd/k8s-manager status

# filter by namespace
go run ./cmd/k8s-manager status -n jeeb-dev

# restart a deployment
go run ./cmd/k8s-manager restart backend -n jeeb-dev

# tail logs
go run ./cmd/k8s-manager logs backend
go run ./cmd/k8s-manager logs backend -f   # follow
```

Build a binary for faster execution:

```bash
cd k8s-manager
go build -o k8s-manager.exe ./cmd/k8s-manager
./k8s-manager.exe setup --charts-dir ../k8s/charts --output-dir ./env
```
