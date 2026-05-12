# Plan: Refactor Workouts page into folder with split section components

## Context
`Workouts.tsx` is 1518 lines — a single monolithic file containing styles, utilities, sub-components, and the page component. The goal is to split it into a `Workouts/` folder where each visible page section is its own file, making the code easier to navigate and maintain. The session list section also needs a scrollable container so it scrolls independently within the page.

---

## Target folder structure

```
frontend/src/pages/Workouts/
  index.tsx               # Main Workouts page component (state + layout)
  styles.ts               # All shared CSSProperties constants + ACC
  utils.ts                # serializeExerciseEntries, filterWorkoutsByPeriod, fmtTime, calculateLoggedActualDuration
  shared.tsx              # HoverableButton, ConfirmDialog, WorkoutFormSection, WorkoutTypeSelector, WorkoutDurationPresets, PeriodBarChart
  WorkoutForm.tsx         # CreateWorkoutForm + EditWorkoutForm
  WorkoutCard.tsx         # WorkoutCard component
  WorkoutTracker.tsx      # WorkoutTracker modal
  HeroSection.tsx         # Top card: title, action buttons, inline stat tiles
  StatsSection.tsx        # Period-aware 3-stat row
  ChartSection.tsx        # Period-aware bar chart
  ControlsSection.tsx     # Filter chips + period toggle
  SessionSection.tsx      # Session list with scroll container + pagination
```

---

## Step-by-step

### 1. Create `styles.ts`
Extract all `const *Style: React.CSSProperties` variables and `ACC`, `WORKOUT_DURATION_PRESETS` into this file. Export each.

Key exports:
- `ACC`, `inputStyle`, `numInputStyle`, `textareaStyle`
- `workoutFormOverlayStyle`, `workoutFormModalStyle`, `workoutFormSectionStyle`
- `workoutFormGridStyle`, `workoutFormTypeGridStyle`, `workoutFormLabelStyle`
- `workoutFormHintStyle`, `workoutFormPresetRowStyle`, `workoutFormFooterStyle`
- `workoutFormToggleButtonStyle`, `WORKOUT_DURATION_PRESETS`

### 2. Create `utils.ts`
Extract pure functions:
- `serializeExerciseEntries`
- `filterWorkoutsByPeriod`
- `fmtTime`
- `calculateLoggedActualDuration`

### 3. Create `shared.tsx`
Extract shared sub-components used across multiple files:
- `HoverableButton`
- `ConfirmDialog`
- `WorkoutFormSection`
- `WorkoutTypeSelector`
- `WorkoutDurationPresets`
- `PeriodBarChart` (+ `PeriodBar` interface)

### 4. Create `WorkoutForm.tsx`
Move `CreateWorkoutForm` and `EditWorkoutForm`. Imports from `styles.ts`, `utils.ts`, `shared.tsx`.

### 5. Create `WorkoutCard.tsx`
Move `WorkoutCard` component. Imports `EditWorkoutForm` from `WorkoutForm.tsx`, `ConfirmDialog`/`HoverableButton` from `shared.tsx`.

### 6. Create `WorkoutTracker.tsx`
Move `WorkoutTracker` component (lines 448–892). Imports from `styles.ts`, `utils.ts`, `shared.tsx`.

### 7. Create `HeroSection.tsx`
Extract the hero `<Card>` block from `Workouts()` (title + action buttons + 3 inline stat tiles).

Props:
```ts
{ filteredPlannedHours: string; filteredLoggedHours: string; streak: number | string;
  onLog: () => void; onPlan: () => void; }
```

### 8. Create `StatsSection.tsx`
Extract the period-aware 3-card stats row.

Props:
```ts
{ period: WorkoutPeriod; stats: WorkoutStats | undefined;
  chartWeekWorkouts: Workout[]; chartMonthWorkouts: Workout[];
  chartYearWorkouts: Workout[]; chartWorkouts: Workout[];
  filteredPlannedHours: string; chartLoggedHours: string; }
```

### 9. Create `ChartSection.tsx`
Extract the bar chart `<Card>` block.

Props:
```ts
{ period: WorkoutPeriod; chartWorkouts: Workout[]; formatDate: (...) => string; }
```
`PeriodBarChart` imported from `shared.tsx`.

### 10. Create `ControlsSection.tsx`
Extract the filter chips + period toggle row.

Props:
```ts
{ filter: WorkoutFilter; period: WorkoutPeriod;
  onFilterChange: (f: WorkoutFilter) => void;
  onPeriodChange: (p: WorkoutPeriod) => void; }
```

### 11. Create `SessionSection.tsx`
Extract the session list + pagination. **Add a scrollable container** around the card list:

```tsx
<div style={{ maxHeight: 640, overflowY: "auto", display: "flex", flexDirection: "column", gap: 12, paddingRight: 4 }}>
  {filtered.map(...)}
</div>
```

Props:
```ts
{ filtered: Workout[]; filteredCompleted: number;
  isLoading: boolean; isError: boolean; refetch: () => void;
  trackerWorkoutId: string | null;
  onDelete: (id: string) => void; onStart: (w: Workout) => void;
  page: number; totalPages: number;
  onPageChange: (p: number) => void; }
```

### 12. Create `index.tsx`
Rewrite `Workouts()` to only hold state + computed values, then compose sections:

```tsx
export function Workouts() {
  // ... state, hooks, computed values (same as current)
  return (
    <div ...>
      <HeroSection ... />
      {createMode && <CreateWorkoutForm ... />}
      <StatsSection ... />
      <ChartSection ... />
      <ControlsSection ... />
      <SessionSection ... />
      {trackerWorkout && <WorkoutTracker ... />}
    </div>
  );
}
```

### 13. Update import in router/app
`Workouts.tsx` is currently imported by path. After the refactor the export moves to `Workouts/index.tsx` — the import path `@/pages/Workouts` resolves to the folder index automatically so **no router change needed** as long as the named export `Workouts` is kept.

Verify the import location:
- `frontend/src/App.tsx` or router file — check that `import { Workouts } from "@/pages/Workouts"` still resolves.

### 14. Delete original `Workouts.tsx`
Once `Workouts/index.tsx` exists and compiles, remove the old flat file.

---

## Verification
1. `cd frontend && npm run build` — must pass TypeScript + Vite build with no errors.
2. Open `http://localhost:30000` (or `npm run dev`), navigate to Workouts page — all sections render.
3. Session list scrolls independently when more than ~5 cards fill the container.
4. Log/Plan buttons open the form modal.
5. Start workout opens the tracker modal.
