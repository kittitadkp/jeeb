import { C, R, T, W } from "../utils/design";

export type Period = "day" | "week" | "month" | "year";

interface PeriodSelectorProps {
  value: Period;
  onChange: (p: Period) => void;
  accent?: string;
}

const PERIODS: Period[] = ["day", "week", "month", "year"];

export function PeriodSelector({ value, onChange, accent }: PeriodSelectorProps) {
  const col = accent || C.primary;
  return (
    <div style={{ display: "flex", gap: 2, background: C.surface2, borderRadius: R.full, padding: 3 }}>
      {PERIODS.map((p) => (
        <button
          key={p}
          onClick={() => onChange(p)}
          style={{
            padding: `${T.xs - 7}px 11px`,
            borderRadius: R.full,
            fontSize: T.xs,
            fontWeight: W.semi,
            background: value === p ? col : "transparent",
            color: value === p ? "#fff" : C.text2,
            border: "none",
            transition: "all 0.15s",
            textTransform: "capitalize",
            letterSpacing: "0.02em",
            cursor: "pointer",
          }}
        >
          {p}
        </button>
      ))}
    </div>
  );
}
