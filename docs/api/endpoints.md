# API Reference

Base URL: `http://localhost:8080`

## Events

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /events | Create event |
| GET | /events | List events |
| DELETE | /events/:id | Delete event |
| POST | /events/:id/sync | Sync to calendar |

## Auth

All endpoints require Bearer token from Keycloak.

```
Authorization: Bearer <token>
```

## Request/Response Examples

### Create Event
```json
POST /events
{
  "title": "Workout",
  "type": "workout",
  "start": "2024-01-15T09:00:00Z",
  "end": "2024-01-15T10:00:00Z"
}
```

### Response
```json
{
  "id": "...",
  "title": "Workout",
  "created_at": "..."
}
```
