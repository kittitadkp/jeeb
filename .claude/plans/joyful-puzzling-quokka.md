# Plan: Lazy Load Workout Cards (Infinite Scroll + Server-Side Filtering)

## Context

Workout cards currently use page-based pagination (prev/next buttons) with client-side type and period filtering. This has two problems:
1. Filtering only applies to the currently loaded page — a "strength" filter on page 1 misses strength workouts on page 2.
2. Page navigation is clunky; infinite scroll is a better fit for a card list.

This plan adds infinite scroll (auto-load-more via IntersectionObserver) and moves type/period filtering to the backend so pagination and filtering are consistent across the full dataset.

---

## Backend Changes

### [x] 1. Add `WorkoutListParams` to the input port
**File:** `backend/internal/port/in/workout.go`

Add a new struct that extends pagination params with filter fields:
```go
type WorkoutListParams struct {
    pagination.Params
    Type string     // optional: "strength" | "cardio" | "flexibility" | ""
    From *time.Time // optional: lower bound on created_at
}
```

Update `WorkoutUseCase` interface:
```go
List(ctx context.Context, userID string, params WorkoutListParams) ([]*domain.Workout, *pagination.Meta, error)
```

### [x] 2. Update repository interface
**File:** `backend/internal/port/out/repositories.go`

Change `WorkoutRepository.FindByUserID` signature:
```go
FindByUserID(ctx context.Context, userID string, params in.WorkoutListParams) ([]*domain.Workout, int64, error)
```

Note: This creates a cross-port import (`port/out` imports `port/in`). Alternative is to define `WorkoutListParams` in a shared `domain` or `pkg` package. Prefer defining it in `port/in` and importing it — the dependency direction (out → in) is acceptable here since `port/out` is consumed by adapters, not domain.

### [x] 3. Update usecase
**File:** `backend/internal/usecase/workout_usecase.go`

Pass `WorkoutListParams` through from use case to repository unchanged.

### [x] 4. Update handler to parse filters
**File:** `backend/internal/adapter/in/http/handler/workout_handler.go`

In `List()`, after calling `pagination.FromRequest(r)`, parse additional params:
```go
opts := pagination.FromRequest(r)
params := in.WorkoutListParams{Params: opts}

if t := r.URL.Query().Get("type"); t != "" {
    params.Type = t
}
if from := r.URL.Query().Get("from"); from != "" {
    t, err := time.Parse(time.RFC3339, from)
    if err == nil {
        params.From = &t
    }
}
```

### [x] 5. Update MongoDB repository
**File:** `backend/internal/adapter/out/mongo/workout_repository.go`

Update `FindByUserID` to build a dynamic filter:
```go
filter := bson.M{"user_id": userID}
if params.Type != "" {
    filter["type"] = params.Type
}
if params.From != nil {
    filter["created_at"] = bson.M{"$gte": *params.From}
}
```

---

## Frontend Changes

### [x] 6. Add `useInfiniteWorkouts` hook
**File:** `frontend/src/hooks/useWorkouts.ts`

Add a new export (keep the old `useWorkouts` to avoid breaking anything in the short term):
```ts
export interface WorkoutFilters {
  type?: string;  // "" | "strength" | "cardio" | "flexibility"
  from?: string;  // ISO 8601 date string, derived from period
}

export function useInfiniteWorkouts(filters: WorkoutFilters = {}) {
  return useInfiniteQuery({
    queryKey: ["workouts", "infinite", filters],
    queryFn: ({ pageParam = 1 }) => {
      const params = new URLSearchParams({ page: String(pageParam), limit: "20" });
      if (filters.type) params.set("type", filters.type);
      if (filters.from) params.set("from", filters.from);
      return api.get<PagedResponse<Workout>>(`/workouts?${params}`);
    },
    initialPageParam: 1,
    getNextPageParam: (last) =>
      last.meta.page < last.meta.total_pages ? last.meta.page + 1 : undefined,
  });
}
```

Add a helper to convert `WorkoutPeriod` → `from` ISO string:
```ts
export function periodToFrom(period: WorkoutPeriod): string | undefined {
  const now = new Date();
  if (period === "week")  return new Date(now.setDate(now.getDate() - 7)).toISOString();
  if (period === "month") return new Date(now.setMonth(now.getMonth() - 1)).toISOString();
  if (period === "year")  return new Date(now.setFullYear(now.getFullYear() - 1)).toISOString();
  return undefined;
}
```

### [x] 7. Update `Workouts/index.tsx`
**File:** `frontend/src/pages/Workouts/index.tsx`

- Replace `useWorkouts(page)` + `page` state with `useInfiniteWorkouts({ type, from })`
- Remove `page` and `totalPages` state
- Derive `workouts` by flattening pages: `data?.pages.flatMap(p => p.data) ?? []`
- Pass `filter !== "all" ? filter : undefined` as `type`
- Pass `periodToFrom(period)` as `from`
- Pass `fetchNextPage`, `hasNextPage`, `isFetchingNextPage` to `SessionSection`
- Keep `filter` and `period` state (they now drive server params)
- Remove `chartWorkouts`/`filtered` computed via period filtering — since `workouts` is already filtered by period, simplify chart/stats computations

### [x] 8. Update `SessionSection.tsx`
**File:** `frontend/src/pages/Workouts/SessionSection.tsx`

- Remove `page`, `totalPages`, `onPageChange` props and the prev/next pagination block
- Add `fetchNextPage`, `hasNextPage`, `isFetchingNextPage` props
- Add a sentinel `<div ref={sentinelRef} />` at the bottom of the list
- Wire `IntersectionObserver` in a `useEffect` to call `fetchNextPage` when sentinel enters the viewport:

```ts
const sentinelRef = useRef<HTMLDivElement>(null);
useEffect(() => {
  if (!sentinelRef.current) return;
  const obs = new IntersectionObserver(
    ([entry]) => { if (entry.isIntersecting && hasNextPage) fetchNextPage(); },
    { threshold: 0.1 }
  );
  obs.observe(sentinelRef.current);
  return () => obs.disconnect();
}, [hasNextPage, fetchNextPage]);
```

- Replace the `maxHeight: 640 / overflowY: auto` scroll container with normal flow (no artificial height cap)
- Show a small spinner or "Loading…" row when `isFetchingNextPage`
- Remove the `totalPages > 1` pagination footer entirely

---

## Critical Files

| File | Change |
|------|--------|
| `backend/internal/port/in/workout.go` | Add `WorkoutListParams`, update `List` interface |
| `backend/internal/port/out/repositories.go` | Update `WorkoutRepository.FindByUserID` signature |
| `backend/internal/usecase/workout_usecase.go` | Pass `WorkoutListParams` through |
| `backend/internal/adapter/in/http/handler/workout_handler.go` | Parse `type` + `from` query params |
| `backend/internal/adapter/out/mongo/workout_repository.go` | Apply type + date filter in MongoDB query |
| `frontend/src/hooks/useWorkouts.ts` | Add `useInfiniteWorkouts` + `periodToFrom` |
| `frontend/src/pages/Workouts/index.tsx` | Switch to infinite hook, remove pagination state |
| `frontend/src/pages/Workouts/SessionSection.tsx` | IntersectionObserver sentinel, remove prev/next |

---

## Verification

1. `cd backend && go build ./...` — confirms all interface changes compile
2. `cd backend && go test ./...` — ensures usecase/repo tests still pass
3. Start the cluster and open `localhost:30000/workouts`
4. Scroll to bottom — next 20 cards load automatically without clicking anything
5. Set type filter to "strength" — list resets and loads only strength workouts across all pages
6. Set period to "week" — list resets and loads only workouts from the last 7 days
7. Change period back to "all" — full list reloads correctly
8. Create a new workout — list invalidates and shows the new item at the top
