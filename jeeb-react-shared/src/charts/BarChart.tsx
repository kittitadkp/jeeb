import { C, R, T } from "../utils/design";

interface BarChartData {
  label: string;
  value: number;
}

interface BarChartProps {
  data: BarChartData[];
  color?: string;
  height?: number;
}

export function BarChart({ data, color = C.primary, height = 80 }: BarChartProps) {
  if (!data || data.length === 0) return null;
  const vals = data.map((d) => d.value);
  const max = Math.max(...vals) || 1;
  const n = data.length;
  const BAR_W = 28, GAP = 10;
  const VW = n * (BAR_W + GAP) - GAP;
  const VH = height + 20;
  return (
    <div style={{ width: "100%", overflowX: "hidden" }}>
      <svg
        width="100%"
        height={VH}
        viewBox={`0 0 ${VW} ${VH}`}
        preserveAspectRatio="none"
        style={{ display: "block" }}
      >
        {data.map((d, i) => {
          const barH = Math.max(3, (d.value / max) * height);
          const x = i * (BAR_W + GAP);
          const y = height - barH;
          const isLast = i === n - 1;
          const fs = Math.min(T.xs - 1, (VW / n) * 0.38);
          return (
            <g key={i}>
              <rect
                x={x} y={y} width={BAR_W} height={barH} rx={R.sm - 1}
                fill={isLast ? color : `${color}45`}
                style={{ transition: "all 0.3s" }}
              />
              <text
                x={x + BAR_W / 2} y={height + 15}
                textAnchor="middle" fontSize={fs}
                fill={C.text2} fontFamily="Inter,sans-serif"
              >
                {d.label}
              </text>
            </g>
          );
        })}
      </svg>
    </div>
  );
}
