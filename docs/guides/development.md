# Development Guide

## Prerequisites

- Go 1.22 for `backend/` and `learning-backend/`
- Go 1.23 for `k8s-manager/`
- Node 20+
- Docker Desktop with Kubernetes enabled for the full stack

## Recommended local workflow

1. Bootstrap the local cluster with `k8s-manager`.
2. Run the service you are changing from its module directory.
3. Point the matching frontend at the expected NodePort or local backend URL.

## Useful commands

```powershell
cd backend; go run ./cmd/api
cd learning-backend; go run ./cmd/api
cd learning-backend; go run ./cmd/seed
cd frontend; npm ci; npm run dev:local
cd learning-frontend; npm ci; npm run dev
cd k8s-manager; go run ./cmd/k8s-manager check
```

## Environment model

- Go services read `env/.env.<GO_ENV>` and default to `env/.env.local`
- Main frontend reads `env/.env.${APP_ENV}` and defaults to `local`
- Learning frontend uses Vite env files and currently points at Kong by default

## Known local pitfalls

- Main backend expects Keycloak on `host.docker.internal:30081` in the checked-in local env.
- Main frontend defaults to `VITE_API_URL=http://localhost:30080`.
- Learning backend and learning frontend use slightly inconsistent checked-in Keycloak client IDs; follow the file you are running unless you are fixing the mismatch intentionally.
