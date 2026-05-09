import { type ReactNode } from "react";

import { C, R, T, W } from "../utils/design";

interface BadgeProps {
  children: ReactNode;
  color?: string;
  className?: string;
}

export function Badge({ children, color, className }: BadgeProps) {
  const col = color || C.primary;
  return (
    <span
      className={className}
      style={{
        background: `${col}20`,
        color: col,
        border: `1px solid ${col}40`,
        borderRadius: R.sm,
        padding: `2px ${T.xs - 3}px`,
        fontSize: T.xs,
        fontWeight: W.semi,
        letterSpacing: "0.02em",
        textTransform: "capitalize",
        display: "inline-block",
      }}
    >
      {children}
    </span>
  );
}
