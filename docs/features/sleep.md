# Sleep Feature

## API

- `GET|POST /sleep/`
- `GET /sleep/stats`
- `GET|PUT|DELETE /sleep/{id}`

## Data model

- `start_time`: required timestamp
- `end_time`: required timestamp
- `quality`: required integer from `1` to `5`
- `notes`: optional

## Stats

The backend returns weekly and monthly totals plus `avg_duration_minutes` and `avg_quality`.
