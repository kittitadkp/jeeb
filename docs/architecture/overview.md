# Architecture Overview

## Clean Architecture

```
┌─────────────────────────────────────┐
│           HTTP Handlers             │  ← adapter/in
├─────────────────────────────────────┤
│            Use Cases                │  ← usecase
├─────────────────────────────────────┤
│             Domain                  │  ← domain
├─────────────────────────────────────┤
│    MongoDB │ Calendar │ Notify      │  ← adapter/out
└─────────────────────────────────────┘
```

## Dependency Rule

Dependencies point inward. Domain has no external dependencies.

## Ports

| Port | Purpose |
|------|---------|
| `CalendarPort` | External calendar sync |
| `NotificationPort` | Push notifications |
| `*Repository` | Data persistence |

## Feature Modules

Each feature (workout, study, sleep, finance) is self-contained:
- Own domain models
- Own use cases
- Own repository interface
