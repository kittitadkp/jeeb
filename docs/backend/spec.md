# Backend Specification

## Main backend

### Module layout

```text
cmd/api                 Entrypoint
internal/domain         Core entities
internal/usecase        Application logic
internal/port           Input and output contracts
internal/adapter        HTTP and MongoDB adapters
pkg/                    Shared app packages
env/                    Local env files
```

### Collections

- `users`
- `workouts`
- `studies`
- `sleep`
- `finance`
- `events`
- `master`

### Runtime behavior

- Loads `env/.env.<GO_ENV>` or `env/.env.local`
- Connects to MongoDB, builds repositories, and wires use cases
- Seeds `master` exercise records when the collection is empty
- Starts OIDC verification against `KEYCLOAK_URL/realms/<realm>`
- Enables OTLP tracing only when `OTEL_EXPORTER_OTLP_ENDPOINT` is set

## Learning backend

### Collections

- `users`
- `topics`
- `items`
- `progress`

### Runtime behavior

- Loads `env/.env.<GO_ENV>` or `env/.env.local`
- Supports direct OIDC validation or Kong-verified upstream auth
- Does not seed data on startup
- Uses a separate CLI seed path: `go run ./cmd/seed`

## Validation and errors

- Request validation is handled at the HTTP layer before use cases run
- Known error codes: `UNAUTHORIZED`, `FORBIDDEN`, `NOT_FOUND`, `VALIDATION_ERROR`, `CONFLICT`, `INTERNAL_ERROR`
- Unhandled errors are returned as `500` with `INTERNAL_ERROR`
