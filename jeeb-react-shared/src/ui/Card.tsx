import { type CSSProperties, type ReactNode, useState } from "react";

import { C, R, S, T, W } from "../utils/design";
import { cn } from "../utils/utils";

interface CardProps {
  children?: ReactNode;
  style?: CSSProperties;
  accent?: string;
  onClick?: () => void;
  className?: string;
}

export function Card({ children, style, accent, onClick, className }: CardProps) {
  const [hov, setHov] = useState(false);
  const interactive = !!onClick;
  return (
    <div
      onClick={onClick}
      onMouseEnter={interactive ? () => setHov(true) : undefined}
      onMouseLeave={interactive ? () => setHov(false) : undefined}
      className={cn(className)}
      style={{
        background: C.surface,
        borderRadius: R.card,
        overflow: "hidden",
        borderTop: accent ? `2px solid ${accent}` : `1px solid ${C.border}`,
        borderRight: `1px solid ${interactive && hov ? C.border2 : C.border}`,
        borderBottom: `1px solid ${interactive && hov ? C.border2 : C.border}`,
        borderLeft: `1px solid ${interactive && hov ? C.border2 : C.border}`,
        transition: interactive ? "border-color 0.15s, box-shadow 0.15s" : undefined,
        boxShadow: interactive && hov ? C.shadowMd : C.shadow,
        cursor: interactive ? "pointer" : "default",
        position: "relative",
        ...style,
      }}
    >
      {children}
    </div>
  );
}

interface SectionProps {
  children?: ReactNode;
  className?: string;
  style?: CSSProperties;
}

export function CardHeader({ children, className, style }: SectionProps) {
  return (
    <div
      className={cn(className)}
      style={{ padding: `${S[3]}px ${S[5]}px`, borderBottom: `1px solid ${C.border}`, ...style }}
    >
      {children}
    </div>
  );
}

export function CardContent({ children, className, style }: SectionProps) {
  return (
    <div className={cn(className)} style={{ padding: `${S[4]}px ${S[5]}px`, ...style }}>
      {children}
    </div>
  );
}

export { T, W, R, S };
