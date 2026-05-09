# Frontend Troubleshooting

## App loads but API calls fail

- Main frontend uses `VITE_API_URL`, not `/api`, by default.
- Check `frontend/env/.env.local` and confirm `http://localhost:30080` is reachable.
- Learning frontend defaults to `http://localhost:30088/learning`.
- If the browser shows `ERR_CERT_AUTHORITY_INVALID` for `https://*.jeeb-dev.local`, run `cd k8s-manager; go run ./cmd/k8s-manager trust-cert` on Windows.

## Login loop or blank screen after login

- Confirm Keycloak is reachable on the configured URL.
- Check the client ID and redirect URLs in the realm export.
- For local cluster access, verify `.local` DNS or use the exposed NodePort directly.

## Changes to frontend env values do not take effect in Kubernetes

- Check that `/app/env/.env.develop` exists in the frontend pod.
- Check that `/usr/share/nginx/html/app-config.js` was regenerated from that file.
- If values are stale, restart the frontend pod or rollout the deployment so the startup script re-renders runtime config.

## Goals, Events, or Settings do not persist

That is current behavior. Those pages are local-state only.
