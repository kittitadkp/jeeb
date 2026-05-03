# Keycloak Auth Troubleshooting

## Invalid token

**Error:**
```
401 Unauthorized: invalid token
```

**Cause:**
Token expired or malformed.

**Solution:**
- Refresh token before expiry
- Check token in jwt.io
- Verify `iss` claim matches Keycloak URL

---

## CORS on token endpoint

**Error:**
```
CORS error on /token
```

**Solution:**
- Add frontend URL to Keycloak client "Web Origins"
- Use `+` to allow all redirect URIs origins

---

## Redirect URI mismatch

**Error:**
```
Invalid redirect_uri
```

**Solution:**
- Add exact redirect URI in Keycloak client settings
- Include port number
- Check http vs https

---

## Client not found

**Error:**
```
Client not found: xxx
```

**Solution:**
- Verify client ID in Keycloak admin
- Check realm is correct
- Client may be disabled

---

## Token verification failed

**Error:**
```
token signature verification failed
```

**Cause:**
Wrong public key or issuer.

**Solution:**
- Fetch latest JWKS from Keycloak
- Check `KEYCLOAK_URL` and `KEYCLOAK_REALM` env vars
- Verify token algorithm matches

---

## User not authorized

**Error:**
```
403 Forbidden
```

**Cause:**
Missing role or permission.

**Solution:**
- Check user roles in Keycloak
- Verify role mapping in client
- Check backend role validation logic

---

## Keycloak 503 — `keycloak-secret` / `mongo-secret` missing from cluster

**Symptoms:**
- Keycloak returns 503 on all requests
- `kubectl get pods -n jeeb-dev` shows `CreateContainerConfigError` on keycloak and mongodb pods
- `kubectl describe pod -n jeeb-dev <keycloak-pod>` shows: `secret "keycloak-secret" not found`

**Root cause:**
`secrets.yaml` was added to the Helm chart but `apply.sh` was never re-run, so the secrets were never rendered into the cluster. Manual recovery via `helm upgrade --reuse-values` was blocked by a secondary issue: Jenkins had previously stored `global.imageTag: <git-sha>` as a user value in the Helm release history. When Helm merged that stale tag via `--reuse-values`, it produced a YAML parse error on `deployment.yaml` and aborted before creating any resources.

**Fix:**
```bash
# 1. Upgrade using explicit value files (NOT --reuse-values) to bypass stale stored values
helm upgrade jeeb-dev k8s/charts/jeeb-app \
  -n jeeb-dev \
  -f k8s/charts/jeeb-app/values.yaml \
  -f k8s/charts/jeeb-app/values-dev.yaml

# 2. Cycle stuck pods so Kubernetes retries secret mounting
kubectl delete pod -n jeeb-dev -l app=keycloak
kubectl delete pod -n jeeb-dev mongodb-0
```

**Prevention:**
- Always run `bash k8s/apply.sh` after adding new templates to the chart.
- In Jenkins, never use `--reuse-values`. Instead pass explicit files plus `--set`:
  ```bash
  helm upgrade --install jeeb-dev k8s/charts/jeeb-app \
    -n jeeb-dev \
    -f k8s/charts/jeeb-app/values.yaml \
    -f k8s/charts/jeeb-app/values-dev.yaml \
    --set global.imageTag=$GIT_SHA
  ```
  This prevents stale value accumulation in Helm release history.

---

## Keycloak container unhealthy

**Error:**
```
keycloak exited with code 1
```

**Solution:**
```bash
docker-compose logs keycloak
# Common: DB connection, memory, port conflict
# Try: increase memory limit in docker-compose
```
