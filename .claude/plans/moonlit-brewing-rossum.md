# Plan: E2E Tests for Workout Page

## Context

The workout page already has Playwright infrastructure (fixtures, mocks, page-object pattern) used by the finance tests. The workout module has test IDs defined in `src/constants/testIds/workouts.ts` but no e2e tests yet. This plan adds workout e2e coverage following the exact same patterns as the finance tests.

---

## Test cases (3)

1. **renders seeded workouts** — navigate to `/workouts`, assert heading visible, log/start buttons visible, seeded workout cards present
2. **logs a cardio workout** — click "Log Workout", select Cardio type, pick 30 min duration, save → assert new card appears in list
3. **deletes a workout** — click Delete on a seeded card, confirm → assert card is gone from list

---

## Files to create

### `e2e/pages/workouts.page.ts`
Page object following `FinancePage` pattern:
- `goto()` → navigate to `/workouts`
- `heading()` / `list()` / `logButton()` / `startWorkoutButton()` → by `WORKOUTS_TEST_IDS.*`
- `workoutCard(id)` → by `workoutCardTestId(id)`
- `createForm()` → by `WORKOUTS_TEST_IDS.createForm`
- `selectWorkoutType(type)` → within form, click button with text `"Strength" | "Cardio" | "Flexibility"`
- `selectDuration(mins)` → within form, click button with text `"{mins} min"`
- `save()` → click Save button within form
- `deleteWorkoutButton(id)` → within card, click the "Delete" button
- `confirmDeleteButton()` → page-level "Delete" button in the confirm dialog (rendered outside card)
- `logWorkout({ type, durationMins })` — composite action

### `e2e/tests/workouts.spec.ts`
```
test.describe("Workouts") {
  test("renders seeded workouts") { goto, assert heading/list/cards }
  test("logs a cardio workout") { logWorkout → assert card }
  test("deletes a workout") { delete first seeded workout → assert card gone }
}
```

---

## Files to modify

### `e2e/mocks/mock-api.ts`
Add to `MockAppState`:
```ts
workouts: Workout[];
nextWorkoutId: number;
workoutStats: WorkoutStats;
```

Add default seeded workouts (2): one strength, one cardio.

Add route handlers:
- `GET /workouts` → `PagedResponse<Workout>` from `state.workouts`
- `GET /workouts/stats` → `state.workoutStats` (recalculated on the fly)
- `POST /workouts` → create from body, push to `state.workouts`, return 201
- `DELETE /workouts/:id` → filter out, return 204

Update `isHandledPath()` to include `/workouts` and `/workouts/stats`.

Export `DEFAULT_WORKOUTS` array so tests can reference seeded IDs.

### `e2e/fixtures/test.ts`
Add `workoutsPage: WorkoutsPage` to the `Fixtures` type and register it:
```ts
workoutsPage: async ({ page, mockAppState }, use) => {
  void mockAppState;
  await use(new WorkoutsPage(page));
},
```

---

## Key implementation notes

- Duration preset buttons render as `"{mins} min"` (e.g. `"30 min"`) — use `getByRole("button", { name: "30 min" })` scoped to the `createForm` container
- Type selector buttons are plain `<button type="button">` with text "Strength" / "Cardio" / "Flexibility"
- Delete confirm dialog is rendered as a sibling to the card (not inside it), so `confirmDeleteButton()` is a page-level locator not scoped to the card
- Workout cards use `workoutCardTestId(id)` → `"workouts-card-{id}"`

---

## Verification

```powershell
cd frontend
npm run test:e2e -- --grep "Workouts"
# or headed for visual check:
npm run test:e2e:headed -- --grep "Workouts"
```
