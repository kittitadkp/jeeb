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
kubectl rollout status deployment/cert-manager            -n "${CERT_NS}" --timeout=120s
kubectl rollout status deployment/cert-manager-webhook    -n "${CERT_NS}" --timeout=120s
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
