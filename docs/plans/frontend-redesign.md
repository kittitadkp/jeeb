# Frontend Redesign — Dark Premium Theme

> **Design source**: `jeeb-redesign.html` (from claude.ai/design handoff bundle)
> **Status**: Complete ✅

## Overview

Full redesign of all 9 pages to a dark-first premium aesthetic. Per-section accent colors, Space Grotesk headings, SVG charts on every data page, and several new UX features (Quick-Log FAB, Notifications panel, Goals & Streaks, Bulk Log).

---

## Design Tokens (what changes from current)

| Token | Current | New |
|---|---|---|
| Page bg | `slate-50` | `#090b10` |
| Card surface | `white` | `#0f1117` |
| Card border | `slate-200` | `rgba(255,255,255,0.06)` |
| Heading font | Inter | Space Grotesk + Inter |
| Primary accent | blue-600 | indigo (dashboard) + per-section |

### Per-section accent colors

| Section | Color |
|---|---|
| Dashboard | Indigo |
| Workouts | Rose |
| Study | Amber |
| Sleep | Violet |
| Finance | Emerald |
| Calendar | Teal |
| Goals | Cyan |
| Events | Purple |

---

## Implementation Checklist

### Phase 1 — Design tokens & theme

- [x] Add Space Grotesk to `index.html` (Google Fonts)
- [x] Update `src/lib/design.ts` — add section accent map, dark surface tokens
- [x] Update `src/store/theme.ts` — default to dark mode
- [x] Update `section-label.tsx` — Space Grotesk heading font

### Phase 2 — Layout shell

- [x] **Sidebar** (`src/components/layout/Sidebar.tsx`)
  - Gradient "Jeeb" wordmark
  - Per-section colored left-border active indicator
  - Section emoji icons (💪 📚 😴 💰 📅 🎯 📌)
  - User avatar at bottom
  - "📋 Bulk Log" button above avatar
- [x] **AppLayout** (`src/components/layout/AppLayout.tsx`)
  - Wired `bulkLogOpen` state → passes `onBulkLog` to Sidebar
  - Renders `<BulkLogModal>` when open

### Phase 3 — Shared UI components

- [x] **BulkLogModal** — spreadsheet-style catch-up log
- [x] **QuickLogFAB** — date chips (Today / Yesterday) on Study, Sleep, Finance forms

### Phase 4 — Pages

#### Dashboard (`src/pages/Dashboard.tsx`)
- [x] 4 stat cards in a row with section accent strips and sparklines
- [x] Period selector (Day / Week / Month / Year)
- [x] Goals & Streaks panel — 4 goal progress bars + streak badges

#### Workouts (`src/pages/Workouts.tsx`)
- [x] Rose accent throughout
- [x] Weekly bar chart (current day highlighted)
- [x] Filter pills (All / Strength / Cardio / Flexibility)

#### Study (`src/pages/Study.tsx`)
- [x] Amber accent
- [x] Weekly bar chart
- [x] Live study timer (functional)

#### Sleep (`src/pages/Sleep.tsx`)
- [x] Violet accent

#### Finance (`src/pages/Finance.tsx`)
- [x] Emerald accent
- [x] Weekly spending bar chart + donut breakdown chart side-by-side

#### Calendar (`src/pages/Calendar.tsx`)
- [x] Teal accent

#### Goals (`src/pages/Goals.tsx`)
- [x] Goal cards with colored progress bars
- [x] Category badge + period badge on each card
- [x] Summary stats: Total / In progress / Completed
- [x] Period chips in modal (Today / This week / This month / This year / Custom)

#### Events (`src/pages/Events.tsx`)
- [x] Overall completion progress bar
- [x] Quick stats: Due today / Upcoming / Completed
- [x] Date subtitle on each event item
- [x] Filter tabs: All / Today / Upcoming / Completed

#### Settings (`src/pages/Settings.tsx`)
- [x] Profile card — gradient avatar, name/email
- [x] Appearance card — Dark/Light/System, 6-color accent picker, week start
- [x] Notifications card — per-category toggles with time inputs
- [x] Default Goal Targets card — steppers for workouts/study/sleep, budget input
- [x] Data & Privacy card — Export JSON, Sync Now, Delete Account, Log out

### Phase 5 — QuickLogModal date chips

- [x] Workout form: Today / Yesterday chips
- [x] Study form: Today / Yesterday chips
- [x] Sleep form: Today / Yesterday chips
- [x] Finance form: Today / Yesterday chips

### Phase 6 — Bulk Log modal

- [x] Spreadsheet-style catch-up log
- [x] Pick category → grid of last 7 days
- [x] Toggle rows on/off, edit fields inline
- [x] "Save All" sends all toggled entries
- [x] Accessible via sidebar "📋 Bulk Log" button

---

## File map (files created or heavily modified)

```
src/
  index.css                            ✅ Space Grotesk Google Fonts
  lib/design.ts                        ✅ section accent colors, dark tokens
  store/theme.ts                       ✅ default dark
  components/
    layout/
      AppLayout.tsx                    ✅ bulk log modal wired
      Sidebar.tsx                      ✅ gradient logo, section accents, bulk log btn
    ui/
      section-label.tsx                ✅ Space Grotesk font
      bulk-log-modal.tsx               ✅ NEW — 7-day spreadsheet UI
  pages/
    Dashboard.tsx                      ✅ period selector, goals & streaks
    Workouts.tsx                       ✅ rose accent, weekly bar chart
    Study.tsx                          ✅ amber accent, weekly bar chart
    Finance.tsx                        ✅ emerald accent, dual charts
    Goals.tsx                          ✅ period field in form, badges on cards
    Events.tsx                         ✅ completion bar, date subtitles
    Settings.tsx                       ✅ 5-section redesign
```
