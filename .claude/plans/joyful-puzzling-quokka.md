# Plan: Workout Page Recommendations

## Context

Audit of the workout page after recent major changes (infinite scroll, server-side filtering, status filter, completed-state design, layout restructure). Goal is to identify what to fix, improve, or build next.

---

## High Priority — UX Breaking

### 1. Workout type color mapping is dead code
**File:** `frontend/src/pages/Workouts/WorkoutCard.tsx:24`
```ts
const typeColors = { strength: C.primary, cardio: C.primary, flexibility: C.primary };
```
All three map to the same color. The badge and icon color (`stateColor`) are identical for all types, making the type distinction purely textual. Either define distinct colors per type (e.g., orange for cardio, teal for flexibility) or delete `typeColors` and let `badgeColor = ACC` always — don't pretend there's differentiation.

**Fix:** Define a proper `WORKOUT_TYPE_COLORS` constant in `constants/` with distinct values and reference it in `WorkoutCard` and `WorkoutTypeSelector`.

---

### 2. EditWorkoutForm has duplicate and unsynced duration inputs
**File:** `frontend/src/pages/Workouts/WorkoutForm.tsx:86–99`

`WorkoutDurationPresets` (presets + number input) and a second `<input type="number">` for duration both render. They bind to the same `form.duration` but the label appears twice as "Duration (minutes)". The second input is redundant — `WorkoutDurationPresets` already includes a number input at the bottom.

**Fix:** Remove the standalone `<input type="number" ...value={form.duration}...>` at line 93 from `EditWorkoutForm`. Keep only `WorkoutDurationPresets`.

---

### 3. Session count is misleading with infinite pagination
**File:** `frontend/src/pages/Workouts/SessionSection.tsx`

The count shows `filtered.length` (number of loaded cards) not the total across all pages. With page size 20, a user with 45 strength sessions filtered to "strength" sees "12 completed, 20 sessions" but there are more on the next pages.

**Fix:** Thread `meta.total` from the last loaded page through to `SessionSection`. Show "X loaded" or use the total from `data.pages.at(-1)?.meta.total` as the denominator.

---

### 4. Empty state doesn't distinguish "no data" vs "no filter results"
**File:** `frontend/src/pages/Workouts/SessionSection.tsx`

When filters are active and return 0 results, the empty state shows the generic message. Users don't know if they have no workouts or if filters caused the empty result.

**Fix:** Pass a `hasActiveFilters` boolean (any filter is non-"all") from `index.tsx` to `SessionSection`. Show "No workouts matching your filters" + a "Clear filters" button when filters are active.

---

### 5. Mutation errors are silently swallowed
**Files:** `WorkoutForm.tsx`, `WorkoutCard.tsx`

`useCreateWorkout`, `useUpdateWorkout`, `useDeleteWorkout` — none have `onError` handlers. If the API returns 4xx/5xx, the form just stops loading with no feedback.

**Fix:** Add `onError` to each mutation call. Show an inline error message below the form footer or use a simple toast/banner. At minimum log to console and display a generic "Something went wrong" message.

---

## Medium Priority — Correctness

### 6. Stats math is inaccurate
**File:** `frontend/src/pages/Workouts/StatsSection.tsx`

- `avgPerWeek = chartMonthWorkouts.length / 4` — hardcoded 4 weeks, wrong for months with 5 weeks or partial months.
- `yearAvgPerMonth = chartYearWorkouts.length / (now.getMonth() + 1)` — biased early in the year (Jan shows 1 as divisor even if data spans multiple months).

**Fix:**
- `avgPerWeek`: count actual calendar weeks in the period using `differenceInCalendarWeeks` or compute weeks since `periodToFrom("month")`.
- `yearAvgPerMonth`: use `Math.max(1, now.getMonth() + 1)` and only count months that have at least one data point, or just show "this year" totals without the misleading per-month average.

---

### 7. Notes are never displayed on WorkoutCard
**File:** `frontend/src/pages/Workouts/WorkoutCard.tsx`

`workout.notes` exists in the domain but nothing renders it. Users who add notes (via form) see no evidence of them on the card.

**Fix:** Add a notes row below the exercises list — only when `workout.notes` is truthy. Truncate to 80 chars with a "…" overflow and full note on hover/expand.

---

### 8. WorkoutTracker goal edits are lost on close
**File:** `frontend/src/pages/Workouts/WorkoutTracker.tsx`

The tracker allows editing sets/reps during a session, but these edits are not persisted to the backend unless the user explicitly saves. Closing the tracker silently drops them.

**Fix:** Either (a) auto-save exercise changes to the backend as the user edits, or (b) show a "Discard changes?" confirmation in the tracker's close handler when exercises have been modified.

---

## Low Priority — Polish

### 9. Edit form accessible while workout is being tracked
**File:** `frontend/src/pages/Workouts/WorkoutCard.tsx:144`

The edit button is always shown, even when `isTracking=true`. Editing a workout mid-session could cause confusion if the tracker and form are out of sync.

**Fix:** Disable or hide the edit button when `isTracking=true`. Already blocked for "Start" — extend the same logic to edit.

---

### 10. Chart has no accessibility labels
**File:** `frontend/src/pages/Workouts/shared.tsx` (`PeriodBarChart`)

SVG bars and text have no ARIA roles, titles, or descriptions. Screen readers see nothing meaningful.

**Fix:** Add `role="img"` and `<title>` to the `<svg>`. Each bar group can have `aria-label={`${d.label}: ${d.completed} completed, ${d.planned} planned`}`.

---

## Summary Table

| # | Priority | File | Change |
|---|----------|------|--------|
| 1 | High | `WorkoutCard.tsx` | Define distinct colors per workout type |
| 2 | High | `WorkoutForm.tsx` | Remove duplicate duration input in EditWorkoutForm |
| 3 | High | `SessionSection.tsx` | Show total count from API meta, not just loaded count |
| 4 | High | `SessionSection.tsx` + `index.tsx` | "No filter results" empty state + clear filters |
| 5 | High | `WorkoutForm.tsx`, `WorkoutCard.tsx` | Add `onError` handlers to all mutations |
| 6 | Medium | `StatsSection.tsx` | Fix avgPerWeek and yearAvgPerMonth math |
| 7 | Medium | `WorkoutCard.tsx` | Display notes field on card |
| 8 | Medium | `WorkoutTracker.tsx` | Persist or confirm discard of goal edits |
| 9 | Low | `WorkoutCard.tsx` | Disable edit button when isTracking |
| 10 | Low | `shared.tsx` | Add ARIA labels to PeriodBarChart |

---

## Suggested order of implementation

1. Fix #2 (duplicate input) — 5 min, no risk
2. Fix #5 (mutation errors) — low effort, high user impact
3. Fix #4 (filter empty state) — clear UX win
4. Fix #3 (session count) — requires threading meta through props
5. Fix #7 (show notes) — visible value, low effort
6. Fix #1 (type colors) — requires design decision on color palette
7. Fix #6 (stats math) — correctness improvement
8. Fix #8 (tracker persistence) — requires backend call or UX flow decision
9. Fix #9 and #10 — polish, low urgency
