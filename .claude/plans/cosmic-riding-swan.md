# Plan: Master Data CRUD — Settings UI

## Context

The backend has a `MasterRecord` system (category + name + defaults map) used to seed exercise data, but only exposes a read-only `GET /master?category=` endpoint. The frontend Settings page has no way to manage this data. The goal is to add full CRUD for master records (starting with `exercise` under workout, with strength/cardio/flexibility muscle groups), wire it end-to-end, and design the UI to be extensible for future master categories (study subjects, sleep types, etc.).

---

## Backend

### 1. Port interfaces

**`backend/internal/port/in/master.go`** — extend `MasterUseCase`:
```go
Create(ctx, category, name string, defaults map[string]interface{}) (*domain.MasterRecord, error)
Update(ctx, id, name string, defaults map[string]interface{}) (*domain.MasterRecord, error)
Delete(ctx, id string) error
```

**New file: `backend/internal/port/out/master.go`** — `MasterRepository` interface:
```go
FindByCategory(ctx, category string) ([]*domain.MasterRecord, error)
FindByID(ctx, id string) (*domain.MasterRecord, error)
Insert(ctx, record *domain.MasterRecord) (*domain.MasterRecord, error)
Update(ctx, id, name string, defaults map[string]interface{}) (*domain.MasterRecord, error)
DeleteByID(ctx, id string) error
Count(ctx) (int64, error)
InsertMany(ctx, records []*domain.MasterRecord) error
```

### 2. Mongo repository

**`backend/internal/adapter/out/mongo/master_repository.go`** — add:
- `FindByID` — find by `_id` ObjectID
- `Insert` — insert one document, return with generated ID
- `Update` — find by ID + update `name` and `defaults`
- `DeleteByID` — delete by `_id` ObjectID

### 3. Use case

**`backend/internal/usecase/master.go`** — implement `Create`, `Update`, `Delete` delegating to repo. `Create` builds a `domain.MasterRecord{Category: category, Name: name, Defaults: defaults}`.

### 4. Handler

**`backend/internal/adapter/in/http/handler/master_handler.go`** — add:
- `Create(w, r)` — decode JSON body `{name, category, defaults}`, call `uc.Create`
- `Update(w, r)` — URL param `{id}`, decode body `{name, defaults}`, call `uc.Update`
- `Delete(w, r)` — URL param `{id}`, call `uc.Delete`

Request structs with `validate` tags (use existing `middleware.DecodeAndValidate` pattern from workout handler).

### 5. Router

**`backend/internal/adapter/in/http/router.go`** — add authenticated routes:
```
POST   /master          → master.Create
PUT    /master/{id}     → master.Update
DELETE /master/{id}     → master.Delete
```

---

## Frontend

### 6. Hooks

**`frontend/src/hooks/useMaster.ts`** — add:
```ts
useCreateMaster()   → POST /master
useUpdateMaster()   → PUT /master/{id}
useDeleteMaster()   → DELETE /master/{id}
```
Each mutation invalidates `["master", category]` on success (optimistic from query key).

### 7. Types

**`frontend/src/types/index.ts`** — add:
```ts
interface CreateMasterRequest { category: string; name: string; defaults: Record<string, unknown> }
interface UpdateMasterRequest { name: string; defaults: Record<string, unknown> }
```

### 8. MasterDataCard component (generic + reusable)

**New file: `frontend/src/components/MasterDataCard.tsx`**

A self-contained card that manages any master category. Props:
```ts
interface MasterCategoryConfig {
  category: string;
  label: string;
  defaultFields: { key: string; label: string; type: 'number' | 'string'; defaultValue: number | string }[];
}
interface MasterDataCardProps { config: MasterCategoryConfig }
```

UI inside the card:
- Header: category label + "Add" button
- List: each record row with name, defaults preview chips, Edit (pencil) + Delete (trash) icon buttons
- Inline form (shown below list on Add/Edit): name input + one input per `defaultFields` entry → Save / Cancel buttons
- Delete: confirm inline (replaces the row with "Are you sure? Delete / Cancel")
- Loading/empty states using existing patterns

### 9. Category configs

**`frontend/src/constants/master.ts`** — define configs:
```ts
export const MASTER_CATEGORY_CONFIGS: MasterCategoryConfig[] = [
  {
    category: "exercise",
    label: "Exercises",
    defaultFields: [
      { key: "muscle_group", label: "Muscle Group", type: "string", defaultValue: "" },
      { key: "sets",         label: "Sets",         type: "number", defaultValue: 3 },
      { key: "reps",         label: "Reps",         type: "number", defaultValue: 10 },
      { key: "rest_seconds", label: "Rest (s)",     type: "number", defaultValue: 60 },
    ],
  },
  // Future: { category: "study_subject", ... }, { category: "sleep_type", ... }
];
```

### 10. Settings page

**`frontend/src/pages/Settings.tsx`** — add a new card after "Default Goal Targets":

```tsx
{/* Master Data */}
{MASTER_CATEGORY_CONFIGS.map(config => (
  <MasterDataCard key={config.category} config={config} />
))}
```

The card spans one column (same as other cards, grid auto-fill handles layout).

---

## File checklist

| File | Change |
|------|--------|
| `backend/internal/port/in/master.go` | Add Create, Update, Delete to interface |
| `backend/internal/port/out/master.go` | **New** — MasterRepository interface |
| `backend/internal/adapter/out/mongo/master_repository.go` | Add FindByID, Insert, Update, DeleteByID |
| `backend/internal/usecase/master.go` | Implement Create, Update, Delete |
| `backend/internal/adapter/in/http/handler/master_handler.go` | Add Create, Update, Delete handlers |
| `backend/internal/adapter/in/http/router.go` | Add POST/PUT/DELETE /master routes |
| `frontend/src/types/index.ts` | Add CreateMasterRequest, UpdateMasterRequest |
| `frontend/src/hooks/useMaster.ts` | Add mutation hooks |
| `frontend/src/constants/master.ts` | **New** — MASTER_CATEGORY_CONFIGS |
| `frontend/src/components/MasterDataCard.tsx` | **New** — generic CRUD card |
| `frontend/src/pages/Settings.tsx` | Add MasterDataCard instances |

---

## Verification

1. **Backend**: `cd backend && go build ./...` — no compile errors
2. **Backend tests**: `go test ./...`
3. **API smoke test** (with valid token):
   - `POST /master` `{"category":"exercise","name":"Test Exercise","defaults":{"muscle_group":"test","sets":3,"reps":10,"rest_seconds":60}}`
   - `PUT /master/{id}` with updated name
   - `DELETE /master/{id}`
   - `GET /master?category=exercise` — verify seeded + created records
4. **Frontend**: `cd frontend && npm run dev` → navigate to Settings → verify "Exercises" card renders with seeded exercises
5. **CRUD flow**: Add new exercise → appears in list; Edit name/defaults → updated; Delete → removed; cancel confirms no change
6. **Extensibility check**: Add a second entry to `MASTER_CATEGORY_CONFIGS` → second card appears in Settings with zero additional code
