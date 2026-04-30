import type { CSSProperties, ReactNode } from "react";

import { C, R, T, W } from "@/lib/design";

interface BadgeProps {
  children: ReactNode;
  color?: string;
  style?: CSSProperties;
}

export function Badge({ children, color = C.primary, style }: BadgeProps) {
  return (
    <span style={{
      display: "inline-flex",
      alignItems: "center",
      padding: "2px 8px",
      borderRadius: R.full,
      fontSize: T.xs,
      fontWeight: W.medium,
      background: `${color}20`,
      color,
      ...style,
    }}>
      {children}
    </span>
  );
}

export function StatusDot({ status }: { status?: string }) {
  const color = status === "mastered" ? "#16A34A" : status === "learning" ? "#D97706" : C.text3;
  return (
    <span style={{
      display: "inline-block",
      width: 8,
      height: 8,
      borderRadius: "50%",
      background: color,
      flexShrink: 0,
    }} />
  );
}
