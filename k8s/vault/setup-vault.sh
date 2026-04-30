#!/bin/bash
# One-time Vault configuration for jeeb services (backend, frontend, learning).
# Usage: VAULT_TOKEN=<root-token> bash k8s/vault/setup-vault.sh
set -euo pipefail

ROOT_TOKEN="${VAULT_TOKEN:?Set VAULT_TOKEN to your Vault root token (from vault-init.json)}"
NS="jeeb-infra"
POD="vault-0"

# Run vault CLI inside the pod so it has cluster-internal access
v() {
  kubectl exec -i -n "$NS" "$POD" -- \
    env VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN="$ROOT_TOKEN" \
    vault "$@"
}

echo "==> Enabling KV v2 secrets engine at secret/"
v secrets enable -path=secret kv-v2 2>/dev/null || echo "   already enabled"

echo "==> Writing backend develop secrets"
v kv put secret/jeeb/backend/develop \
  PORT=8080 \
  LOG_LEVEL=INFO \
  READ_TIMEOUT=10 \
  WRITE_TIMEOUT=10 \
  MONGO_DATABASE=jeeb \
  "MONGO_URI=mongodb://jeeb:jeeb123@mongodb.jeeb-dev.svc.cluster.local:27017/jeeb?authSource=admin" \
  "KEYCLOAK_URL=http://keycloak.jeeb-dev.svc.cluster.local:8080" \
  KEYCLOAK_REALM=jeeb \
  KEYCLOAK_CLIENT_ID=jeeb-app

echo "==> Writing frontend develop secrets"
v kv put secret/jeeb/frontend/develop \
  "VITE_KEYCLOAK_URL=http://localhost:30081" \
  VITE_KEYCLOAK_REALM=jeeb \
  VITE_KEYCLOAK_CLIENT_ID=jeeb-app \
  "VITE_API_URL=http://localhost:30080"

echo "==> Writing learning develop secrets"
v kv put secret/jeeb/learning/develop \
  PORT=8080 \
  LOG_LEVEL=INFO \
  MONGO_DATABASE=jeeb_learning \
  "MONGO_URI=mongodb://jeeb:jeeb123@mongodb.jeeb-dev.svc.cluster.local:27017/jeeb_learning?authSource=admin" \
  "KEYCLOAK_URL=http://keycloak.jeeb-dev.svc.cluster.local:8080" \
  KEYCLOAK_REALM=jeeb \
  KEYCLOAK_CLIENT_ID=jeeb-app

echo "==> Enabling Kubernetes auth method"
v auth enable kubernetes 2>/dev/null || echo "   already enabled"

echo "==> Configuring Kubernetes auth (uses in-cluster API server)"
v write auth/kubernetes/config \
  kubernetes_host="https://kubernetes.default.svc:443"

echo "==> Writing backend policy"
v policy write backend-policy - << 'POLICY'
path "secret/data/jeeb/backend/develop" {
  capabilities = ["read"]
}
POLICY

echo "==> Creating Kubernetes auth role for backend"
v write auth/kubernetes/role/backend \
  bound_service_account_names=backend \
  bound_service_account_namespaces=jeeb-dev \
  policies=backend-policy \
  ttl=1h

echo "==> Writing frontend policy"
v policy write frontend-policy - << 'POLICY'
path "secret/data/jeeb/frontend/develop" {
  capabilities = ["read"]
}
POLICY

echo "==> Creating Kubernetes auth role for frontend"
v write auth/kubernetes/role/frontend \
  bound_service_account_names=frontend \
  bound_service_account_namespaces=jeeb-dev \
  policies=frontend-policy \
  ttl=1h

echo "==> Writing learning policy"
v policy write learning-policy - << 'POLICY'
path "secret/data/jeeb/learning/develop" {
  capabilities = ["read"]
}
POLICY

echo "==> Creating Kubernetes auth role for learning"
v write auth/kubernetes/role/learning \
  bound_service_account_names=learning \
  bound_service_account_namespaces=jeeb-dev \
  policies=learning-policy \
  ttl=1h

echo ""
echo "Done. Vault is configured for backend, frontend, and learning."
echo ""
echo "Next steps:"
echo "  bash k8s/apply.sh          to deploy all services"
echo "  kubectl get pods -n jeeb-dev   to check status"
