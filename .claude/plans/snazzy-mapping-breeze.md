# Plan: Period Selector — Analysis Cards + Chart (Revised)

## Context

The workout page has a period selector (Week / Month / All) that currently only filters the workout list. The user wants:
1. **Stats cards** to show period-relevant metrics (not always this_week/this_month/streak)
2. **Chart** to show hours (not counts) in Week mode and visually distinguish planned vs started workouts

No backend changes. All data is computed from the existing stats API (`WorkoutStats`) + the fetched workout list (page 1, sorted newest-first).

---

## Files to modify

| File | Change |
|------|--------|
| `frontend/src/pages/Workouts.tsx` | New inline `PeriodBarChart`, new card design per period, new chart data logic |
| `frontend/src/i18n/locales/en.ts` | Add `workouts.stats.totalHours`, `workouts.stats.avgPerWeek`, `workouts.stats.completionRate`, `workouts.chartTitleMonth`, `workouts.chartTitleAll` |
| `frontend/src/i18n/locales/th.ts` | Same keys in Thai |

---

## 1. New inline `PeriodBarChart` component

The existing `BarChart` (shared lib) only supports a single value per bar. We need a dual-layer bar showing **planned** (full height, dimmed) vs **started/completed** (colored fill within).

Write a local `PeriodBarChart` in `Workouts.tsx` (no shared lib change):

```tsx
interface PeriodBar { label: string; planned: number; completed: number; }

function PeriodBarChart({ data, color, height = 64, unit }: {
  data: PeriodBar[];
  color: string;
  height?: number;
  unit: "h" | "count";
}) {
  if (!data.length) return null;
  const max = Math.max(...data.map((d) => d.planned), 1);
  const BAR_W = 28, GAP = 10;
  const n = data.length;
  const VW = n * (BAR_W + GAP) - GAP;
  const VH = height + 24;
  return (
    <div style={{ width: "100%", maxHeight: height + 40 }}>
      <svg width="100%" height="100%" viewBox={`0 0 ${VW} ${VH}`} preserveAspectRatio="xMidYMid meet" style={{ display: "block" }}>
        {data.map((d, i) => {
          const plannedH = Math.max(3, (d.planned / max) * height);
          const completedH = d.planned > 0 ? Math.max(0, (d.completed / d.planned) * plannedH) : 0;
          const x = i * (BAR_W + GAP);
          const fs = Math.min(10, (VW / n) * 0.38);
          return (
            <g key={i}>
              {/* planned (background, dimmed) */}
              <rect x={x} y={height - plannedH} width={BAR_W} height={plannedH} rx={4}
                fill={`color-mix(in srgb, ${color} 22%, transparent)`} />
              {/* completed (foreground, full color) */}
              {completedH > 0 && (
                <rect x={x} y={height - completedH} width={BAR_W} height={completedH} rx={4}
                  fill={color} style={{ transition: "all 0.3s" }} />
              )}
              <text x={x + BAR_W / 2} y={height + 15} textAnchor="middle" fontSize={fs}
                fill={C.text2} fontFamily="Inter,sans-serif">{d.label}</text>
            </g>
          );
        })}
      </svg>
    </div>
  );
}
```

---

## 2. Chart data per period

Build bars from `workouts` (full page, unfiltered by type so bars reflect all activity):

**Week** — 7 bars (days), unit = hours:
```ts
bars = last 7 days, each bar:
  label: weekday short (Mon, Tue…)
  planned: sum(w.duration) for workouts on that day / 60  (hours, 1 decimal)
  completed: sum(w.actual_duration) for workouts on that day with actual_duration / 60
```

**Month** — 4 bars (weeks), unit = count:
```ts
bars = last 4 weeks (oldest → newest), each bar:
  label: "Wk 1", "Wk 2", "Wk 3", "Wk 4"
  planned: workouts created in that week's date range
  completed: workouts in range that have actual_duration > 0
```

**All** — 6 bars (months), unit = count:
```ts
bars = last 6 months (oldest → newest), each bar:
  label: month short (Jan, Feb…)
  planned: workouts created in that calendar month
  completed: workouts in month with actual_duration > 0
```

Chart title:
- week → `t("workouts.chartTitle")` ("This Week")
- month → `t("workouts.chartTitleMonth")` ("Last 4 Weeks")
- all → `t("workouts.chartTitleAll")` ("Last 6 Months")

---

## 3. Stats cards per period

Compute from `filtered` (period + type filtered list) and `stats` API:

| Period | Card 1 | Card 2 | Card 3 |
|--------|--------|--------|--------|
| **week** | Sessions: `stats.this_week` · "This Week" | Planned hours: `sum(filtered.duration)/60` · "Hours planned" | Streak 🔥 |
| **month** | Sessions: `stats.this_month` · "This Month" | Avg/week: `stats.this_month / 4 \| 1 decimal` · "Avg / week" | Streak 🔥 |
| **all** | Sessions: `stats.total` · "All Time" | Total hours: `sum(workouts.duration)/60` · "Hours logged" | Streak 🔥 |

Hours computed client-side as `(sum / 60).toFixed(1)`.

---

## 4. i18n additions

**en.ts** under `workouts.stats`:
```ts
hoursPlanned: "Hours planned",
avgPerWeek: "Avg / week",
hoursLogged: "Hours logged",
allTime: "All Time",
```

**en.ts** under `workouts`:
```ts
chartTitleMonth: "Last 4 Weeks",
chartTitleAll: "Last 6 Months",
```

Same for **th.ts**:
```ts
hoursPlanned: "ชั่วโมงที่วางแผน",
avgPerWeek: "เฉลี่ย / สัปดาห์",
hoursLogged: "ชั่วโมงที่ทำไป",
allTime: "ทั้งหมด",
chartTitleMonth: "4 สัปดาห์ที่ผ่านมา",
chartTitleAll: "6 เดือนที่ผ่านมา",
```

---

## Verification

1. **Week** → cards show this-week count / planned-hours / streak; chart = 7 day bars in hours with dimmed planned + colored completed portion
2. **Month** → cards show this-month count / avg per week / streak; chart = 4 weekly count bars
3. **All** → cards show total count / total hours / streak; chart = 6-month count bars
4. Switching type filter updates the cards' hours/avg but not the chart bars
5. `npm run build` — no TypeScript errors
