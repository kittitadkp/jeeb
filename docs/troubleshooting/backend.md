# Backend Troubleshooting

## Main backend does not start

- Check `backend/env/.env.local` or `GO_ENV` selection.
- Confirm MongoDB on `localhost:30017`.
- Confirm Keycloak on `http://host.docker.internal:30081`.

## `401 unauthorized`

- Verify the frontend is sending a bearer token.
- Check that the token issuer matches the configured Keycloak realm.
- For the learning backend, confirm whether `UPSTREAM_AUTH` should be `true` or `false`.

## `503 calendar integration not configured`

This is expected in the current main backend. `POST /events/{id}/sync` is exposed, but the service is started with `NewEventUseCase(eventRepo, nil)`.

## Metrics confusion

- Main backend: `GET /metrics` exists.
- Learning backend: deployment has scrape annotations, but the service does not expose `/metrics`.
