# Testing Guide

## What is automated today

- Go modules support `go test ./...`
- `k8s-manager` currently has the clearest committed unit-test coverage
- Frontends do not expose a `test` script

## Commands

```powershell
cd backend; go test ./...
cd learning-backend; go test ./...
cd k8s-manager; go test ./...
cd frontend; npm run lint; npm run build
cd learning-frontend; npm run lint; npm run build
```

## Expectations for new work

- Add focused `*_test.go` files for new Go logic
- Keep frontend changes buildable and lint-clean
- For cluster changes, run `go run ./cmd/k8s-manager check` after deployment when practical

## Gaps

- No committed end-to-end test suite
- No frontend unit-test harness
- Sparse backend coverage outside the operational tooling
