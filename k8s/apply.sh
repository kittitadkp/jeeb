#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "==> Deploying jeeb-infra (Vault, Jenkins, Nexus, SonarQube, Kong)"
helm upgrade --install jeeb-infra "$SCRIPT_DIR/charts/jeeb-infra" \
  --namespace jeeb-infra --create-namespace

echo "==> Deploying jeeb-app (dev)"
helm upgrade --install jeeb-dev "$SCRIPT_DIR/charts/jeeb-app" \
  --namespace jeeb-dev --create-namespace \
  -f "$SCRIPT_DIR/charts/jeeb-app/values-dev.yaml"

echo "==> Deploying jeeb-obs (Prometheus, Loki, Tempo, Grafana)"
helm upgrade --install jeeb-obs "$SCRIPT_DIR/charts/jeeb-obs" \
  --namespace jeeb-obs --create-namespace

echo ""
echo "All releases applied. Access:"
echo ""
echo "  jeeb-infra:"
echo "    Jenkins    http://jenkins.jeeb-infra.local    (NodePort: http://localhost:30082)"
echo "    SonarQube  http://sonarqube.jeeb-infra.local  (NodePort: http://localhost:30090)"
echo "    Nexus      http://nexus.jeeb-infra.local      (NodePort: http://localhost:30083)"
echo "    Vault      http://vault.jeeb-infra.local      (NodePort: http://localhost:30091)"
echo "    Kong       http://localhost:30088"
echo ""
echo "  jeeb-dev:"
echo "    App        http://jeeb-dev.local"
echo "    API        http://api.jeeb-dev.local"
echo "    Auth       http://auth.jeeb-dev.local"
echo "    Learning   http://localhost:30086"
echo ""
echo "  jeeb-obs:"
echo "    Grafana    http://grafana.jeeb.local  (NodePort: http://localhost:30092)"
echo "    Prometheus http://localhost:30093"
echo ""
echo "  helm list -A   to see all releases"
