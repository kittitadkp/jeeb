import { C, T, W } from "../utils/design";

interface DonutSegment {
  label: string;
  value: number;
  color: string;
}

interface DonutChartProps {
  segments: DonutSegment[];
  size?: number;
  centerLabel?: string;
  centerSubLabel?: string;
}

export function DonutChart({ segments, size = 120, centerLabel, centerSubLabel }: DonutChartProps) {
  const total = segments.reduce((s, x) => s + x.value, 0) || 1;
  const r = 44, cx = size / 2, cy = size / 2;
  const circ = 2 * Math.PI * r;
  let offset = 0;
  const slices = segments.map((seg) => {
    const pct = seg.value / total;
    const dash = pct * circ;
    const slice = { ...seg, dash, offset };
    offset += dash;
    return slice;
  });
  const totalK = (segments.reduce((s, x) => s + x.value, 0) / 1000).toFixed(1);
  const label = centerLabel ?? `$${totalK}k`;
  const subLabel = centerSubLabel ?? "total";
  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
      <circle cx={cx} cy={cy} r={r} fill="none" stroke={C.surface3} strokeWidth="18" />
      {slices.map((s, i) => (
        <circle
          key={i} cx={cx} cy={cy} r={r} fill="none"
          stroke={s.color} strokeWidth="18"
          strokeDasharray={`${s.dash} ${circ - s.dash}`}
          strokeDashoffset={-s.offset}
          style={{ transform: "rotate(-90deg)", transformOrigin: `${cx}px ${cy}px` }}
        />
      ))}
      <text x={cx} y={cy - 4} textAnchor="middle" fontSize={T.base} fontWeight={W.bold} fill={C.text}>
        {label}
      </text>
      <text x={cx} y={cy + 14} textAnchor="middle" fontSize={T.xs - 2} fill={C.text2} fontFamily="Inter,sans-serif">
        {subLabel}
      </text>
    </svg>
  );
}
