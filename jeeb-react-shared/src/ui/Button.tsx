import { type CSSProperties, forwardRef, type ReactNode, useState } from "react";

import { C, R, T, W } from "../utils/design";
import { cn } from "../utils/utils";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "default" | "outline" | "ghost" | "danger" | "secondary";
  size?: "sm" | "default" | "lg";
  color?: string;
  children?: ReactNode;
  style?: CSSProperties;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "default", size = "default", color, children, style, disabled, onClick, ...props }, ref) => {
    const [hov, setHov] = useState(false);

    const col = color || C.primary;
    const pad = size === "sm" ? `${T.xs - 1}px 14px` : size === "lg" ? "12px 24px" : "9px 18px";
    const fs  = size === "sm" ? T.sm : size === "lg" ? T.md : T.base;

    let bg: string, border: string, textColor: string;

    if (variant === "outline") {
      bg = hov ? `${col}18` : "transparent";
      border = `1px solid ${hov ? col : col + "60"}`;
      textColor = col;
    } else if (variant === "danger") {
      bg = hov ? `${C.danger}25` : `${C.danger}15`;
      border = `1px solid ${C.danger}50`;
      textColor = C.danger;
    } else if (variant === "ghost") {
      bg = hov ? C.surface2 : "transparent";
      border = "1px solid transparent";
      textColor = C.text2;
    } else if (variant === "secondary") {
      bg = hov ? C.surface3 : C.surface2;
      border = `1px solid ${C.border2}`;
      textColor = C.text;
    } else {
      bg = hov ? `${col}dd` : col;
      border = "1px solid transparent";
      textColor = "#fff";
    }

    return (
      <button
        ref={ref}
        onClick={onClick}
        disabled={disabled}
        onMouseEnter={() => setHov(true)}
        onMouseLeave={() => setHov(false)}
        className={cn(className)}
        style={{
          background: bg,
          border,
          color: textColor,
          padding: pad,
          borderRadius: R.lg,
          fontSize: fs,
          fontWeight: W.semi,
          display: "inline-flex",
          alignItems: "center",
          gap: T.xs - 5,
          transition: "all 0.15s",
          opacity: disabled ? 0.5 : 1,
          cursor: disabled ? "not-allowed" : "pointer",
          letterSpacing: "0.01em",
          ...style,
        }}
        {...props}
      >
        {children}
      </button>
    );
  },
);
Button.displayName = "Button";
