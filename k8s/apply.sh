#!/bin/bash
set -e

echo "==> Creating namespace"
kubectl apply -f k8s/00-namespace.yaml

echo "==> Deploying Jenkins"
kubectl apply -f k8s/jenkins/

echo "==> Deploying SonarQube"
kubectl apply -f k8s/sonarqube/

echo "==> Deploying Nexus"
kubectl apply -f k8s/nexus/

echo "==> Deploying App (MongoDB, Keycloak, Backend, Frontend)"
kubectl apply -f k8s/app/secrets.yaml
kubectl apply -f k8s/app/mongodb/
kubectl apply -f k8s/app/keycloak/
kubectl apply -f k8s/app/backend/
kubectl apply -f k8s/app/frontend/

echo ""
echo "All resources applied. Access:"
echo "  Jenkins    http://localhost:30082"
echo "  SonarQube  http://localhost:30090"
echo "  Nexus      http://localhost:30083"
echo "  Backend    http://localhost:30080"
echo "  Frontend   http://localhost:30000"
echo "  Keycloak   http://localhost:30081"
