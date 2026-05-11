# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Shell preference

Always use the Bash tool for shell commands. Never use PowerShell.

## Repository layout

This is a full monorepo — all application and infrastructure code lives here:

```
jeeb/
  backend/            # Go API (main app — workouts, study, sleep, finance, events)
  frontend/           # React app (main app)
  learning-backend/   # Go API (learning platform — topics, items, progress)
  learning-frontend/  # React app (learning platform)
  jeeb-react-shared/  # Shared React component library (@jeeb/react-shared on Nexus)
  k8s/                # Helm charts + cluster config
  k8s-manager/        # Go CLI for bootstrapping and operating the K8s stack
  jenkins/            # CI/CD pipelines and Groovy job definitions
  docs/               # Architecture, API reference, feature specs, guides
```

## Kubernetes

All dev services run in the `jeeb-dev` namespace on Docker Desktop Kubernetes.

### Cluster management (k8s-manager CLI)

`k8s-manager` is the primary tool for cluster operations — prefer it over running kubectl/helm directly.

```powershell
# First-time setup
Copy-Item k8s-manager/env/secrets.yaml.example k8s-manager/env/secrets.yaml
go run ./cmd/k8s-manager validate

# Full bootstrap
go run ./cmd/k8s-manager setup

# Deploy specific charts
go run ./cmd/k8s-manager deploy app
go run ./cmd/k8s-manager deploy learning
go run ./cmd/k8s-manager deploy infra
go run ./cmd/k8s-manager deploy obs

# Operations
go run ./cmd/k8s-manager check
go run ./cmd/k8s-manager maintain
go run ./cmd/k8s-manager status
go run ./cmd/k8s-manager restart <deployment>
go run ./cmd/k8s-manager logs <deployment>
go run ./cmd/k8s-manager trust-cert
go run ./cmd/k8s-manager kong-key
go run ./cmd/k8s-manager seed
```

### NodePort map

**jeeb-dev namespace**

| Service | NodePort | In-cluster host |
|---------|----------|-----------------|
| frontend | 30000 | `frontend.jeeb-dev.svc.cluster.local:80` |
| backend | 30080 | `backend.jeeb-dev.svc.cluster.local:8080` |
| keycloak | 30081 | `keycloak.jeeb-dev.svc.cluster.local:8080` |
| mongodb | 30017 | `mongodb.jeeb-dev.svc.cluster.local:27017` |
| learning-backend | 30086 | `learning-backend.jeeb-dev.svc.cluster.local:8080` |
| learning-frontend | 30087 | `learning-frontend.jeeb-dev.svc.cluster.local:80` |

**jeeb-infra namespace**

| Service | NodePort | In-cluster host |
|---------|----------|-----------------|
| jenkins | 30082 | `jenkins.jeeb-infra.svc.cluster.local:8080` |
| nexus (ui) | 30083 | `nexus.jeeb-infra.svc.cluster.local:8081` |
| nexus (registry) | 30050 | `nexus.jeeb-infra.svc.cluster.local:5000` |
| kong | 30088 | `kong.jeeb-infra.svc.cluster.local:8000` |
| sonarqube | 30090 | `sonarqube.jeeb-infra.svc.cluster.local:9000` |
| vault | 30091 | `vault.jeeb-infra.svc.cluster.local:8200` |

**jeeb-obs namespace**

| Service | NodePort | In-cluster host |
|---------|----------|-----------------|
| grafana | 30092 | `grafana.jeeb-obs.svc.cluster.local:3000` |
| prometheus | 30093 | `prometheus.jeeb-obs.svc.cluster.local:9090` |

**cattle-system namespace**

| Service | NodePort | In-cluster host |
|---------|----------|-----------------|
| rancher | 30443 | `rancher.cattle-system.svc.cluster.local:443` |

### k8s structure

```
k8s/
  coredns-patch.yaml    # kubectl apply to kube-system for in-cluster .local DNS
  tls-issuer.yaml       # cert-manager ClusterIssuer
  charts/
    jeeb-app/           # frontend, backend (jeeb-dev ns)
    jeeb-data/          # mongodb, keycloak — shared data services (jeeb-dev ns)
    jeeb-learning/      # learning-backend, learning-frontend (jeeb-dev ns)
    jeeb-infra/         # jenkins, nexus, sonarqube, vault, kong (jeeb-infra ns)
    jeeb-obs/           # prometheus, loki, tempo, grafana, promtail (jeeb-obs ns)
```

## Backend (Go — main app)

```bash
cd backend
go test ./...
go build ./cmd/api/...
go run ./cmd/api/main.go   # requires .env
```

**Stack:** Go 1.22, Chi router, go-oidc/v3, mongo-driver, envconfig, validator/v10, slog

**Architecture — Clean/Hexagonal:**
```
internal/
  domain/           # Pure structs + business rules
  usecase/          # Orchestrates domain logic, depends on port interfaces
  port/in/          # UseCase interfaces
  port/out/         # Repository + integration interfaces
  adapter/in/http/  # Chi handlers → use cases
  adapter/out/
    mongo/          # MongoDB repositories
    integration/    # External integrations
  config/           # envconfig structs
cmd/api/            # Entry point + exercise master seeding
```

**Features:** workout, study, sleep, finance, event — each has files at every layer.

**Routes** (all except `/health` require Bearer token from Keycloak):
- `GET /health`, `GET /me`
- `/workouts`, `/study`, `/sleep`, `/finance`, `/events` — CRUD + `/stats`
- `POST /events/:id/sync` — calendar sync

## Learning Backend (Go)

```bash
cd learning-backend
go test ./...
go run ./cmd/api/main.go
```

**Stack:** Same as backend (Go 1.22, Chi, go-oidc/v3, mongo-driver, slog)

**Features:** topic, item, progress, user

**Architecture:** Same Clean/Hexagonal pattern as backend.

**Routes:**
- `GET /health`, `GET /me`
- `/topics`, `/items`, `/progress`

## Frontend (React — main app)

```bash
cd frontend
npm run dev           # dev server → http://localhost:3000
npm run dev:local     # APP_ENV=local
npm run dev:docker    # APP_ENV=docker
npm run build         # tsc + vite build
npm run lint:fix
```

**Stack:** React 19, TypeScript 6, Vite, TanStack Query v5, Tailwind CSS v3, Radix UI, keycloak-js v26, react-router-dom v7, Lucide React, `@jeeb/react-shared`

**Structure:**
```
src/
  hooks/      # TanStack Query hooks — one file per feature
  pages/      # Dashboard, Workouts, Study, Sleep, Finance, Calendar, Settings
  components/
  types/      # TypeScript interfaces matching backend snake_case JSON
  lib/        # API client, utilities
  store/      # Global state (theme, etc.)
```

All API calls go through hooks in `src/hooks/`. Types in `src/types/` must match backend JSON field names exactly.

**Design system:** Blue-600 primary, Slate neutrals, Lucide icons (24px, 1.5px stroke), `rounded-lg shadow-sm` cards, fixed 240px left sidebar + fixed header. Import shared components from `@jeeb/react-shared`.

## Learning Frontend (React)

```bash
cd learning-frontend
npm run dev
npm run build
```

**Stack:** React 19, TypeScript 6, Vite, TanStack Query v5, Tailwind CSS v3, keycloak-js v26, react-router-dom v7, `@jeeb/react-shared`

**Pages:** Home (topic list), Topic (item list + progress tracking)

## Shared React Library (jeeb-react-shared)

Published as `@jeeb/react-shared` to the Nexus npm registry. Both frontends consume it.

```bash
cd jeeb-react-shared
npm run build    # tsup build → dist/
npm run dev      # tsup --watch
```

**Exports:**
- `@jeeb/react-shared` — re-exports all below
- `@jeeb/react-shared/ui` — Button, Card, Badge, StatCard, States, SectionLabel
- `@jeeb/react-shared/charts` — shared chart components
- `@jeeb/react-shared/auth` — AuthProvider, useAuth, keycloak setup
- `@jeeb/react-shared/utils` — shared utilities

After changing the library, bump the version, build, and publish to Nexus before frontend changes take effect in the cluster.

## k8s-manager CLI (Go)

```bash
cd k8s-manager
go run ./cmd/k8s-manager <command>
```

**Stack:** Go 1.23, Cobra, k8s client-go, yaml

**Commands:** `setup`, `deploy`, `seed`, `check`, `maintain`, `trust-cert`, `validate`, `kong-key`, `patch-jenkins-creds`, `redeploy-jenkins`, `namespace`, `status`, `restart`, `logs`, `rancher`

Config lives in `k8s-manager/env/secrets.yaml` (gitignored) and `env/config.yaml`.

## CI/CD

Jenkins polls GitHub (`H/5 * * * *`) and runs pipelines in `jenkins/pipelines/`. Pipeline flow: test → SonarQube → Kaniko build → push to Nexus → `kubectl set image`.

```
jenkins/
  pipelines/    # Jenkinsfiles per service
  jobs/         # Groovy seed job definitions
  vars/         # Shared pipeline library steps
  resources/    # Pipeline resource files
```

Images push to `nexus.jeeb-infra.svc.cluster.local:5000/jeeb/<service>` (in-cluster) and pull via `localhost:30050` from outside.

## Wiki (Obsidian knowledge base)

The Obsidian vault is the **project root** (`D:\personal\jeeb`). Wiki pages live in `wiki/` (gitignored — generated content, not source code).

Run `/wiki-update` after:
- Adding a new feature or domain (new route group, new page, new service)
- Introducing a new architectural pattern or abstraction
- Making a significant infrastructure change (new Helm chart, new K8s component)
- Fixing a bug that revealed a non-obvious constraint or lesson learned
- Any decision that would confuse you returning to the codebase in 3 months

Do **not** run `/wiki-update` for routine bug fixes, UI tweaks, or dependency bumps.

## Go conventions

**New feature = all 5 hexagonal layers.** Adding a domain concept requires files at every layer:
`domain/` → `port/in/` → `port/out/` → `usecase/` → `adapter/in/http/handler/` + `adapter/out/mongo/`
Never import across layer boundaries (e.g. handler must not import domain directly, only through usecase interfaces).

**Error handling:** Wrap errors with context using `fmt.Errorf("action: %w", err)`. Never `panic` in domain or usecase layers. Only `log.Fatal` in `cmd/` entry points.

**Tests:** Usecase logic with branching requires table-driven tests. Repositories and handlers don't need unit tests — the running cluster covers integration paths.

**Avoid:** global state, `init()` functions, embedding structs for behavior (use interfaces instead).

## TypeScript / React conventions

**API calls:** All API calls go through hooks in `src/hooks/`. Components never call `api.*` directly.

**Types:** All interfaces in `src/types/index.ts` use `snake_case` field names matching the backend JSON exactly. Never add a frontend type that doesn't have a corresponding backend struct.

**No `any`:** Use `unknown` + narrowing, or define a proper interface. `as any` only for unavoidable third-party interop.

**Mutations:** Body type is always `Omit<Entity, 'id' | 'user_id' | 'created_at' | 'updated_at'>` — the server assigns those fields.

**Design tokens:** Use `C`, `T`, `W`, `R`, `S` from `src/lib/design.ts` for inline styles. Don't scatter raw hex values or pixel numbers directly in components.

## Git commit format

```
type: short description (imperative, lowercase, no period)
```

Types: `feat` `fix` `chore` `refactor` `docs` `test` `infra`

Examples: `feat: add sleep stats endpoint` · `fix: workout duration not saved on update` · `infra: add vault policy for learning-backend`

## Troubleshooting docs

Whenever a bug or infrastructure problem is fixed, update `docs/troubleshooting/` before closing the task:

- Pick the file matching the affected service (`keycloak.md`, `backend.md`, `frontend.md`, `mongodb.md`, `docker.md`) or create a new one.
- Add a section: **Symptoms**, **Root cause**, **Fix** (with exact commands), **Prevention**.
- Do this automatically — do not wait for the user to ask.

## Plan progress tracking

When working from a plan file in `.claude/plans/`, mark each step complete by updating `[ ]` to `[x]` immediately after it succeeds. Mark blocked steps `[-]` with a one-line note.

## Slash commands

| Command | Purpose |
|---------|---------|
| `/backend` | Go backend tasks with full arch context |
| `/frontend` | React tasks with design system context |
| `/k8s` | Kubernetes manifest work |
| `/k8s-manager` | k8s-manager CLI development |
| `/jenkins` | CI/CD pipeline help |
| `/docs` | Write/update documentation |
| `/status` | Check all pod health |
| `/logs <svc>` | Tail service logs |
| `/deploy <svc>` | Restart and watch a deployment |
