#!/bin/bash
# Store Vault unseal keys in a Kubernetes Secret for auto-unseal on restart.
# Usage: bash k8s/vault/store-unseal-keys.sh <path-to-vault-init.json>
# vault-init.json is produced by: vault operator init -format=json
set -euo pipefail

INIT_FILE="${1:?Usage: bash k8s/vault/store-unseal-keys.sh <vault-init.json>}"
NS="jeeb"

if ! command -v jq &>/dev/null; then
  echo "Error: jq is required. Install it with: apt install jq / brew install jq" >&2
  exit 1
fi

KEY1=$(jq -r '.unseal_keys_b64[0]' "$INIT_FILE")
KEY2=$(jq -r '.unseal_keys_b64[1]' "$INIT_FILE")
KEY3=$(jq -r '.unseal_keys_b64[2]' "$INIT_FILE")

if [ "$KEY1" = "null" ] || [ -z "$KEY1" ]; then
  echo "Error: could not parse unseal_keys_b64 from $INIT_FILE" >&2
  exit 1
fi

kubectl create secret generic vault-unseal-keys \
  -n "$NS" \
  --from-literal=key1="$KEY1" \
  --from-literal=key2="$KEY2" \
  --from-literal=key3="$KEY3" \
  --dry-run=client -o yaml | kubectl apply -f -

echo "Unseal keys stored in secret vault-unseal-keys (namespace: $NS)"
echo "Restart the Vault pod to activate auto-unseal:"
echo "  kubectl rollout restart statefulset/vault -n $NS"
