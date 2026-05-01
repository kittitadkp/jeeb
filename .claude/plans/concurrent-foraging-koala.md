
# Plan: jeeb-learning — Full-Stack Learning Platform

## Context
Build a new standalone project **`jeeb-learning`** — a web app for learning topics (starting with IPA phonetics). Master data (topics + items) and user progress are stored in MongoDB database `jeeb_learning`. Each topic supports multiple **study tools** (Flashcard, Recall) — all tools feed into the same unified progress. Follows exact same stack and patterns as jeeb backend (Go/Chi/mongo-driver) and frontend (React 19/TanStack Query/design tokens).

---

## Architecture

```
jeeb-learning-backend  (Go API, port 30086)
jeeb-learning-frontend (React, port 30087)
MongoDB: jeeb_learning database (same pod, separate DB)
```

---

## Study Tools

Tools are frontend-only modes — all write to the same `UserProgress` via `PUT /progress/:itemId`. No backend changes needed per tool; adding a new tool is a new React component only.

| Tool | Description | Interaction |
|---|---|---|
| **Flashcard** | See term → flip card → see meaning + example | Know it ✓ / Still learning ✗ |
| **Recall** | See meaning + example → type the term from memory | Submit → correct/incorrect auto-check |

### Flashcard UI
```
┌──────────────────────────────┐
│  Plosives                    │
│                              │
│          /p/                 │   ← front (term)
│                              │
└──────────────────────────────┘
         [ Reveal ]

┌──────────────────────────────┐
│  voiceless bilabial plosive  │
│  Example: "pit"              │   ← back (after reveal)
└──────────────────────────────┘
    [✗ Still learning]  [✓ Know it]
```

### Recall UI
```
┌──────────────────────────────┐
│  voiceless bilabial plosive  │
│  Example: "pit"              │   ← prompt (meaning shown)
└──────────────────────────────┘
  Type the IPA symbol:
  [ __________ ] [ Submit ]

  ✓ Correct! — /p/             ← or ✗ Incorrect — answer was /p/
  [ Next → ]
```

### Practice Session Flow (both tools)
- Draws a shuffled queue from items not yet mastered (falls back to all if all mastered)
- Session progress bar: X / Y items done
- On completion: summary card (correct count, time taken)
- Progress auto-saved after each item via `PUT /progress/:itemId`

---

## Data Models (`internal/domain/`)

### `topic.go`
```go
type Topic struct {
    ID          string    `bson:"_id,omitempty" json:"id"`
    Name        string    `bson:"name" json:"name"`
    Description string    `bson:"description" json:"description"`
    Category    string    `bson:"category" json:"category"`
    Icon        string    `bson:"icon" json:"icon"`
    CreatedAt   time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
```

### `item.go`
```go
type Item struct {
    ID        string    `bson:"_id,omitempty" json:"id"`
    TopicID   string    `bson:"topic_id" json:"topic_id"`
    Term      string    `bson:"term" json:"term"`         // "/p/"
    Meaning   string    `bson:"meaning" json:"meaning"`   // "voiceless bilabial plosive"
    Example   string    `bson:"example" json:"example"`   // "pit"
    Hint      string    `bson:"hint" json:"hint"`
    Category  string    `bson:"category" json:"category"` // "Plosives"
    SortOrder int       `bson:"sort_order" json:"sort_order"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
```

### `progress.go`
```go
type UserProgress struct {
    ID             string    `bson:"_id,omitempty" json:"id"`
    UserID         string    `bson:"user_id" json:"user_id"`
    TopicID        string    `bson:"topic_id" json:"topic_id"`
    ItemID         string    `bson:"item_id" json:"item_id"`
    Status         string    `bson:"status" json:"status"` // "learning" | "mastered"
    ReviewCount    int       `bson:"review_count" json:"review_count"`
    LastReviewedAt time.Time `bson:"last_reviewed_at" json:"last_reviewed_at"`
    CreatedAt      time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}
```

---

## Backend — File Structure

```
learning-backend/
├── cmd/
│   ├── api/main.go           # wire: config → mongo → repos → usecases → handlers → router
│   └── seed/main.go          # seeds IPA topic + 44 items
├── internal/
│   ├── domain/               # topic.go, item.go, progress.go
│   ├── port/
│   │   ├── in/               # TopicUseCase, ItemUseCase, ProgressUseCase interfaces + DTOs
│   │   └── out/repositories.go
│   ├── usecase/              # topic_usecase.go, item_usecase.go, progress_usecase.go
│   ├── adapter/
│   │   ├── in/http/
│   │   │   ├── router.go
│   │   │   ├── handler/      # topic_handler.go, item_handler.go, progress_handler.go
│   │   │   └── middleware/   # RespondJSON, RespondError, Auth (same as jeeb-backend)
│   │   └── out/mongo/        # topic_repository.go, item_repository.go, progress_repository.go
│   └── config/config.go
├── pkg/apperror/ + pkg/pagination/   # copied from jeeb-backend
└── go.mod  (module: github.com/kittitadkp/jeeb-learning)
```

## Backend — API Routes

```
GET  /health
GET  /me

GET    /topics
POST   /topics
GET    /topics/:id
PUT    /topics/:id
DELETE /topics/:id

GET    /topics/:id/items          (?category=, ?page=, ?limit=)
POST   /topics/:id/items
PUT    /topics/:id/items/:itemId
DELETE /topics/:id/items/:itemId

GET    /topics/:id/progress       user progress map {item_id → status}
PUT    /progress/:itemId          upsert {status: "learning"|"mastered"}
DELETE /topics/:id/progress       reset topic progress
GET    /stats                     [{topic_id, name, mastered, learning, total}]
```

MongoDB collections: `topics`, `items`, `progress`
Indexes: `progress(user_id, item_id)` unique; `progress(user_id, topic_id)`.

---

## IPA Seed Data (`cmd/seed/main.go`)

44 items across 8 categories:

| Category | Examples |
|---|---|
| Plosives (6) | /p/ pit, /b/ bit, /t/ tip, /d/ dip, /k/ cat, /ɡ/ gap |
| Fricatives (9) | /f/ fat, /v/ vat, /θ/ thin, /ð/ this, /s/ sat, /z/ zap, /ʃ/ ship, /ʒ/ vision, /h/ hat |
| Affricates (2) | /tʃ/ chip, /dʒ/ jam |
| Nasals (3) | /m/ map, /n/ nap, /ŋ/ sing |
| Approximants (4) | /l/ lip, /r/ rip, /j/ yes, /w/ wet |
| Short Vowels (7) | /ɪ/ bit, /e/ bet, /æ/ bat, /ʌ/ but, /ɒ/ bot, /ʊ/ book, /ə/ about |
| Long Vowels (5) | /iː/ beat, /ɑː/ bar, /ɔː/ bore, /uː/ boot, /ɜː/ bird |
| Diphthongs (8) | /eɪ/ bait, /aɪ/ bite, /ɔɪ/ boy, /əʊ/ boat, /aʊ/ bout, /ɪə/ beer, /eə/ bear, /ʊə/ tour |

---

## Frontend — File Structure

```
learning-frontend/src/
├── lib/api.ts, design.ts, auth.tsx, utils.ts   (same patterns as jeeb-frontend)
├── types/index.ts           # Topic, Item, UserProgress, TopicStats
├── hooks/
│   ├── useTopics.ts         # useTopics, useTopic, useCreateTopic, useUpdateTopic, useDeleteTopic
│   ├── useItems.ts          # useItems, useCreateItem, useUpdateItem, useDeleteItem
│   └── useProgress.ts       # useProgress, useUpsertProgress, useResetProgress, useStats
├── pages/
│   ├── Home.tsx             # topic cards + overall stats
│   └── Topic.tsx            # 4 tabs: Browse | Flashcard | Recall | Progress
├── components/
│   ├── study/
│   │   ├── FlashcardTool.tsx    # flashcard session component
│   │   └── RecallTool.tsx       # recall/typing session component
│   └── ui/                      # Button, Card, Badge, StatCard, SectionLabel, States
├── store/theme.ts
├── App.tsx                  # / → Home, /topics/:id → Topic
└── main.tsx
```

## Frontend — Page Designs

### Home (`/`)
- `SectionLabel` "🎓 Learning"
- Overall `StatCard`s: total mastered, active topics, study streak
- Topic cards grid: name, icon, description, progress bar `mastered/total`, [Study →]

### Topic (`/topics/:id`) — 4 tabs

**Browse tab:**
- Search + category filter pills
- Item card grid: large term, meaning, example, `Badge(category)`, status dot (●mastered / ○learning / ·new)

**Flashcard tab:**
- `FlashcardTool` component (see tool design above)
- Term on front, meaning + example on back
- Writes `PUT /progress/:itemId` on each Know it / Still learning

**Recall tab:**
- `RecallTool` component (see tool design above)
- Meaning + example shown, user types term
- Case/whitespace-insensitive match
- Writes `PUT /progress/:itemId` on correct (mastered) or incorrect (learning)

**Progress tab:**
- `StatCard`s: Mastered / Learning / Not Started
- Progress bar
- Mastered item badges
- [Reset Progress]

---

## K8s Changes — `k8s/charts/jeeb-app/`

### New templates
```
templates/learning/
  deployment.yaml   (Vault sidecar, app: learning)
  service.yaml      (NodePort 30086)
  configmap.yaml    (MONGO_URI → jeeb_learning DB, KEYCLOAK_URL)
  serviceaccount.yaml
```

### `values.yaml` additions
```yaml
learning:
  replicas: 1
  nodePort: 30086
  image: ""
  vault:
    path: secret/data/jeeb/learning/develop
    envFile: .env.develop
    role: learning
  database: jeeb_learning
```

---

## Verification
1. `go run ./cmd/seed/main.go` → `topics` and `items` seeded in MongoDB
2. `curl localhost:8080/topics` → returns IPA topic
3. `curl localhost:8080/topics/:id/items` → returns 44 items
4. `npm run dev` → Home lists IPA topic, `/topics/:id` Browse shows all 44 symbols
5. Flashcard: flip works, Know it → status=mastered in DB, Progress tab reflects it
6. Recall: correct answer → mastered, wrong → learning, case-insensitive match works
7. Reset: all progress cleared, Progress tab shows 0/44
8. Dark mode: all components respect CSS vars
