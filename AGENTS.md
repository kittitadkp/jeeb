# Repository Guidelines

## Project Structure & Module Organization
`backend/` and `learning-backend/` are Go services organized around `cmd/`, `internal/`, `pkg/`, and `env/`. In the main backend, keep Clean Architecture boundaries intact: `internal/domain`, `internal/usecase`, `internal/port`, and `internal/adapter`. `frontend/` and `learning-frontend/` are Vite + React apps; most work lands in `src/components`, `src/pages`, `src/hooks`, `src/lib`, and `src/store`. `k8s/` contains Helm charts and cluster manifests. `k8s-manager/` is the operational Go CLI for cluster bootstrap and maintenance. CI/CD definitions live under `jenkins/pipelines/`, and design, API, and troubleshooting notes belong in `docs/`.

## Build, Test, and Development Commands
Run commands from the relevant module directory:

- `cd frontend && npm run dev` starts the main UI; `npm run build` produces a production bundle; `npm run lint` runs ESLint and Prettier-backed checks.
- `cd learning-frontend && npm run dev|build|lint` does the same for the learning app.
- `cd backend && go test ./...` validates Go packages; use the same command in `learning-backend/` and `k8s-manager/`.
- `cd k8s-manager && go run ./cmd/k8s-manager validate` checks `env/secrets.yaml`; `go run ./cmd/k8s-manager setup` bootstraps the cluster; `go run ./cmd/k8s-manager check` runs health checks.

## Coding Style & Naming Conventions
Go code should stay `gofmt`-clean, use tabs, and keep package names lowercase. Prefer small, focused packages and wire new backend features through domain -> usecase -> port -> adapter. React/TypeScript uses 2-space indentation, semicolons, and double quotes as formatted today. Use `PascalCase` for components and pages (`Dashboard.tsx`), `camelCase` for hooks and utilities, and keep shared UI in `src/components`.

## Testing Guidelines
Add Go tests as `*_test.go` beside the package under test and run `go test ./...` before opening a PR. Current automated coverage is light; new backend or CLI logic should include focused unit tests. Frontend packages do not currently expose a `test` script, so UI changes must at least pass `npm run lint` and `npm run build`.

## Commit & Pull Request Guidelines
Recent history shows short, fix-focused subjects such as `fix: keycloak`, `fix realm`, and `fix: deployment`. Follow that pattern with concise, imperative summaries, but avoid placeholder subjects like `1` or `2`. PRs should list affected modules, validation commands run, linked issues/tasks, and screenshots for frontend, Kubernetes, or Jenkins changes.

## Security & Configuration Tips
Do not commit real secrets from `*/env/`, `k8s-manager/env/secrets.yaml`, `values-secrets.yaml`, or `vault-init.json`. Use the example files, keep credentials local, and document any new required variables in `docs/` and the owning module README.
