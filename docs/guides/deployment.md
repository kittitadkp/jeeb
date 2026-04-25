# Deployment Guide

## Docker Compose (Production)

```bash
docker-compose -f docker-compose.prod.yml up -d
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| MONGO_URI | MongoDB connection string |
| KEYCLOAK_URL | Keycloak base URL |
| KEYCLOAK_REALM | Auth realm |
| KEYCLOAK_CLIENT_ID | OAuth client ID |

## Health Checks

```bash
curl http://localhost:8080/health
```
