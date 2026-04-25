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
