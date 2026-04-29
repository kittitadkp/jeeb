# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository layout

This is a monorepo containing infrastructure only. Application code lives in separate folders (excluded from git):

```
jeeb/
  k8s/          # All Kubernetes manifests — the primary concern of this repo
  docs/         # Architecture, API reference, feature specs, guides
  backend/      # Go API (not committed here — has its own repo)
  frontend/     # React app (not committed here — has its own repo)
  jenkins/      # CI/CD pipelines and setup scripts (not committed here)
```

## Kubernetes

All services run in the `jeeb` namespace on Docker Desktop Kubernetes.

```bash
# Apply everything
bash k8s/apply.sh

# Apply a single service
kubectl apply -f k8s/app/backend/

# Common operations
kubectl get pods -n jeeb
kubectl logs -n jeeb deployment/backend
kubectl rollout restart deployment/backend -n jeeb
kubectl exec -it -n jeeb deployment/backend -- sh
```

### NodePort map

| Service | NodePort | In-cluster host |
|---------|----------|-----------------|
| frontend | 30000 | `frontend.jeeb.svc.cluster.local:80` |
| backend | 30080 | `backend.jeeb.svc.cluster.local:8080` |
| keycloak | 30081 | `keycloak.jeeb.svc.cluster.local:8080` |
| jenkins | 30082 | `jenkins.jeeb.svc.cluster.local:8080` |
| nexus (ui) | 30083 | `nexus.jeeb.svc.cluster.local:8081` |
| sonarqube | 30090 | `sonarqube.jeeb.svc.cluster.local:9000` |
| mongodb | 30017 | `mongodb.jeeb.svc.cluster.local:27017` |
| nexus (registry) | 30050 | `nexus.jeeb.svc.cluster.local:5000` |
| vault | 30091 | `vault.jeeb.svc.cluster.local:8200` |

### k8s structure

```
k8s/
  00-namespace.yaml
  apply.sh
  app/
    secrets.yaml              # mongo-secret, keycloak-secret (base64)
    backend/                  # configmap + deployment + service
    frontend/                 # deployment + service
    keycloak/                 # deployment + service
    mongodb/                  # statefulset + service
  jenkins/                    # deployment + pvc + rbac + service
  nexus/                      # deployment + pvc + service
  sonarqube/                  # deployment + pvc + service
  vault/                      # statefulset + pvc + configmap + rbac + service
```

Backend env vars come from `k8s/app/backend/configmap.yaml` (non-secret) and `mongo-secret` (MONGO_URI).

## Backend (Go)

```bash
cd backend
go test ./...                        # all tests
go test ./internal/usecase/...       # single package
go build ./cmd/api/...               # build binary
go run ./cmd/api/main.go             # run locally (requires .env)
```

**Stack:** Go 1.22, Chi router, go-oidc/v3, mongo-driver, envconfig, validator/v10

**Architecture — Clean/Hexagonal:**
```
internal/
  domain/           # Pure structs + business rules, no imports
  usecase/          # Orchestrates domain logic, depends on port interfaces
  port/in/          # UseCase interfaces (called by handlers)
  port/out/         # Repository + integration interfaces (implemented by adapters)
  adapter/in/http/  # Chi handlers → call use cases via port/in
  adapter/out/      # MongoDB repos + external integrations
  config/           # envconfig structs loaded from environment
```

Each feature (workout, study, sleep, finance, event) has its own files at every layer. Adding a feature means creating one file per layer and wiring it in `cmd/api/main.go`.

All domain structs use `json:"snake_case"` tags. Handlers use `middleware.RespondJSON` for all responses.

**Routes** (all except `/health` require Bearer token from Keycloak):
- `GET /health`, `GET /me`
- `/workouts`, `/study`, `/sleep`, `/finance`, `/events` — all have CRUD + `/stats`
- `POST /events/:id/sync` — calendar sync

## Frontend (React)

```bash
cd frontend
npm run dev       # dev server → http://localhost:3000
npm run build     # tsc + vite build
npm run lint      # eslint
npm run lint:fix  # eslint --fix
```

**Stack:** React 19, TypeScript, Vite, TanStack Query v5, Tailwind CSS, Radix UI, keycloak-js, react-router-dom v7, Lucide React

**Structure:**
```
src/
  hooks/    # TanStack Query hooks — one file per feature (useWorkouts, useStudy, etc.)
  pages/    # Full-page components (Dashboard, Workouts, Study, Sleep, Finance, Calendar, Settings)
  components/
  types/    # TypeScript interfaces matching backend snake_case JSON
  lib/      # API client, utilities
  store/    # Global state
```

All API calls go through hooks in `src/hooks/` — never fetch directly in components. Types in `src/types/` must match backend JSON field names exactly.

**Design system:** Blue-600 primary, Slate neutrals, Lucide icons only (24px, 1.5px stroke), `rounded-lg shadow-sm` cards, fixed 240px left sidebar + fixed header layout.

## CI/CD

Jenkins polls GitHub (`H/5 * * * *`) and runs pipelines defined in `jenkins/backend/Jenkinsfile` and `jenkins/frontend/Jenkinsfile`. Each pipeline: test → SonarQube → Kaniko build → push to Nexus → `kubectl set image`.

Jenkins is fully bootstrapped via the `jeeb-infra` Helm chart. Plugins, credentials, and jobs are configured declaratively — no manual setup step needed:
```bash
# Fill in credentials in values.yaml (or pass via --set), then:
bash k8s/apply.sh
```

Pipelines use in-cluster service URLs. Images are pushed to `nexus.jeeb.svc.cluster.local:5000/jeeb/<service>` and pulled via `localhost:30050` from outside the cluster.

## Plan progress tracking

When working from a plan file in `.claude/plans/`, mark each step complete as you finish it by updating `[ ]` to `[x]` in the file. Do this immediately after each step succeeds — not at the end. This lets the plan survive context resets and session restarts.

If a step is partially done or blocked, mark it `[-]` and add a one-line note explaining why.

## Slash commands

| Command | Purpose |
|---------|---------|
| `/backend` | Go backend tasks with full arch context |
| `/frontend` | React tasks with design system context |
| `/k8s` | Kubernetes manifest work |
| `/jenkins` | CI/CD pipeline help |
| `/docs` | Write/update documentation |
| `/status` | Check all pod health |
| `/logs <svc>` | Tail service logs |
| `/deploy <svc>` | Restart and watch a deployment |
