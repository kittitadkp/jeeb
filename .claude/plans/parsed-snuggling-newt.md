# Plan: Setup Rancher

## Context

Adding Rancher to the local Docker Desktop Kubernetes cluster to provide a GUI for cluster management. Rancher requires cert-manager as a prerequisite. The project already has a pattern of separate `apply-*.sh` scripts + `k8s/charts/` Helm chart directories — this plan follows that exact pattern.

**NodePort 30443** is available (all ports 30000–30093 and 30500 are taken).

---

## Files to Create / Modify

| Action | File |
|--------|------|
| Create | `k8s/charts/jeeb-rancher/Chart.yaml` |
| Create | `k8s/charts/jeeb-rancher/values.yaml` |
| Create | `k8s/apply-rancher.sh` |
| Update | `CLAUDE.md` — add Rancher row to NodePort table |

---

## Steps

- [x] **1. Create `k8s/charts/jeeb-rancher/Chart.yaml`**

```yaml
apiVersion: v2
name: jeeb-rancher
description: Rancher cluster management GUI — cert-manager + Rancher (cattle-system)
type: application
version: 0.1.0
```

- [x] **2. Create `k8s/charts/jeeb-rancher/values.yaml`**

```yaml
certManager:
  version: v1.15.3
  namespace: cert-manager

rancher:
  version: 2.9.3
  namespace: cattle-system
  hostname: rancher.jeeb-infra.local
  tlsSource: rancher      # self-signed via cert-manager
  replicas: 1
  nodePort: 30443

ingress:
  className: nginx
```

- [x] **3. Create `k8s/apply-rancher.sh`**

```bash
#!/bin/bash
set -e

CERT_MANAGER_VERSION="v1.15.3"
RANCHER_VERSION="2.9.3"
RANCHER_HOSTNAME="rancher.jeeb-infra.local"
RANCHER_NAMESPACE="cattle-system"
CERT_NS="cert-manager"
RANCHER_NP="30443"

echo "==> Adding / updating Helm repos"
helm repo add jetstack https://charts.jetstack.io --force-update
helm repo add rancher-stable https://releases.rancher.com/server-charts/stable
helm repo update

echo "==> Installing cert-manager ${CERT_MANAGER_VERSION}"
helm upgrade --install cert-manager jetstack/cert-manager \
  --namespace "${CERT_NS}" \
  --create-namespace \
  --version "${CERT_MANAGER_VERSION}" \
  --set crds.enabled=true

echo "==> Waiting for cert-manager pods to be ready (up to 120s)..."
kubectl rollout status deployment/cert-manager           -n "${CERT_NS}" --timeout=120s
kubectl rollout status deployment/cert-manager-webhook   -n "${CERT_NS}" --timeout=120s
kubectl rollout status deployment/cert-manager-cainjector -n "${CERT_NS}" --timeout=120s

echo "==> Installing Rancher ${RANCHER_VERSION}"
helm upgrade --install rancher rancher-stable/rancher \
  --namespace "${RANCHER_NAMESPACE}" \
  --create-namespace \
  --version "${RANCHER_VERSION}" \
  --set hostname="${RANCHER_HOSTNAME}" \
  --set ingress.tls.source=rancher \
  --set ingress.ingressClassName=nginx \
  --set replicas=1 \
  --set bootstrapPassword=admin

echo "==> Waiting for Rancher rollout (up to 5 minutes)..."
kubectl rollout status deployment/rancher -n "${RANCHER_NAMESPACE}" --timeout=300s

echo "==> Patching rancher service to NodePort ${RANCHER_NP}"
kubectl patch svc rancher -n "${RANCHER_NAMESPACE}" \
  --type='json' \
  -p="[
    {\"op\": \"replace\", \"path\": \"/spec/type\", \"value\": \"NodePort\"},
    {\"op\": \"add\", \"path\": \"/spec/ports/0/nodePort\", \"value\": ${RANCHER_NP}}
  ]"

echo ""
echo "Rancher deployed. Access:"
echo "  Via NodePort https://localhost:${RANCHER_NP}  (accept self-signed cert warning)"
echo "  Via ingress  https://${RANCHER_HOSTNAME}      (requires hosts file entry)"
echo ""
echo "  Bootstrap password: admin  (you will be prompted to change it on first login)"
echo ""
echo "  Add to C:\\Windows\\System32\\drivers\\etc\\hosts (as Administrator):"
echo "    127.0.0.1  ${RANCHER_HOSTNAME}"
echo ""
echo "  kubectl get pods -n ${RANCHER_NAMESPACE}   to check status"
echo "  kubectl get pods -n ${CERT_NS}             to check cert-manager"
```

- [x] **4. Update `CLAUDE.md` NodePort map**

Append one row to the existing table (line 48, after vault):

```
| rancher | 30443 | `rancher.cattle-system.svc.cluster.local:443` |
```

---

## Verification

```bash
# Run the script
bash k8s/apply-rancher.sh

# Check cert-manager (should see 3 Running pods)
kubectl get pods -n cert-manager

# Check Rancher (should see rancher-* pod Running)
kubectl get pods -n cattle-system

# Confirm NodePort is set
kubectl get svc rancher -n cattle-system

# Access UI
# Browser → https://localhost:30443
# Accept self-signed cert → enter bootstrap password "admin" → set permanent password
```

---

## Notes

- Docker Desktop should have **≥ 6 GB RAM** allocated (Settings → Resources) — Rancher + cert-manager adds ~600 MB.
- The `--set crds.enabled=true` flag is the modern cert-manager CRD install method (≥ v1.15). Do not use the old `kubectl apply -f crds.yaml` approach.
- The NodePort patch changes the service type to `NodePort` and assigns port 30443. Re-running `apply-rancher.sh` is safe — `helm upgrade --install` and `helm repo add --force-update` are both idempotent. The `kubectl patch` may emit a warning on re-run if the port is already set, but will not fail.
- Rancher version: always run `helm search repo rancher-stable/rancher` first to confirm the latest stable version, then update `values.yaml` and the script variable accordingly.
