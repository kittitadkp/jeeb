# Plan: Delete jeeb-react-shared entirely + update frontend

## Context

The previous session stripped `jeeb-react-shared` down to only `utils` (design tokens, `cn`, `useDarkMode`). This session finishes the job: move those three tiny utilities directly into the frontend and remove the shared library package and directory altogether. The goal is zero external shared-lib dependency.

---

## Step 1 — Inline utils source into frontend wrapper files

### `frontend/src/lib/design.ts`
Replace the two re-export lines with the full source from `jeeb-react-shared/src/utils/design.ts` (no import adjustments needed — it has no relative imports):

```ts
export const C = { bg: "var(--c-bg)", surface: "var(--c-surface)", ... } as const;
export const T = { xs: 11, sm: 12, ... } as const;
export const W = { normal: 400, medium: 500, semi: 600, bold: 700 } as const;
export const R = { sm: 6, md: 8, lg: 10, card: 14, full: 9999 } as const;
export const S = { 1: 4, 2: 8, 3: 12, 4: 16, 5: 20, 6: 24, 8: 32, 10: 40 } as const;
export const SECTION_COLORS = { dashboard: "#7c6ef5", workouts: "#f43f5e", ... } as const;
export type SectionKey = keyof typeof SECTION_COLORS;
export type ThemeMode = 'light' | 'dark' | 'system';
export type PrimaryColor = 'blue' | 'purple' | 'green' | 'rose' | 'amber';
export interface UserPreferences { theme: ThemeMode; primary_color: PrimaryColor; currency: string; week_start: 'monday' | 'sunday'; }
export const DEFAULT_PREFERENCES: UserPreferences = { theme: 'system', primary_color: 'blue', currency: 'THB', week_start: 'monday' };
export const PRIMARY_COLORS: { value: PrimaryColor; label: string; hex: string }[] = [...];
```

### `frontend/src/lib/utils.ts`
Replace `export { cn } from "@jeeb/react-shared/utils"` with the inline implementation; keep the existing local date helpers:

```ts
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function getTodayStr(): string { ... }
export function getYesterdayStr(): string { ... }
```

(`clsx` and `tailwind-merge` are already in `frontend/package.json`.)

### `frontend/src/store/theme.ts`
Replace `export { useDarkMode } from "@jeeb/react-shared/utils"` with the inline implementation from `jeeb-react-shared/src/utils/theme.ts`:

```ts
import { useEffect, useState } from "react";

export function useDarkMode() {
  const [dark, setDark] = useState(() => {
    const stored = localStorage.getItem("theme");
    return stored !== "light";
  });
  useEffect(() => {
    document.documentElement.classList.toggle("dark", dark);
    localStorage.setItem("theme", dark ? "dark" : "light");
  }, [dark]);
  return { dark, toggleDark: () => setDark((d) => !d) };
}
```

---

## Step 2 — Fix direct `@jeeb/react-shared` root imports in 5 files

All of these import utils exports directly from the shared lib root. Redirect each to the local wrapper at `@/lib/design`:

| File | Current import | Fix |
|------|---------------|-----|
| `src/hooks/useAppSettings.ts:2` | `import type { UserPreferences } from "@jeeb/react-shared"` | → `@/lib/design` |
| `src/hooks/useUser.ts:2` | `import type { UserPreferences } from "@jeeb/react-shared"` | → `@/lib/design` |
| `src/providers/ProfileProvider.tsx:2-3` | `UserPreferences, DEFAULT_PREFERENCES, PRIMARY_COLORS from "@jeeb/react-shared"` | → `@/lib/design` |
| `src/pages/Settings.tsx:5-6` | `PrimaryColor, ThemeMode, UserPreferences, PRIMARY_COLORS from "@jeeb/react-shared"` | → `@/lib/design` |
| `e2e/mocks/mock-api.ts:2` | `DEFAULT_PREFERENCES from "@jeeb/react-shared"` | → `../../src/lib/design` (relative, since `@/` alias may not work in e2e) |

---

## Step 3 — Fix dynamic imports in `src/types/index.ts`

Two lines use `import("@jeeb/react-shared").UserPreferences` as an inline type. Replace both with a static import at the top of the file and use the type directly:

```ts
// Add at top:
import type { UserPreferences } from "@/lib/design";

// Replace both occurrences of:
//   import("@jeeb/react-shared").UserPreferences
// with:
//   UserPreferences
```

---

## Step 4 — Remove dependency & delete the package

**`frontend/package.json`** — remove the `@jeeb/react-shared` line from dependencies.

```powershell
cd frontend && npm install --legacy-peer-deps
```

**Delete the directory:**
```powershell
Remove-Item -Recurse -Force "D:/personal/jeeb/jeeb-react-shared"
```

---

## Verification

```powershell
# 1. No remaining references to @jeeb/react-shared in frontend
grep -r "@jeeb/react-shared" D:/personal/jeeb/frontend --include="*.ts" --include="*.tsx"

# 2. Frontend builds clean
cd frontend && npm run build

# 3. Directory is gone
Test-Path "D:/personal/jeeb/jeeb-react-shared"  # should be False
```
