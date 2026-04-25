# Jeeb

Personal management app: workouts, study, sleep, finance, calendar notifications.

## Stack
- Frontend: React | Backend: Go | DB: MongoDB | Auth: Keycloak | Infra: Kubernetes

## Ports (NodePort)
| Service | URL |
|---------|-----|
| Frontend | http://localhost:30000 |
| Backend API | http://localhost:30080 |
| Keycloak | http://localhost:30081 |
| Jenkins | http://localhost:30082 |
| Nexus UI | http://localhost:30083 |
| SonarQube | http://localhost:30090 |
| MongoDB | mongodb://localhost:30017 |
| Nexus Registry | localhost:30050 |

## Commands
```bash
# Apply all k8s manifests
bash k8s/apply.sh

# Check pods
kubectl get pods -n jeeb

# Logs
kubectl logs -n jeeb deployment/backend
kubectl logs -n jeeb deployment/frontend

# Restart
kubectl rollout restart deployment/backend -n jeeb
```

## Architecture (Clean/Hexagonal)
```
cmd/                 # Entry points
internal/
  domain/            # Entities, business rules
  usecase/           # Application logic
  port/in/           # Use case interfaces
  port/out/          # External service interfaces
  adapter/in/http/   # REST handlers
  adapter/out/       # MongoDB, integrations
pkg/                 # Shared utilities
```

## Key Interfaces
- `CalendarPort` - Google Calendar, etc.
- `NotificationPort` - LINE, Slack, Discord, Email

## API
```
POST   /events          # Create
GET    /events          # List
DELETE /events/:id      # Delete
POST   /events/:id/sync # Sync to calendars
```

## Collections
users, workouts, studies, sleep, finance, events, integrations

## Conventions
- Each feature (workout/study/sleep/finance) follows same pattern: domain → usecase → adapters
- Go: standard error handling, context propagation
- React: functional components, hooks
