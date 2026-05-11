# Frontend Refactor Plan

## Context

The frontend has grown organically and now suffers from three core problems:
1. **Oversized page files** — Workouts.tsx is 1518 lines with 10+ inline sub-components. Dashboard (391), Settings (452), and all other pages have the same anti-pattern.
2. **Duplicated modal/form patterns** — every page hand-rolls its own modal shell, confirm dialog, and input styles. No shared primitives exist.
3. **Scattered inline styles** — 50-100+ `CSSProperties` objects per page with heavy duplication.

Goal: split every page into a feature folder, extract shared primitives, and centralize common styles — without changing behavior or visual output.

---

## Phase 1 — Foundation: Shared Primitives + Styles

### 1a. `frontend/src/lib/styles.ts`
Centralize all repeated inline style objects. Reuse throughout pages.

Key exports:
```ts
export const modalOverlay: CSSProperties   // fixed inset-0 dark overlay
export const modalCard: CSSProperties      // centered card shell
export const inputStyle: CSSProperties     // consistent form input
export const selectStyle: CSSProperties
export const formRow: CSSProperties        // label + input row
export const sectionDivider: CSSProperties
export const listItem: CSSProperties
```
Import design tokens from `lib/design.ts` (already has all colors, spacing, radius).

### 1b. `frontend/src/components/ui/FormModal.tsx`
Reusable modal shell used by all pages.

Props:
```ts
interface FormModalProps {
  title: string
  open: boolean
  onClose: () => void
  children: React.ReactNode
  footer?: React.ReactNode   // action buttons
  width?: number             // default 480
}
```
Replaces the ad-hoc modal overlays in Workouts, Finance, Sleep, Study, Events, Calendar, Settings.

### 1c. `frontend/src/components/ui/ConfirmDialog.tsx`
Reusable confirm dialog (currently duplicated in Workouts, Events, Finance).

Props:
```ts
interface ConfirmDialogProps {
  open: boolean
  title: string
  message: string
  confirmLabel?: string      // default "Delete"
  onConfirm: () => void
  onCancel: () => void
  danger?: boolean           // red confirm button
}
```

---

## Phase 2 — Split Workouts.tsx (1518 lines → feature folder)

**Target:** `frontend/src/pages/workouts/`

| File | Lines (est.) | Content |
|------|-------------|---------|
| `index.tsx` | ~120 | Page shell, list rendering, open/close state wiring |
| `WorkoutCard.tsx` | ~80 | Single workout row card |
| `WorkoutCreateForm.tsx` | ~200 | Create workout form in FormModal |
| `WorkoutEditForm.tsx` | ~200 | Edit workout form in FormModal |
| `WorkoutTracker.tsx` | ~350 | Active tracker (phases: plan/active/summary) |
| `useWorkoutPage.ts` | ~80 | Page-level state: editing ID, confirming delete, tracker open |

`App.tsx` import changes: `import WorkoutsPage from './pages/workouts'`

**State extraction into `useWorkoutPage.ts`:**
- `editingId`, `setEditingId`
- `confirmDeleteId`, `setConfirmDeleteId`
- `trackerOpen`, `setTrackerOpen`
- `activeWorkoutId`

---

## Phase 3 — Split Dashboard.tsx (391 lines → feature folder)

**Target:** `frontend/src/pages/dashboard/`

| File | Content |
|------|---------|
| `index.tsx` | Layout shell, composes section components |
| `StatsSummary.tsx` | 4 stat cards + period selector |
| `RecentActivitySection.tsx` | Activity feed list |
| `UpcomingEventsSection.tsx` | Upcoming events list |
| `GoalProgressSection.tsx` | Goal progress bars |
| `useDashboardData.ts` | All `useMemo` logic: sparklines, activity aggregation, event filtering, goal computation |

`useDashboardData.ts` takes raw hook data as inputs and returns derived display values — removes the useMemo chains from the JSX.

---

## Phase 4 — Split Settings.tsx (452 lines → feature folder)

**Target:** `frontend/src/pages/settings/`

| File | Content |
|------|---------|
| `index.tsx` | Section list layout |
| `ProfileCard.tsx` | Avatar, name, email display |
| `AppearanceCard.tsx` | Theme/color selector, extracted Toggle + ColorSwatch |
| `SystemCard.tsx` | Currency, language, week start |
| `NotificationsCard.tsx` | Notification toggles (uses shared Toggle) |
| `GoalsCard.tsx` | Default goal fields with Stepper |
| `PrivacyCard.tsx` | Data export/delete actions |
| `Toggle.tsx` | Extracted Toggle primitive (currently inline in Settings) |
| `Stepper.tsx` | Extracted Stepper primitive (currently inline in Settings) |

`Toggle.tsx` and `Stepper.tsx` go in `components/ui/` — they are generic enough to be shared.

---

## Phase 5 — Remaining Pages (Calendar, Events, Finance, Sleep, Study)

Each page gets a feature folder. Pattern is the same:

### Finance → `pages/finance/`
- `index.tsx`
- `TransactionItem.tsx`
- `CreateTransactionForm.tsx` (uses FormModal)
- `EditTransactionForm.tsx` (uses FormModal)

### Sleep → `pages/sleep/`
- `index.tsx`
- `SleepChart.tsx` (bar visualization, currently inline)
- `CreateSleepForm.tsx` (uses FormModal)

### Study → `pages/study/`
- `index.tsx`
- `SubjectPicker.tsx` (uses FormModal)
- `StudyTimer.tsx`

### Calendar → `pages/calendar/`
- `index.tsx`
- `CalendarGrid.tsx` (7-column grid layout)
- `EventModal.tsx` (create/edit/view modes, uses FormModal)
- `DayEventsPanel.tsx`

### Events → `pages/events/`
- `index.tsx`
- `AddEventForm.tsx` (uses FormModal)
- `EventListItem.tsx`

---

## Phase 6 — Style Centralization (applies across all phases)

As each file is written/rewritten, replace all local `const style: CSSProperties = { ... }` blocks with imports from `lib/styles.ts`.

Exceptions: unique one-off styles can stay local; only shared patterns (modal overlay, input, list item) move to `lib/styles.ts`.

---

## File Change Summary

**New files:**
- `frontend/src/lib/styles.ts`
- `frontend/src/components/ui/FormModal.tsx`
- `frontend/src/components/ui/ConfirmDialog.tsx`
- `frontend/src/components/ui/Toggle.tsx`
- `frontend/src/components/ui/Stepper.tsx`
- `frontend/src/pages/workouts/` (6 files)
- `frontend/src/pages/dashboard/` (6 files)
- `frontend/src/pages/settings/` (7 files)
- `frontend/src/pages/finance/` (4 files)
- `frontend/src/pages/sleep/` (3 files)
- `frontend/src/pages/study/` (3 files)
- `frontend/src/pages/calendar/` (4 files)
- `frontend/src/pages/events/` (3 files)

**Deleted files:**
- `frontend/src/pages/Workouts.tsx`
- `frontend/src/pages/Dashboard.tsx`
- `frontend/src/pages/Settings.tsx`
- `frontend/src/pages/Calendar.tsx`
- `frontend/src/pages/Events.tsx`
- `frontend/src/pages/Finance.tsx`
- `frontend/src/pages/Sleep.tsx`
- `frontend/src/pages/Study.tsx`

**Modified files:**
- `frontend/src/App.tsx` — update import paths
- `frontend/src/lib/styles.ts` — (new, exports common styles)

---

## Constraints

- No behavior changes — only structural refactoring
- No new API calls or state additions
- `Goals.tsx` is out of scope (API integration is a separate task)
- Do not remove or alter any existing hook logic in `hooks/`
- Existing components in `components/ui/` stay unchanged — only add new ones
- TypeScript must compile with no new errors after each phase

---

## Verification

After each phase:
1. `npm run build` in `frontend/` — must pass with 0 errors
2. `npm run lint:fix` — no new lint errors
3. Open the app in browser (port 30000 or `npm run dev`) and exercise the affected page:
   - Navigate, create a record, edit it, delete it (confirm dialog)
   - Check responsive layout (mobile bottom nav)
4. After Phase 2 (Workouts): complete a full tracker flow (plan → active → summary)
5. After Phase 3 (Dashboard): verify all 4 stat period selectors update sparklines
6. After Phase 4 (Settings): change theme, change currency, save profile

---

## Execution Order

- [ ] Phase 1 — Foundation (styles.ts, FormModal, ConfirmDialog, Toggle, Stepper)
- [ ] Phase 2 — Workouts split
- [ ] Phase 3 — Dashboard split
- [ ] Phase 4 — Settings split
- [ ] Phase 5 — Finance, Sleep, Study, Calendar, Events splits
- [ ] Phase 6 — Style centralization sweep (lib/styles.ts imports applied)
