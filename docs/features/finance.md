# Finance Feature

## API

- `GET|POST /finance/`
- `GET /finance/stats`
- `GET /finance/categories`
- `GET|PUT|DELETE /finance/{id}`

## Data model

- `type`: `income` or `expense`
- `amount`: required, greater than `0`
- `category`: required
- `date`: required timestamp
- `notes`: optional

## Stats

Finance stats include current-month `income`, `expense`, and `net`, plus totals grouped by category.
