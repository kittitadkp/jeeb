import { C, T, W } from "../utils/design";

interface SectionLabelProps {
  color: string;
  children: string;
}

export function SectionLabel({ color, children }: SectionLabelProps) {
  return (
    <div style={{ display: "flex", alignItems: "center", gap: 8, marginBottom: 0 }}>
      <div style={{ width: 3, height: 22, borderRadius: 2, background: color, flexShrink: 0 }} />
      <h1 style={{ fontSize: T.xl, fontWeight: W.bold, color: C.text, margin: 0, fontFamily: '"Space Grotesk", sans-serif' }}>
        {children}
      </h1>
    </div>
  );
}
