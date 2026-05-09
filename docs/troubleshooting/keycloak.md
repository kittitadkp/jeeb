# Keycloak Troubleshooting

## Realm or issuer mismatch

- Main backend expects `KEYCLOAK_URL/realms/KEYCLOAK_REALM`.
- Kong defaults to issuer `http://auth.jeeb-dev.local/realms/jeeb`.
- If signing keys changed, run:

```powershell
cd k8s-manager
go run ./cmd/k8s-manager kong-key
```

## Client ID confusion

The committed realm export contains `jeeb-app` and `learning-app`, but checked-in local env files are not fully aligned:

- `frontend/env/.env.local` uses `jeeb-app`
- `learning-frontend/env/.env.local` uses `jeeb-app`
- `learning-backend/env/.env.local` uses `jeeb-client`

Treat this as an implementation inconsistency, not a documentation typo.

## Keycloak is up but auth still fails

- Check `.local` DNS resolution in the cluster.
- Confirm the backend can reach the issuer URL it was configured with.
- Verify the browser client redirect URLs in `realm-jeeb.json`.
