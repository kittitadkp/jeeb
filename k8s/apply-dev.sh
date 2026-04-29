#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "==> Deploying jeeb-app (dev)"
helm upgrade --install jeeb-dev "$SCRIPT_DIR/charts/jeeb-app" \
  --namespace jeeb-dev --create-namespace \
  -f "$SCRIPT_DIR/charts/jeeb-app/values-dev.yaml" \
  "$@"

echo ""
echo "Dev deployed. Access:"
echo "  http://jeeb-dev.local"
echo "  http://api.jeeb-dev.local"
echo "  http://auth.jeeb-dev.local"
