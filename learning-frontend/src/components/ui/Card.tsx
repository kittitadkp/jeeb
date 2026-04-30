import type { CSSProperties, ReactNode } from "react";

import { C, R, S } from "@/lib/design";

interface CardProps {
  children: ReactNode;
  style?: CSSProperties;
  onClick?: () => void;
}

export function Card({ children, style, onClick }: CardProps) {
  return (
    <div
      onClick={onClick}
      style={{
        background: C.surface,
        border: `1px solid ${C.border}`,
        borderRadius: R.card,
        padding: S[5],
        boxShadow: C.shadow,
        cursor: onClick ? "pointer" : undefined,
        transition: onClick ? "border-color .15s" : undefined,
        ...style,
      }}
    >
      {children}
    </div>
  );
}
