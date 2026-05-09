# API Endpoints

## Base URLs

- Main backend, direct local process: `http://localhost:8080`
- Main backend, cluster NodePort: `http://localhost:30080`
- Learning backend, direct local process: `http://localhost:8081`
- Learning backend, cluster NodePort: `http://localhost:30086`
- Learning API through Kong: `http://localhost:30088/learning`

## Authentication

- All routes except `GET /health` and main-backend `GET /metrics` require `Authorization: Bearer <access-token>`.
- Both services create or update the local user record from token claims on authenticated requests.

## Shared conventions

- Paginated list routes accept `page`, `limit`, and `sort`.
- Defaults: `page=1`, `limit=20`, `sort=-created_at`.
- Paginated responses use:

```json
{
  "data": [],
  "meta": { "page": 1, "limit": 20, "total": 0, "total_pages": 0 }
}
```

- Error responses use:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "invalid request",
    "details": "..."
  }
}
```

## Main backend

| Method | Path | Notes |
|---|---|---|
| `GET` | `/health` | public |
| `GET` | `/metrics` | public Prometheus endpoint |
| `GET` | `/me` | current user |
| `GET,POST` | `/workouts/` | list or create |
| `GET` | `/workouts/stats` | workout stats |
| `GET,PUT,DELETE` | `/workouts/{id}` | item operations |
| `GET,POST` | `/study/` | list or create |
| `GET` | `/study/stats` | study stats |
| `GET,PUT,DELETE` | `/study/{id}` | item operations |
| `GET,POST` | `/sleep/` | list or create |
| `GET` | `/sleep/stats` | sleep stats |
| `GET,PUT,DELETE` | `/sleep/{id}` | item operations |
| `GET,POST` | `/finance/` | list or create |
| `GET` | `/finance/stats` | finance stats |
| `GET` | `/finance/categories` | distinct user categories |
| `GET,PUT,DELETE` | `/finance/{id}` | item operations |
| `GET,POST` | `/events/` | list or create |
| `GET,PUT,DELETE` | `/events/{id}` | item operations |
| `POST` | `/events/{id}/sync` | currently returns `503` unless a calendar provider is added |
| `GET` | `/master?category=exercise` | required `category` query param |

### Create payload examples

```json
{ "type": "strength", "duration": 60, "exercises": [], "notes": "Upper body" }
```

```json
{ "subject": "Algorithms", "duration": 90, "notes": "DP review" }
```

```json
{ "start_time": "2026-05-10T22:30:00Z", "end_time": "2026-05-11T06:30:00Z", "quality": 4, "notes": "" }
```

```json
{ "type": "expense", "amount": 12.5, "category": "food", "date": "2026-05-10T00:00:00Z", "notes": "Lunch" }
```

```json
{ "title": "Leg day", "type": "workout", "start": "2026-05-10T09:00:00Z", "end": "2026-05-10T10:00:00Z" }
```

## Learning backend

| Method | Path | Notes |
|---|---|---|
| `GET` | `/health` | public |
| `GET` | `/me` | current user |
| `GET` | `/stats` | progress summary by topic |
| `GET,POST` | `/topics/` | list or create |
| `GET,PUT,DELETE` | `/topics/{id}` | topic operations |
| `GET,POST` | `/topics/{id}/items` | list or create topic items |
| `PUT,DELETE` | `/topics/{id}/items/{itemId}` | item operations |
| `GET,DELETE` | `/topics/{id}/progress` | topic progress or reset |
| `PUT` | `/progress/{itemId}` | upsert item progress |

### Learning payload examples

```json
{ "name": "IPA Phonetics", "description": "International Phonetic Alphabet", "category": "Language", "icon": "IPA" }
```

```json
{ "term": "/p/", "meaning": "voiceless bilabial plosive", "example": "pit", "hint": "", "category": "Plosives", "sort_order": 1 }
```

```json
{ "topic_id": "<topic-id>", "status": "mastered" }
```

`status` must be either `learning` or `mastered`.
