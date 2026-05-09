# Deployment Guide

## Supported path

The maintained deployment path is Kubernetes with Helm, Jenkins, Nexus, Vault, and Keycloak. There is no current root Docker Compose deployment flow.

## Release flow

1. Jenkins checks out the service repository and the top-level `k8s` repository.
2. Optional test and SonarQube stages run when `SKIP_SONAR=false`.
3. Kaniko builds and pushes images to the Nexus registry.
4. On `main`, Jenkins runs `helm upgrade --install` and waits for rollout.

## Branch behavior

- `main`: build, push, and deploy
- `develop`: build and push only
- other branches: no image publish or deployment through the shared pipeline

## Chart mapping

| Service | Helm chart | Values key |
|---|---|---|
| `backend` | `jeeb-app` | `backend.imageTag` |
| `frontend` | `jeeb-app` | `frontend.imageTag` |
| `learning-backend` | `jeeb-learning` | `backend.imageTag` |
| `learning-frontend` | `jeeb-learning` | `frontend.imageTag` |

## Secrets and config

- Backends read Vault-rendered env files during container startup.
- Frontends currently do not consume the mounted Vault env files.
- Keycloak realm configuration is shipped in `k8s/charts/jeeb-data/files/realm-jeeb.json`.

## Manual rollout

After pipelines have pushed images:

```powershell
cd k8s-manager
go run ./cmd/k8s-manager deploy app learning
go run ./cmd/k8s-manager check
```
