# Jeeb

Personal management app: workouts, study, sleep, finance, calendar notifications.

## Stack
- Frontend: React | Backend: Go | DB: MongoDB | Auth: Keycloak | Infra: Docker Compose

## Commands
```bash
docker-compose up --build    # Build and run
docker-compose up -d         # Background
```

## Ports
Frontend:3000 | API:8080 | MongoDB:27017 | Keycloak:8081

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
