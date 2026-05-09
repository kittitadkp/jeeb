# Workout Feature

## Backing API

- `GET|POST /workouts/`
- `GET /workouts/stats`
- `GET|PUT|DELETE /workouts/{id}`
- `GET /master?category=exercise`

## Data model

- `type`: `strength`, `cardio`, or `flexibility`
- `duration`: required, greater than `0`
- `exercises`: optional array with `name`, `sets`, `reps`, `weight`, `rest_seconds`
- `notes`: optional

## Runtime notes

- Dashboard and workout pages read persisted workout data.
- Exercise reference data comes from the auto-seeded `master` collection.
