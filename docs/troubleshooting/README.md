# Troubleshooting

Use this section when the current implementation does not match expected local or cluster behavior.

## Quick triage

```powershell
cd k8s-manager
go run ./cmd/k8s-manager check
go run ./cmd/k8s-manager maintain
```

## Focus areas

- [backend.md](backend.md)
- [frontend.md](frontend.md)
- [keycloak.md](keycloak.md)
- [mongodb.md](mongodb.md)
- [docker.md](docker.md)

## Common facts

- There is no maintained root Docker Compose workflow.
- Main backend metrics live at `/metrics`; learning-backend has no matching endpoint.
- Frontend runtime config is mostly build-time config today.
