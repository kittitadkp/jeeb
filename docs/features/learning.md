# Learning Feature

## Scope

The learning product is separate from the main study tracker. It uses `learning-frontend/` and `learning-backend/`.

## API

- Topics: `GET|POST /topics/`, `GET|PUT|DELETE /topics/{id}`
- Items: `GET|POST /topics/{id}/items`, `PUT|DELETE /topics/{id}/items/{itemId}`
- Progress: `GET|DELETE /topics/{id}/progress`, `PUT /progress/{itemId}`
- Summary: `GET /stats`

## Progress rules

- Allowed statuses: `learning`, `mastered`
- Progress is stored per `user_id` and `item_id`
- Resetting topic progress deletes that user's records for the topic

## Seed content

`go run ./cmd/seed` creates one `IPA Phonetics` topic and 44 learning items.
