#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "==> Deploying jeeb-app (uat)"
helm upgrade --install jeeb-uat "$SCRIPT_DIR/charts/jeeb-app" \
  --namespace jeeb-uat --create-namespace \
  -f "$SCRIPT_DIR/charts/jeeb-app/values-uat.yaml" \
  "$@"

echo ""
echo "UAT deployed. Access:"
echo "  http://jeeb-uat.local"
echo "  http://api.jeeb-uat.local"
echo "  http://auth.jeeb-uat.local"
