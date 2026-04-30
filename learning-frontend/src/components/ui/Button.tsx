import type { ButtonHTMLAttributes, CSSProperties } from "react";

import { C, R, T, W } from "@/lib/design";

type Variant = "primary" | "secondary" | "ghost" | "danger";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  size?: "sm" | "md" | "lg";
  fullWidth?: boolean;
}

const variantStyles: Record<Variant, CSSProperties> = {
  primary:   { background: C.primary, color: "#fff", border: "none" },
  secondary: { background: C.surface2, color: C.text, border: `1px solid ${C.border}` },
  ghost:     { background: "transparent", color: C.text2, border: "none" },
  danger:    { background: C.danger, color: "#fff", border: "none" },
};

const sizeStyles = {
  sm: { padding: "4px 12px", fontSize: T.sm, borderRadius: R.md },
  md: { padding: "7px 16px", fontSize: T.base, borderRadius: R.md },
  lg: { padding: "10px 20px", fontSize: T.md, borderRadius: R.lg },
};

export function Button({ variant = "primary", size = "md", fullWidth, style, children, ...props }: ButtonProps) {
  return (
    <button
      style={{
        display: "inline-flex",
        alignItems: "center",
        justifyContent: "center",
        gap: 6,
        fontWeight: W.medium,
        cursor: "pointer",
        transition: "opacity .15s",
        width: fullWidth ? "100%" : undefined,
        ...variantStyles[variant],
        ...sizeStyles[size],
        ...style,
      }}
      {...props}
    >
      {children}
    </button>
  );
}
