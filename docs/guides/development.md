# Development Guide

## Prerequisites

- Docker Desktop with Kubernetes enabled
- kubectl
- Go 1.22+
- Node.js 20+

## Start All Services

```bash
# Apply all k8s manifests in order
bash k8s/apply.sh
```

Or apply individually:

```bash
kubectl apply -f k8s/00-namespace.yaml
kubectl apply -f k8s/app/secrets.yaml
kubectl apply -f k8s/app/mongodb/
kubectl apply -f k8s/app/keycloak/
kubectl apply -f k8s/app/backend/
kubectl apply -f k8s/app/frontend/
```

## Services

| Service | URL |
|---------|-----|
| Frontend | http://localhost:30000 |
| Backend API | http://localhost:30080 |
| Keycloak | http://localhost:30081 |
| Jenkins | http://localhost:30082 |
| Nexus UI | http://localhost:30083 |
| SonarQube | http://localhost:30090 |

## Useful kubectl Commands

```bash
# List all pods
kubectl get pods -n jeeb

# Watch pod status
kubectl get pods -n jeeb -w

# Logs
kubectl logs -n jeeb deployment/backend
kubectl logs -n jeeb deployment/frontend

# Restart a deployment
kubectl rollout restart deployment/backend -n jeeb

# Describe a pod (for events/errors)
kubectl describe pod -n jeeb <pod-name>

# Shell into a pod
kubectl exec -it -n jeeb deployment/backend -- sh
```

## Adding a Feature

1. Define domain model in `backend/internal/domain/`
2. Create use case in `backend/internal/usecase/`
3. Define port interface in `backend/internal/port/`
4. Implement adapter in `backend/internal/adapter/`
5. Register route in `backend/internal/adapter/in/http/router.go`
6. Add React hook in `frontend/src/hooks/`
7. Build UI in `frontend/src/pages/` or `frontend/src/components/`

## Testing

```bash
# Backend
cd backend
go test ./...

# Frontend
cd frontend
npm test
```

## Code Style

- Go: `gofmt`, `golangci-lint`
- React: ESLint, Prettier
