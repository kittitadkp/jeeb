import { C, R, S, T, W } from "@/lib/design";

interface StatCardProps {
  label: string;
  value: number | string;
  color?: string;
}

export function StatCard({ label, value, color = C.primary }: StatCardProps) {
  return (
    <div style={{
      background: C.surface,
      border: `1px solid ${C.border}`,
      borderRadius: R.card,
      padding: `${S[4]}px ${S[5]}px`,
      display: "flex",
      flexDirection: "column",
      gap: 4,
    }}>
      <div style={{ fontSize: T.xs, color: C.text2, fontWeight: W.medium, textTransform: "uppercase", letterSpacing: "0.06em" }}>{label}</div>
      <div style={{ fontSize: T["2xl"], fontWeight: W.bold, color }}>{value}</div>
    </div>
  );
}
