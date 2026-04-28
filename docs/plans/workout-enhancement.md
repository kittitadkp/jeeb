# Workout Feature Enhancement Plan

## Context

The workout feature has basic CRUD but the tracker uses hardcoded exercises and there's no exercise library. This plan adds a generic master collection — pure reference/library data only, never stores user logs. Workout logs continue to go to the workouts collection unchanged. The master collection supports exercise, study, and future domains via a category field. The tracker is redesigned into a 3-phase wizard (plan → active → summary) where users pick from the master list and can edit goals mid-workout.

---

## Backend Changes

### 1. Domain — `backend/internal/domain/master.go` (new)

Generic master record — reference data only, no user logs stored here:

```go
type MasterRecord struct {
    ID       string                 `bson:"_id,omitempty" json:"id"`
    Category string                 `bson:"category" json:"category"`   // "exercise", "study_topic", etc.
    Name     string                 `bson:"name" json:"name"`
    Defaults map[string]interface{} `bson:"defaults" json:"defaults"`   // flexible per category
}
```

Example exercise record defaults:
```json
{ "muscle_group": "chest", "sets": 4, "reps": 8, "rest_seconds": 90 }
```

> **Rule:** master collection = library/reference data only. User workout logs → workouts collection (unchanged).

### 2. Domain — `backend/internal/domain/workout.go`

- Add `RestSeconds int` to `Exercise` struct (`bson:"rest_seconds" json:"rest_seconds"`)
- This enhances the workouts collection: each saved exercise subdocument now stores the actual rest time used during the session. Existing documents without `rest_seconds` default to 0 (zero-value, backward compatible — no migration needed).

### 3. Port/In — `backend/internal/port/in/master.go` (new)

```go
type MasterUseCase interface {
    ListByCategory(ctx context.Context, category string) ([]*domain.MasterRecord, error)
}
```

### 4. Port/In — `backend/internal/port/in/workout.go`

- Add `RestSeconds int \`json:"rest_seconds" validate:"min=0"\`` to `ExerciseRequest`

### 5. Port/Out — `backend/internal/port/out/repositories.go`

Add:

```go
type MasterRepository interface {
    FindByCategory(ctx context.Context, category string) ([]*domain.MasterRecord, error)
    Count(ctx context.Context) (int64, error)
    InsertMany(ctx context.Context, records []*domain.MasterRecord) error
}
```

### 6. Usecase — `backend/internal/usecase/master.go` (new)

- Implements `MasterUseCase`, delegates `ListByCategory` to repo

### 7. Usecase — `backend/internal/usecase/workout_usecase.go`

- In the `ExerciseRequest` → `domain.Exercise` mapping, add `RestSeconds: req.RestSeconds`

### 8. Mongo Repo — `backend/internal/adapter/out/mongo/master_repository.go` (new)

- Collection: `"master"`
- `FindByCategory`: `coll.Find(ctx, bson.M{"category": category})`
- `Count`: `coll.CountDocuments(ctx, bson.D{})`
- `InsertMany`: convert `[]*domain.MasterRecord` → `[]interface{}` then `coll.InsertMany`
- Create index on `{ category: 1, name: 1 }` at init

### 9. Handler — `backend/internal/adapter/in/http/handler/master_handler.go` (new)

```go
// GET /master?category=exercise
func (h *MasterHandler) ListByCategory(w http.ResponseWriter, r *http.Request) {
    category := r.URL.Query().Get("category")
    if category == "" { respondError(w, 400, errors.New("category required")); return }
    records, err := h.uc.ListByCategory(r.Context(), category)
    ...respondJSON(w, 200, records)
}
```

### 10. Router — `backend/internal/adapter/in/http/router.go`

- Add `Master *handler.MasterHandler` to `Handlers` struct
- Add inside the authenticated group:
  ```go
  r.Get("/master", h.Master.ListByCategory)
  ```

### 11. Main — `backend/cmd/api/main.go`

- Wire: `masterRepo` → `masterUC` → `MasterHandler`
- Add seed function `seedMaster(ctx, repo)`:
  - If `Count > 0` → skip
  - Insert exercise records (`category: "exercise"`) — 18 exercises across muscle groups
  - Future: add study topics, sleep routines etc. here
  - On error: `log.Printf("WARN: seed master: %v", err)` — do NOT fatal
- Call seed before `http.ListenAndServe`

**Seed exercise records** (`category: "exercise"`):

| Name | Muscle Group |
|---|---|
| Bench Press | chest |
| Incline Dumbbell Press | chest |
| Pull Ups | back |
| Barbell Row | back |
| Deadlift | back |
| Overhead Press | shoulders |
| Lateral Raise | shoulders |
| Face Pull | shoulders |
| Squat | legs |
| Romanian Deadlift | legs |
| Leg Press | legs |
| Hip Thrust | legs |
| Bicep Curl | arms |
| Hammer Curl | arms |
| Tricep Pushdown | arms |
| Dip | arms |
| Plank | core |
| Cable Crunch | core |

---

## Frontend Changes

### 1. Types — `frontend/src/types/index.ts`

- Add `rest_seconds: number` to `Exercise`
- Add:

```ts
interface MasterRecord {
  id: string;
  category: string;
  name: string;
  defaults: Record<string, unknown>;
}
```

### 2. Hook — `frontend/src/hooks/useMaster.ts` (new)

```ts
export function useMaster(category: string) {
  return useQuery({
    queryKey: ['master', category],
    queryFn: () => api.get<MasterRecord[]>(`/master?category=${category}`),
    staleTime: Infinity,
    gcTime: Infinity,
  });
}
```

### 3. Shared Component — `frontend/src/components/ExercisePicker.tsx` (new)

Props:
```ts
interface ExerciseFormEntry { name: string; sets: number; reps: number; rest_seconds: number; weight: number; }
interface ExercisePickerProps { selected: ExerciseFormEntry[]; onChange: (e: ExerciseFormEntry[]) => void; }
```

- Calls `useMaster('exercise')` internally
- Search input (Lucide `Search`) filters by name client-side
- Results grouped by `defaults.muscle_group`, headers in `text-slate-500 text-xs uppercase`
- Each row: name + `+` button (Lucide `Plus`)
- Selected cards: name + 3 number inputs (Sets / Reps / Rest s) initialized from defaults + remove (Lucide `X`)

### 4. Tracker Redesign — `frontend/src/pages/Workouts.tsx`

Local state:
```ts
type TrackerPhase = 'plan' | 'active' | 'summary';
interface TrackerExercise {
  name: string;
  goalSets: number; goalReps: number; goalRestSeconds: number;
  completedSets: number;
  setLog: Array<{ reps: number; weight: number }>;
}
```

**Phase 1 — Plan:**
- Renders `<ExercisePicker>` → maps output to `TrackerExercise[]`
- "Start Workout" disabled until ≥1 exercise → sets `phase = 'active'`, records `startedAt`

**Phase 2 — Active:**
- Shows current exercise/set: `"Set 2 of 3 · 10 reps · 90s rest"`
- Lucide `Pencil` → inline edit `goalSets`/`goalReps`/`goalRestSeconds` (min `goalSets = completedSets+1`)
- Actual reps + weight inputs
- "Complete Set" → logs set → rest countdown (`useRef` interval, MM:SS `text-4xl font-mono`, "Skip Rest", auto-advance at 0)
- Clear interval in `useEffect` cleanup to prevent leaks
- All sets done → next exercise; all exercises done → `phase = 'summary'`

**Phase 3 — Summary:**
- Per-exercise: goal vs actual, set-by-set log
- "Save Workout" → `useCreateWorkout` with:
  - `duration: Math.max(1, Math.round(elapsed / 60000))`
  - `exercises` from `setLog`

### 5. Log Form — `frontend/src/pages/Workouts.tsx`

- Add collapsible "Exercises" section (Lucide `ChevronDown`/`ChevronUp`) below Notes in create form
- Renders `<ExercisePicker>` — output maps to `ExerciseRequest[]` on submit

---

## API Summary

| Method | Path | Description |
|--------|------|-------------|
| GET | `/master?category=exercise` | List exercise master records |
| GET | `/master?category=study_topic` | List study master records (future) |

---

## Verification

1. **Seed:** `kubectl exec -n jeeb deployment/backend -- wget -qO- 'http://localhost:8080/master?category=exercise'` → 18 records
2. **Tracker:** Start Workout → search/pick exercises → set goals → Start → complete sets with rest countdown → summary → Save → appears in list
3. **Edit goals mid-workout:** Pencil icon → change reps → header updates immediately
4. **Log form:** Create → expand Exercises → add 2 → save → workout card shows exercises
5. **Tests:** `go test ./internal/usecase/... ./internal/adapter/...`
