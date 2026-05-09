# Study Feature

## Main tracker study

### API

- `GET|POST /study/`
- `GET /study/stats`
- `GET|PUT|DELETE /study/{id}`

### Data model

- `subject`: required
- `duration`: required, greater than `0`
- `notes`: optional

## Learning app

The learning application is a separate feature set. See [learning.md](learning.md) for topic, item, and progress behavior.
