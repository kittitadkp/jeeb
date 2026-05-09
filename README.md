# Jeeb

Jeeb is a personal tracking workspace split across four runnable apps and a Kubernetes-first operations stack:

- `backend/`: main Go API for workouts, study, sleep, finance, calendar events, and seeded exercise master data
- `frontend/`: main React app for the personal tracker UI
- `learning-backend/`: Go API for topic, item, and learning-progress management
- `learning-frontend/`: React app for the learning experience

This repository is the umbrella workspace. Several subdirectories are also separate Git repositories and are referenced that way by Jenkins.

## Repository layout

```text
backend/            Main API (Go 1.22, MongoDB, Keycloak)
frontend/           Main UI (Vite, React 19, TypeScript)
learning-backend/   Learning API (Go 1.22, MongoDB, Keycloak)
learning-frontend/  Learning UI (Vite, React 19, TypeScript)
k8s/                Helm charts, realm export, DNS/TLS manifests
k8s-manager/        Cluster bootstrap and maintenance CLI
jenkins/            Jenkins shared library, seed job, pipeline definitions
docs/               Architecture, API, ops, feature, and troubleshooting docs
```

## Local development

There is no root `docker-compose.yml`. The current workflow is either:

1. Bootstrap the local Kubernetes stack with `k8s-manager`, then run apps against its NodePorts.
2. Run services module-by-module and provide your own MongoDB and Keycloak instances that match the checked-in env files.

Typical commands:

```powershell
cd backend; go run ./cmd/api
cd frontend; npm ci; npm run dev:local
cd learning-backend; go run ./cmd/api
cd learning-frontend; npm ci; npm run dev
```

Validation commands:

```powershell
cd backend; go test ./...
cd learning-backend; go test ./...
cd k8s-manager; go test ./...
cd frontend; npm run lint; npm run build
cd learning-frontend; npm run lint; npm run build
```

## Cluster workflow

The supported deployment path is Kubernetes on Docker Desktop:

```powershell
cd k8s-manager
Copy-Item env/secrets.yaml.example env/secrets.yaml
go run ./cmd/k8s-manager validate
go run ./cmd/k8s-manager setup
go run ./cmd/k8s-manager check
```

After bootstrap, Jenkins builds and publishes images to Nexus. App workloads are then rolled out with:

```powershell
go run ./cmd/k8s-manager deploy app learning
```

## Current implementation notes

- Main backend exposes `/metrics`; learning-backend does not, even though its deployment has Prometheus scrape annotations.
- Main backend wires event sync without a calendar provider, so `POST /events/{id}/sync` currently returns `503`.
- `frontend/` uses `VITE_API_URL` directly. The Nginx `/api` proxy exists, but the app does not rely on it by default.
- `Goals`, `Events`, and `Settings` pages in the main frontend are local-state UI only; they are not backed by the API.

## Documentation map

- [docs/README.md](docs/README.md)
- [backend/README.md](backend/README.md)
- [frontend/README.md](frontend/README.md)
- [learning-backend/README.md](learning-backend/README.md)
- [learning-frontend/README.md](learning-frontend/README.md)
- [k8s-manager/README.md](k8s-manager/README.md)
