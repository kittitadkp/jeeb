# Frontend Troubleshooting

## App loads but API calls fail

- Main frontend uses `VITE_API_URL`, not `/api`, by default.
- Check `frontend/env/.env.local` and confirm `http://localhost:30080` is reachable.
- Learning frontend defaults to `http://localhost:30088/learning`.

## Login loop or blank screen after login

- Confirm Keycloak is reachable on the configured URL.
- Check the client ID and redirect URLs in the realm export.
- For local cluster access, verify `.local` DNS or use the exposed NodePort directly.

## Changes to frontend env values do not take effect in Kubernetes

That is current behavior. The deployed Nginx images serve a static build and do not read the Vault-rendered env files mounted into the pods.

## Goals, Events, or Settings do not persist

That is current behavior. Those pages are local-state only.
