#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "==> Deploying jeeb-obs (Prometheus, Loki, Tempo, Grafana)"
helm upgrade --install jeeb-obs "$SCRIPT_DIR/charts/jeeb-obs" \
  --namespace jeeb-obs --create-namespace

echo ""
echo "Access:"
echo "  Grafana    http://grafana.jeeb.local  (NodePort: http://localhost:30092)"
echo "  Prometheus http://localhost:30093"
echo ""
echo "  kubectl get pods -n jeeb-obs   to check status"
