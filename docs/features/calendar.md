# Calendar and Events

## API-backed calendar

The `Calendar` page in the main frontend uses the main backend `/events` resource:

- `GET|POST /events/`
- `GET|PUT|DELETE /events/{id}`
- `POST /events/{id}/sync`

Event types are `workout`, `study`, `sleep`, `finance`, and `custom`.

## Current limitation

`POST /events/{id}/sync` is implemented in the API surface but does not have a real calendar provider at runtime, so it currently returns `503`.

## Separate local Events page

The `/events` page in the main frontend is not the same feature. It is a local to-do style page with seeded browser state and no backend persistence.
