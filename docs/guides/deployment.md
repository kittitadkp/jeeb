# Deployment Guide

## Prerequisites

- Docker Desktop with Kubernetes enabled
- kubectl configured (`kubectl cluster-info`)
- Nexus running at `localhost:30050` (Docker registry)

## Apply All Manifests

```bash
bash k8s/apply.sh
```

## Environment Variables (Backend)

Managed via `k8s/app/backend/configmap.yaml` and `k8s/app/secrets.yaml`.

| Variable | Source | Description |
|----------|--------|-------------|
| `MONGO_URI` | Secret `mongo-secret` | MongoDB connection string |
| `KEYCLOAK_URL` | ConfigMap | Keycloak base URL |
| `KEYCLOAK_REALM` | ConfigMap | Auth realm |
| `KEYCLOAK_CLIENT_ID` | ConfigMap | OAuth client ID |

## CI/CD (Jenkins)

Pipelines are defined in `jenkins/backend/Jenkinsfile` and `jenkins/frontend/Jenkinsfile`.

On every push to `main`:
1. Jenkins polls SCM and detects the change
2. Kubernetes pod agent spins up with the required containers
3. Tests run, SonarQube analysis runs
4. Kaniko builds and pushes image to Nexus (`localhost:30050/jeeb/<service>:latest`)
5. `kubectl set image` updates the deployment
6. `kubectl rollout status` waits for rollout to complete

## Manual Image Build & Deploy

```bash
# Build and push image to Nexus
docker build -t localhost:30050/jeeb/backend:latest ./backend
docker push localhost:30050/jeeb/backend:latest

# Update deployment
kubectl set image deployment/backend backend=localhost:30050/jeeb/backend:latest -n jeeb
kubectl rollout status deployment/backend -n jeeb
```

## Health Checks

```bash
# Backend
curl http://localhost:30080/health

# Check pod readiness
kubectl get pods -n jeeb

# Describe deployment for events
kubectl describe deployment backend -n jeeb
```

## Secrets

Secrets are defined in `k8s/app/secrets.yaml` (base64-encoded). Never commit real values — use placeholders and apply manually:

```bash
kubectl apply -f k8s/app/secrets.yaml
```

## Nexus Docker Registry

Jenkins uses the in-cluster registry at `nexus.jeeb.svc.cluster.local:5000`.
From your local machine use `localhost:30050`.

```bash
# Login
docker login localhost:30050 -u admin

# Pull an image
docker pull localhost:30050/jeeb/backend:latest
```
