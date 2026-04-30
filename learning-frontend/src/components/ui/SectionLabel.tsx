import type { ReactNode } from "react";

import { C, T, W } from "@/lib/design";

export function SectionLabel({ children }: { children: ReactNode }) {
  return (
    <div style={{ fontSize: T.xs, fontWeight: W.semi, color: C.text2, textTransform: "uppercase", letterSpacing: "0.08em" }}>
      {children}
    </div>
  );
}
