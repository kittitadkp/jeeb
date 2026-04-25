# Development Guide

## Setup

```bash
# Clone and start
git clone <repo>
docker-compose up --build
```

## Services

| Service | Port | URL |
|---------|------|-----|
| Frontend | 3000 | http://localhost:3000 |
| API | 8080 | http://localhost:8080 |
| MongoDB | 27017 | mongodb://localhost:27017 |
| Keycloak | 8081 | http://localhost:8081 |

## Adding a Feature

1. Define domain in `internal/domain/`
2. Create use case in `internal/usecase/`
3. Define port in `internal/port/`
4. Implement adapter in `internal/adapter/`
5. Register routes in `cmd/`

## Testing

```bash
go test ./internal/...
cd frontend && npm test
```

## Code Style

- Go: `gofmt`, `golangci-lint`
- React: ESLint, Prettier
