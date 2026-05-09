# AI Agent Guide

## Repository shape

- `backend/`: Go backend services. In the main backend, keep `domain -> usecase -> port -> adapter` boundaries intact.
- `frontend/`: frontend application code. Match existing component, page, hook, and state patterns before adding new structure.
- `k8s/`: Kubernetes manifests and Helm-related deployment config.
- `jenkins/`: CI/CD pipelines and delivery automation.
- `docs/`: architecture notes, runbooks, troubleshooting, and process docs.

## Default working rules

- Prefer minimal, safe changes over broad rewrites.
- Read the nearby code, config, tests, and docs before editing.
- Match existing naming, structure, and style unless the task is explicitly a cleanup.
- If code and docs disagree, treat code as source of truth and update docs when needed.
- Do not change secrets, environment-specific values, or deployment credentials unless asked.

## Debugging

- Reproduce the issue first when possible. Capture the exact error, URL, command, log line, or failing request.
- Isolate the problem by layer: client, auth, API, database, network, then infrastructure.
- Prefer evidence from logs, configs, tests, manifests, and runtime behavior over guesses.
- Propose the smallest fix that addresses the verified root cause.

## Refactoring

- Keep refactors scoped to the task.
- Preserve behavior unless the user asked for a behavior change.
- Avoid unnecessary file moves, renames, or interface churn.
- If a cleanup adds risk without clear payoff, do not do it.

## Testing

- Run the smallest relevant checks first, then broader validation if the change warrants it.
- Backend changes: `go test ./...` in the affected Go module.
- Frontend changes: `npm run lint` and `npm run build` in the affected frontend module.
- Infra or pipeline changes: run the relevant validation command if one exists.
- If you cannot run a useful check, state that clearly.

## Documentation

- Update docs when setup steps, behavior, debugging steps, or operator workflows change.
- Keep docs short, actionable, and aligned with the repository as it exists now.
- Put reusable prompts, troubleshooting playbooks, and templates under `docs/`.
