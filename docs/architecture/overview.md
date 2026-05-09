# Architecture Overview

## Workspace shape

Jeeb is a multi-app workspace with two user-facing products:

- Personal tracker: `frontend/` + `backend/`
- Learning app: `learning-frontend/` + `learning-backend/`

Infrastructure, deployment, and operations live in `k8s/`, `k8s-manager/`, and `jenkins/`.

## Service boundaries

### Main backend

- Clean Architecture layout under `internal/domain`, `internal/usecase`, `internal/port`, and `internal/adapter`
- Stores users, workouts, studies, sleep records, finance transactions, events, and master data in MongoDB
- Verifies Keycloak access tokens directly
- Exposes `/health` and `/metrics`

### Learning backend

- Similar layered Go layout, but narrower domain
- Stores topics, items, progress, and users in MongoDB
- Can verify Keycloak tokens itself or trust upstream verification from Kong with `UPSTREAM_AUTH=true`
- Exposes `/health`, but not `/metrics`

### Frontends

- Both frontends are Vite + React + TypeScript apps with Keycloak login on load
- Main frontend talks directly to `VITE_API_URL`
- Learning frontend defaults to Kong at `http://localhost:30088/learning`

## Request flow

1. User signs in through Keycloak.
2. Frontend stores the access token in memory.
3. API calls send `Authorization: Bearer <token>`.
4. Backend resolves or creates the user record from JWT claims.
5. Use cases operate on MongoDB repositories and return JSON responses.

## Deployment architecture

- Helm charts deploy `jeeb-data`, `jeeb-infra`, `jeeb-app`, `jeeb-learning`, and `jeeb-obs`
- Vault renders backend env files before process startup
- Jenkins builds images with Kaniko, pushes to Nexus, and deploys Helm releases on `main`
- Grafana, Loki, Tempo, and Prometheus run in `jeeb-obs`

## Current implementation caveats

- Main frontend pages `Goals`, `Events`, and `Settings` are local UI only.
- Event calendar sync is exposed in the main backend, but no calendar provider is wired.
- Frontend deployments mount Vault-rendered env files, but the shipped Nginx images do not read them. Frontend config is effectively decided at build time.
