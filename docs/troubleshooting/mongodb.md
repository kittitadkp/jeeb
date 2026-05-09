# MongoDB Troubleshooting

## Expected local ports and databases

- MongoDB NodePort: `30017`
- Main backend database: `jeeb`
- Learning backend database: `jeeb_learning`

## Main backend has no exercise master data

The main backend seeds `master` records on startup when the collection is empty. If the API started before MongoDB was ready, restart the service and check logs.

## Learning content is missing

The learning backend does not auto-seed. Run:

```powershell
cd learning-backend
go run ./cmd/seed
```

## Collections

- Main backend: `users`, `workouts`, `studies`, `sleep`, `finance`, `events`, `master`
- Learning backend: `users`, `topics`, `items`, `progress`
