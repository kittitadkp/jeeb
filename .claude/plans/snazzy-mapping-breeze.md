# Plan: Workout Feature Enhancements

## Context

Extending the workout feature with per-exercise duration tracking, plan vs actual duration visibility on cards, type-aware forms (cardio/flexibility skip exercise picking), a confirm dialog before starting the tracker, and a simpler settings page that only manages strength exercises.

---

## Changes

### 1. Backend — `backend/internal/domain/workout.go`

Add two fields:

```go
type Exercise struct {
    Name            string  `bson:"name" json:"name"`
    Sets            int     `bson:"sets" json:"sets"`
    Reps            int     `bson:"reps" json:"reps"`
    Weight          float64 `bson:"weight" json:"weight"`
    RestSeconds     int     `bson:"rest_seconds" json:"rest_seconds"`
    DurationSeconds int     `bson:"duration_seconds" json:"duration_seconds"` // planned per-exercise duration
}

type Workout struct {
    ...
    Duration       int `bson:"duration" json:"duration"`               // planned (user-entered, minutes)
    ActualDuration int `bson:"actual_duration" json:"actual_duration"` // tracker-recorded (minutes)
    ...
}
```

No changes needed to usecase or handler — they pass structs through generically.

---

### 2. Frontend types — `frontend/src/types/index.ts`

```ts
export interface Exercise {
  name: string;
  sets: number;
  reps: number;
  weight: number;
  rest_seconds: number;
  duration_seconds?: number;   // optional, 0 = not set
}

export interface Workout {
  ...
  duration: number;          // plan (minutes)
  actual_duration?: number;  // tracker time (minutes), absent on older records
  ...
}
```

---

### 3. ExercisePicker — `frontend/src/components/ExercisePicker.tsx`

Add `duration_seconds` to `ExerciseFormEntry` (default `0`).

In the selected-exercise form, add a 4th input "Duration (min)" after Sets/Reps/Rest:

```tsx
{ label: "Duration (min)", field: "duration_seconds" as const, min: 0 }
```

Show `duration_seconds / 60` in the input (store as seconds, display as minutes — convert on read/write). Or just store/display in seconds and label it "sec". Keep consistent with existing `rest_seconds`.

Actually: store as seconds (field name is `duration_seconds`), input label "Duration (s)" for consistency.

---

### 4. Settings — `frontend/src/constants/master.ts`

Add `fixedDefaults?: Record<string, string | number>` to `MasterCategoryConfig`.

Remove `exercise_type` from `defaultFields` (settings is strength-only — no need to select type). Add `fixedDefaults: { exercise_type: "strength" }` so every created/displayed exercise is strength.

```ts
export interface MasterCategoryConfig {
  slug: string;
  category: string;
  label: string;
  defaultFields: MasterDefaultField[];
  fixedDefaults?: Record<string, string | number>;
}

// config:
{
  slug: "workout",
  category: "exercise",
  label: "Exercises",
  fixedDefaults: { exercise_type: "strength" },
  defaultFields: [
    { key: "muscle_group", label: "Muscle Group", type: "string", defaultValue: "" },
    { key: "sets", label: "Sets", type: "number", defaultValue: 3 },
    { key: "reps", label: "Reps", type: "number", defaultValue: 10 },
    { key: "rest_seconds", label: "Set Rest (s)", type: "number", defaultValue: 60 },
    { key: "transition_rest_seconds", label: "Transition Rest (s)", type: "number", defaultValue: 30 },
  ],
}
```

---

### 5. MasterDataCard — `frontend/src/components/MasterDataCard.tsx`

Two changes:
- **Filter**: only display records where every `fixedDefaults` key matches (e.g., `record.defaults.exercise_type === "strength"`)
- **Save**: merge `config.fixedDefaults` into `defaults` before calling create/update

```ts
// filter
const visible = records.filter((r) =>
  !config.fixedDefaults ||
  Object.entries(config.fixedDefaults).every(([k, v]) => String(r.defaults[k] ?? "") === String(v))
);

// save
const defaults = {
  ...Object.fromEntries(config.defaultFields.map((f) => [f.key, formDefaults[f.key]])),
  ...config.fixedDefaults,
};
```

---

### 6. Workouts page — `frontend/src/pages/Workouts.tsx`

#### 6a. WorkoutTracker — simplify props

Remove `initialExercises`, `autoStart`, `onSave`. Add `workout: Workout`.

```ts
function WorkoutTracker({ workout, onClose }: {
  workout: Workout;
  onClose: () => void;
})
```

- Initialize directly in `"active"` phase (always autoStart from card)
- Initialize `trackerExs` from `workout.exercises`
- Initialize `workoutType` from `workout.type`
- On save: call `useUpdateWorkout({ id: workout.id, actual_duration: duration })` internally (add `useUpdateWorkout` hook inside tracker)
- Remove `workoutType` state (derived from `workout.type`)

#### 6b. WorkoutCard — confirm dialog + plan/actual duration

**Header area** — show plan and actual duration:
```tsx
<div style={{ fontSize: 11, color: C.text2 }}>
  {workout.actual_duration
    ? `Plan: ${workout.duration}min · Actual: ${workout.actual_duration}min`
    : `${workout.duration}min`}
</div>
```

**Start button** — show confirm dialog before opening tracker:

Add local state `[confirming, setConfirming] = useState(false)`.

When `confirming`:
```tsx
// small overlay modal
<div style={modal}>
  <div style={container}>
    <div>Start "{t(`common.workoutTypes.${workout.type}`)}" workout?</div>
    <div>exercises list...</div>
    <Button onClick={() => { setConfirming(false); onStart(workout); }}>▶ Start</Button>
    <Button variant="ghost" onClick={() => setConfirming(false)}>Cancel</Button>
  </div>
</div>
```

`onStart` prop replaces `onPlanAgain` — receives the full `Workout` object.

**Exercises list** — show `duration_seconds` if set:
```
• Bench Press  3×10 @ 60kg  (120s)
```

**Start button** — only shown when `workout.type === "strength"` AND exercises exist.

#### 6c. CreateWorkoutForm — type-aware

Only show ExercisePicker when `form.type === "strength"`:
```tsx
{form.type === "strength" && (
  <div>
    <button type="button" onClick={() => setShowExercises((v) => !v)}>...</button>
    {showExercises && <ExercisePicker ... />}
  </div>
)}
```

Pass `duration_seconds` in exercise mapping:
```ts
const apiExercises = exercises.map((ex) => ({
  name: ex.name, sets: ex.sets, reps: ex.reps,
  rest_seconds: ex.rest_seconds, weight: ex.weight,
  duration_seconds: ex.duration_seconds ?? 0,
}));
```

#### 6d. EditWorkoutForm — same type-aware logic

Same as 6c. Only show exercises section when `form.type === "strength"`. Initialize exercises from `workout.exercises` (map `duration_seconds`).

#### 6e. Workouts component — simplified tracker state

Replace `trackerInitialExercises`, `trackerAutoStart` states with single `trackerWorkout: Workout | null`.

```ts
const [trackerWorkout, setTrackerWorkout] = useState<Workout | null>(null);
```

```tsx
{trackerWorkout && (
  <WorkoutTracker
    workout={trackerWorkout}
    onClose={() => setTrackerWorkout(null)}
  />
)}
```

Pass `onStart` to WorkoutCard:
```tsx
<WorkoutCard
  ...
  onStart={(w) => setTrackerWorkout(w)}
/>
```

---

## Files to modify

| File | Change |
|------|--------|
| `backend/internal/domain/workout.go` | Add `DurationSeconds` to Exercise, `ActualDuration` to Workout |
| `frontend/src/types/index.ts` | Add `duration_seconds?` to Exercise, `actual_duration?` to Workout |
| `frontend/src/components/ExercisePicker.tsx` | Add `duration_seconds` to form entry + UI input |
| `frontend/src/constants/master.ts` | Add `fixedDefaults` to interface; remove `exercise_type` field; set `fixedDefaults: { exercise_type: "strength" }` |
| `frontend/src/components/MasterDataCard.tsx` | Filter records by `fixedDefaults`; merge `fixedDefaults` on save |
| `frontend/src/pages/Workouts.tsx` | Tracker simplification; confirm dialog; type-aware forms; plan/actual duration display |

---

## Verification

1. `backend`: `go build ./cmd/api/...` — no compile errors
2. Settings `/settings/workout` — only strength exercises shown; new exercises auto-tagged as strength
3. Create workout (cardio/flexibility) — exercise picker hidden; only type + duration + notes
4. Create workout (strength) — exercise picker shown with duration_seconds input
5. WorkoutCard strength — shows "Plan: Xmin", Start button visible
6. Start button — confirm dialog appears with exercise list → Start opens tracker → End → card shows "Plan: Xmin · Actual: Ymin"
7. WorkoutCard cardio/flexibility — no Start button, shows only duration
