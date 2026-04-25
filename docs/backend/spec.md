# Backend Specification

## Stack

| Layer | Choice | Why |
|-------|--------|-----|
| Language | Go 1.22+ | Performance, simplicity |
| Router | Chi | Lightweight, stdlib compatible |
| Auth | go-oidc | Keycloak JWT validation |
| Database | mongo-go-driver | Official MongoDB driver |
| Validation | go-playground/validator | Struct tag validation |
| Config | envconfig | Env-based config |
| Logging | slog | Stdlib structured logging |
| Testing | testify | Assertions, mocks |

## Project Structure

```
jeeb/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go            # App configuration
│   │
│   ├── domain/                  # Business entities
│   │   ├── user.go
│   │   ├── workout.go
│   │   ├── study.go
│   │   ├── sleep.go
│   │   ├── finance.go
│   │   └── event.go
│   │
│   ├── port/
│   │   ├── in/                  # Input ports (use case interfaces)
│   │   │   ├── user.go
│   │   │   ├── workout.go
│   │   │   ├── study.go
│   │   │   ├── sleep.go
│   │   │   ├── finance.go
│   │   │   └── event.go
│   │   │
│   │   └── out/                 # Output ports (repository interfaces)
│   │       ├── user_repository.go
│   │       ├── workout_repository.go
│   │       ├── study_repository.go
│   │       ├── sleep_repository.go
│   │       ├── finance_repository.go
│   │       ├── event_repository.go
│   │       ├── calendar_port.go
│   │       └── notification_port.go
│   │
│   ├── usecase/                 # Application logic
│   │   ├── user_usecase.go
│   │   ├── workout_usecase.go
│   │   ├── study_usecase.go
│   │   ├── sleep_usecase.go
│   │   ├── finance_usecase.go
│   │   └── event_usecase.go
│   │
│   └── adapter/
│       ├── in/
│       │   └── http/            # HTTP handlers
│       │       ├── router.go
│       │       ├── middleware/
│       │       │   ├── auth.go
│       │       │   ├── logging.go
│       │       │   └── recovery.go
│       │       ├── handler/
│       │       │   ├── user_handler.go
│       │       │   ├── workout_handler.go
│       │       │   ├── study_handler.go
│       │       │   ├── sleep_handler.go
│       │       │   ├── finance_handler.go
│       │       │   └── event_handler.go
│       │       ├── request/
│       │       │   └── *.go     # Request DTOs
│       │       └── response/
│       │           └── *.go     # Response DTOs
│       │
│       └── out/
│           ├── mongo/           # MongoDB implementations
│           │   ├── client.go
│           │   ├── user_repository.go
│           │   ├── workout_repository.go
│           │   ├── study_repository.go
│           │   ├── sleep_repository.go
│           │   ├── finance_repository.go
│           │   └── event_repository.go
│           │
│           └── integration/     # External services
│               ├── google_calendar.go
│               └── line_notify.go
│
├── pkg/                         # Shared utilities
│   ├── apperror/
│   │   └── error.go             # Application errors
│   └── pagination/
│       └── pagination.go
│
├── go.mod
├── go.sum
└── Dockerfile
```

## Configuration

```go
type Config struct {
    Server   ServerConfig
    MongoDB  MongoConfig
    Keycloak KeycloakConfig
}

type ServerConfig struct {
    Port         string `envconfig:"PORT" default:"8080"`
    ReadTimeout  int    `envconfig:"READ_TIMEOUT" default:"10"`
    WriteTimeout int    `envconfig:"WRITE_TIMEOUT" default:"10"`
}

type MongoConfig struct {
    URI      string `envconfig:"MONGO_URI" required:"true"`
    Database string `envconfig:"MONGO_DATABASE" default:"jeeb"`
}

type KeycloakConfig struct {
    URL      string `envconfig:"KEYCLOAK_URL" required:"true"`
    Realm    string `envconfig:"KEYCLOAK_REALM" required:"true"`
    ClientID string `envconfig:"KEYCLOAK_CLIENT_ID" required:"true"`
}
```

---

## Domain Models

### User
```go
type User struct {
    ID          string    `bson:"_id,omitempty"`
    KeycloakID  string    `bson:"keycloak_id"`
    Email       string    `bson:"email"`
    DisplayName string    `bson:"display_name"`
    CreatedAt   time.Time `bson:"created_at"`
    UpdatedAt   time.Time `bson:"updated_at"`
}
```

### Workout
```go
type Workout struct {
    ID        string     `bson:"_id,omitempty"`
    UserID    string     `bson:"user_id"`
    Type      WorkoutType `bson:"type"`
    Duration  int        `bson:"duration"` // minutes
    Exercises []Exercise `bson:"exercises"`
    Notes     string     `bson:"notes"`
    CreatedAt time.Time  `bson:"created_at"`
}

type WorkoutType string
const (
    WorkoutStrength    WorkoutType = "strength"
    WorkoutCardio      WorkoutType = "cardio"
    WorkoutFlexibility WorkoutType = "flexibility"
)

type Exercise struct {
    Name   string  `bson:"name"`
    Sets   int     `bson:"sets"`
    Reps   int     `bson:"reps"`
    Weight float64 `bson:"weight"` // kg
}
```

### Study
```go
type StudySession struct {
    ID        string    `bson:"_id,omitempty"`
    UserID    string    `bson:"user_id"`
    Subject   string    `bson:"subject"`
    Duration  int       `bson:"duration"` // minutes
    Notes     string    `bson:"notes"`
    CreatedAt time.Time `bson:"created_at"`
}
```

### Sleep
```go
type SleepRecord struct {
    ID        string    `bson:"_id,omitempty"`
    UserID    string    `bson:"user_id"`
    StartTime time.Time `bson:"start_time"`
    EndTime   time.Time `bson:"end_time"`
    Quality   int       `bson:"quality"` // 1-5
    Notes     string    `bson:"notes"`
    CreatedAt time.Time `bson:"created_at"`
}
```

### Finance
```go
type Transaction struct {
    ID        string          `bson:"_id,omitempty"`
    UserID    string          `bson:"user_id"`
    Type      TransactionType `bson:"type"`
    Amount    float64         `bson:"amount"`
    Category  string          `bson:"category"`
    Date      time.Time       `bson:"date"`
    Notes     string          `bson:"notes"`
    CreatedAt time.Time       `bson:"created_at"`
}

type TransactionType string
const (
    TransactionIncome  TransactionType = "income"
    TransactionExpense TransactionType = "expense"
)
```

### Event
```go
type Event struct {
    ID         string    `bson:"_id,omitempty"`
    UserID     string    `bson:"user_id"`
    Title      string    `bson:"title"`
    Type       EventType `bson:"type"`
    Start      time.Time `bson:"start"`
    End        time.Time `bson:"end"`
    ExternalID string    `bson:"external_id"` // Google Calendar ID
    CreatedAt  time.Time `bson:"created_at"`
}

type EventType string
const (
    EventWorkout  EventType = "workout"
    EventStudy    EventType = "study"
    EventSleep    EventType = "sleep"
    EventFinance  EventType = "finance"
    EventCustom   EventType = "custom"
)
```

---

## API Endpoints

### Auth
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /health | Health check (public) |
| GET | /me | Get current user |

### Workouts
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /workouts | List workouts |
| POST | /workouts | Create workout |
| GET | /workouts/:id | Get workout |
| PUT | /workouts/:id | Update workout |
| DELETE | /workouts/:id | Delete workout |
| GET | /workouts/stats | Get workout stats |

### Study
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /study | List study sessions |
| POST | /study | Create study session |
| GET | /study/:id | Get study session |
| PUT | /study/:id | Update study session |
| DELETE | /study/:id | Delete study session |
| GET | /study/stats | Get study stats |

### Sleep
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /sleep | List sleep records |
| POST | /sleep | Create sleep record |
| GET | /sleep/:id | Get sleep record |
| PUT | /sleep/:id | Update sleep record |
| DELETE | /sleep/:id | Delete sleep record |
| GET | /sleep/stats | Get sleep stats |

### Finance
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /finance | List transactions |
| POST | /finance | Create transaction |
| GET | /finance/:id | Get transaction |
| PUT | /finance/:id | Update transaction |
| DELETE | /finance/:id | Delete transaction |
| GET | /finance/stats | Get finance stats |
| GET | /finance/categories | List categories |

### Events
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /events | List events |
| POST | /events | Create event |
| GET | /events/:id | Get event |
| PUT | /events/:id | Update event |
| DELETE | /events/:id | Delete event |
| POST | /events/:id/sync | Sync to external calendar |

### Integrations
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /integrations | List integrations |
| POST | /integrations/google/connect | Connect Google Calendar |
| DELETE | /integrations/google | Disconnect Google |
| POST | /integrations/line/connect | Connect LINE Notify |
| DELETE | /integrations/line | Disconnect LINE |

---

## Request/Response Formats

### Pagination Request
```
GET /workouts?page=1&limit=20&sort=-created_at
```

### Pagination Response
```json
{
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Error Response
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request",
    "details": [
      {"field": "duration", "message": "must be greater than 0"}
    ]
  }
}
```

### Stats Response (Workout)
```json
{
  "this_week": 5,
  "this_month": 18,
  "total": 156,
  "streak": 7,
  "by_type": {
    "strength": 10,
    "cardio": 6,
    "flexibility": 2
  }
}
```

---

## Auth Flow

```
┌─────────┐     ┌─────────┐     ┌──────────┐     ┌─────────┐
│ Frontend│────>│Keycloak │────>│ Frontend │────>│ Backend │
│         │login│         │token│          │ API │         │
└─────────┘     └─────────┘     └──────────┘     └─────────┘
                                     │               │
                                     │ Bearer token  │
                                     └───────────────┘
                                           │
                                    ┌──────▼──────┐
                                    │ Validate JWT│
                                    │ (go-oidc)   │
                                    └──────┬──────┘
                                           │
                                    ┌──────▼──────┐
                                    │ Get/Create  │
                                    │ User        │
                                    └─────────────┘
```

### Auth Middleware
```go
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractBearerToken(r)
        if token == "" {
            respondError(w, ErrUnauthorized)
            return
        }

        idToken, err := m.verifier.Verify(r.Context(), token)
        if err != nil {
            respondError(w, ErrUnauthorized)
            return
        }

        var claims Claims
        idToken.Claims(&claims)

        ctx := context.WithValue(r.Context(), userCtxKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## Error Codes

| Code | HTTP | Description |
|------|------|-------------|
| UNAUTHORIZED | 401 | Missing/invalid token |
| FORBIDDEN | 403 | Not allowed |
| NOT_FOUND | 404 | Resource not found |
| VALIDATION_ERROR | 400 | Invalid input |
| CONFLICT | 409 | Duplicate resource |
| INTERNAL_ERROR | 500 | Server error |

---

## Port Interfaces

### Repository (Output Port)
```go
type WorkoutRepository interface {
    Create(ctx context.Context, workout *domain.Workout) error
    FindByID(ctx context.Context, id string) (*domain.Workout, error)
    FindByUserID(ctx context.Context, userID string, opts QueryOpts) ([]*domain.Workout, int64, error)
    Update(ctx context.Context, workout *domain.Workout) error
    Delete(ctx context.Context, id string) error
}
```

### Use Case (Input Port)
```go
type WorkoutUseCase interface {
    Create(ctx context.Context, userID string, req CreateWorkoutRequest) (*domain.Workout, error)
    GetByID(ctx context.Context, userID, id string) (*domain.Workout, error)
    List(ctx context.Context, userID string, opts QueryOpts) ([]*domain.Workout, *Pagination, error)
    Update(ctx context.Context, userID, id string, req UpdateWorkoutRequest) (*domain.Workout, error)
    Delete(ctx context.Context, userID, id string) error
    GetStats(ctx context.Context, userID string, period string) (*WorkoutStats, error)
}
```

### External Integration (Output Port)
```go
type CalendarPort interface {
    CreateEvent(ctx context.Context, userID string, event *domain.Event) (externalID string, err error)
    UpdateEvent(ctx context.Context, userID string, event *domain.Event) error
    DeleteEvent(ctx context.Context, userID, externalID string) error
}

type NotificationPort interface {
    Send(ctx context.Context, userID string, message string) error
}
```

---

## Dependencies

```go
// go.mod
module github.com/yourusername/jeeb

go 1.22

require (
    github.com/go-chi/chi/v5 v5.0.12
    github.com/go-chi/cors v1.2.1
    github.com/coreos/go-oidc/v3 v3.9.0
    go.mongodb.org/mongo-driver v1.14.0
    github.com/go-playground/validator/v10 v10.18.0
    github.com/kelseyhightower/envconfig v1.4.0
    github.com/stretchr/testify v1.8.4
)
```

---

## Testing Strategy

| Layer | Type | Coverage |
|-------|------|----------|
| Domain | Unit | Business rules |
| Use Case | Unit | Mock repositories |
| Handler | Integration | HTTP requests |
| Repository | Integration | Test MongoDB |

### Test Structure
```
internal/
  usecase/
    workout_usecase.go
    workout_usecase_test.go    # Unit tests with mocks
  adapter/
    in/http/handler/
      workout_handler.go
      workout_handler_test.go  # HTTP integration tests
    out/mongo/
      workout_repository.go
      workout_repository_test.go  # MongoDB integration tests
```
