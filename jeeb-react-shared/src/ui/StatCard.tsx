import { C, S, T, W } from "../utils/design";
import { Sparkline } from "../charts/Sparkline";
import { Card } from "./Card";

interface StatCardProps {
  emoji: string;
  label: string;
  value: string | number;
  change?: string;
  trend?: "up" | "down" | "neutral";
  accent?: string;
  sparkData?: number[];
  className?: string;
}

export function StatCard({ emoji, label, value, change, trend, accent, sparkData, className }: StatCardProps) {
  const trendColor = trend === "up" ? C.success : trend === "down" ? C.negative : C.text2;
  const arrow = trend === "up" ? "↑" : trend === "down" ? "↓" : "";
  const col = accent || C.primary;

  return (
    <Card accent={col} className={className} style={{ padding: `${S[5]}px ${S[5]}px ${S[4]}px`, flex: 1, minWidth: 140 }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start" }}>
        <div style={{ flex: 1 }}>
          <div style={{ fontSize: T.lg + 2, marginBottom: S[2] }}>{emoji}</div>
          <div style={{ fontSize: T["2xl"], fontWeight: W.bold, color: C.text, lineHeight: 1 }}>{value}</div>
          <div style={{ fontSize: T.xs, color: C.text2, marginTop: S[1] + 1, textTransform: "uppercase", letterSpacing: "0.06em" }}>
            {label}
          </div>
          {change && (
            <div style={{ fontSize: T.xs, fontWeight: W.semi, color: trendColor, marginTop: S[1] + 2 }}>
              {arrow} {change}
            </div>
          )}
        </div>
        {sparkData && sparkData.length > 1 && (
          <div style={{ alignSelf: "flex-end", marginBottom: S[1], opacity: 0.9 }}>
            <Sparkline data={sparkData} color={col} width={80} height={32} />
          </div>
        )}
      </div>
    </Card>
  );
}
