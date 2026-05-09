# Plan: jeeb-react-shared — Shared React Component Library

## Context

Both `frontend` and `learning-frontend` duplicate 6 UI components, the entire auth/keycloak stack, design tokens, and utility helpers. Every time a primitive changes it must be updated in two places. This plan extracts those into a standalone npm package (`@jeeb/react-shared`) published to the existing Nexus registry, so both apps consume a single source of truth.

---

## Scope

### What goes into the shared library

| Category | Files to extract |
|----------|-----------------|
| **UI primitives** | Button, Card, Badge, SectionLabel, StatCard, States |
| **Charts** | bar-chart, donut-chart, sparkline, period-selector |
| **Auth** | AuthProvider, AuthContext, useAuth, keycloak.ts |
| **Utils/tokens** | design.ts, utils.ts (cn helper), theme store |

### What stays app-specific
- Layout (AppLayout, Sidebar, Header, BottomNav)
- Domain hooks (useWorkouts, useStudy, useTopics, etc.)
- Domain components (ExercisePicker, FlashcardTool, RecallTool)
- Pages

---

## New Repository: `jeeb-react-shared`

Create a new standalone repo (sibling to `jeeb/frontend`, `jeeb/learning-frontend`).

### Directory layout

```
jeeb-react-shared/
  src/
    ui/
      Button.tsx
      Card.tsx
      Badge.tsx
      SectionLabel.tsx
      StatCard.tsx
      States.tsx
      index.ts
    charts/
      BarChart.tsx
      DonutChart.tsx
      Sparkline.tsx
      PeriodSelector.tsx
      index.ts
    auth/
      AuthContext.ts
      AuthProvider.tsx
      useAuth.ts
      keycloak.ts
      index.ts
    utils/
      design.ts          # color/spacing tokens
      utils.ts           # cn() clsx+tailwind-merge helper
      theme.ts           # dark-mode store
      index.ts
    index.ts             # barrel — re-exports all four categories
  package.json
  tsconfig.json
  tsup.config.ts
  .npmrc
  Jenkinsfile
  .gitignore
```

---

## Step-by-Step Implementation

### Step 1 — Scaffold the repository
- [x] `git init jeeb-react-shared && cd jeeb-react-shared`
- [x] Create `package.json`:
  ```json
  {
    "name": "@jeeb/react-shared",
    "version": "1.0.0",
    "main": "dist/index.js",
    "module": "dist/index.mjs",
    "types": "dist/index.d.ts",
    "exports": {
      ".": { "import": "./dist/index.mjs", "require": "./dist/index.js" },
      "./ui": { "import": "./dist/ui/index.mjs", "require": "./dist/ui/index.js" },
      "./charts": { "import": "./dist/charts/index.mjs", "require": "./dist/charts/index.js" },
      "./auth": { "import": "./dist/auth/index.mjs", "require": "./dist/auth/index.js" },
      "./utils": { "import": "./dist/utils/index.mjs", "require": "./dist/utils/index.js" }
    },
    "peerDependencies": {
      "react": ">=19",
      "react-dom": ">=19",
      "tailwindcss": ">=3"
    },
    "dependencies": {
      "@radix-ui/react-dialog": "...",
      "@radix-ui/react-dropdown-menu": "...",
      "@radix-ui/react-select": "...",
      "class-variance-authority": "0.7.1",
      "clsx": "2.1.1",
      "keycloak-js": "26.2.4",
      "lucide-react": "1.9.0",
      "tailwind-merge": "3.5.0",
      "@tanstack/react-query": "5.100.1"
    },
    "devDependencies": {
      "tsup": "^8",
      "typescript": "~6.0.2",
      "@types/react": "^19"
    },
    "scripts": {
      "build": "tsup",
      "typecheck": "tsc --noEmit"
    }
  }
  ```

### Step 2 — Configure tsup (build tool)
- [x] Create `tsup.config.ts`:
  ```ts
  import { defineConfig } from 'tsup'
  export default defineConfig({
    entry: {
      index: 'src/index.ts',
      'ui/index': 'src/ui/index.ts',
      'charts/index': 'src/charts/index.ts',
      'auth/index': 'src/auth/index.ts',
      'utils/index': 'src/utils/index.ts',
    },
    format: ['esm', 'cjs'],
    dts: true,
    sourcemap: true,
    external: ['react', 'react-dom', 'tailwindcss'],
    treeshake: true,
  })
  ```

### Step 3 — Configure tsconfig
- [x] Create `tsconfig.json` with `"jsx": "react-jsx"`, `"moduleResolution": "bundler"`, `"declaration": true`

### Step 4 — Copy and clean source files
- [x] Copy from `frontend/src/components/ui/` → `src/ui/` (6 components)
- [x] Copy from `frontend/src/components/ui/` chart files → `src/charts/` (4 files)
- [x] Copy from `frontend/src/lib/auth.tsx`, `auth-context.ts`, `useAuth.ts`, `keycloak.ts` → `src/auth/`
- [x] Copy from `frontend/src/lib/design.ts`, `utils.ts` + `frontend/src/store/theme.ts` → `src/utils/`
- [x] Write barrel `index.ts` files for each category
- [x] Write root `src/index.ts` re-exporting all categories
- [x] Remove any app-specific imports (API client, router, page references)

### Step 5 — Configure Nexus publishing
- [ ] Create `.npmrc`:
  ```
  registry=http://localhost:30050/repository/npm-hosted/
  //localhost:30050/repository/npm-hosted/:_authToken=${NEXUS_NPM_TOKEN}
  ```
- [ ] Add `.npmrc` to `.gitignore` (token is injected by Jenkins)

### Step 6 — Jenkinsfile for CI/CD
- [x] Create `Jenkinsfile`:
  ```groovy
  pipeline {
    agent any
    environment {
      NEXUS_NPM_TOKEN = credentials('nexus-npm-token')
    }
    stages {
      stage('Install')   { steps { sh 'npm ci' } }
      stage('Typecheck') { steps { sh 'npm run typecheck' } }
      stage('Build')     { steps { sh 'npm run build' } }
      stage('Publish')   { steps { sh 'npm publish --registry http://nexus.jeeb-infra.svc.cluster.local:8081/repository/npm-hosted/' } }
    }
  }
  ```

### Step 7 — Update consuming frontends

**In both `frontend/` and `learning-frontend/`:**
- [x] Add to `.npmrc`: point to Nexus npm registry
- [x] `npm install @jeeb/react-shared`
- [x] Update `tailwind.config.js` content array to include:
  ```js
  content: [
    './src/**/*.{ts,tsx}',
    './node_modules/@jeeb/react-shared/dist/**/*.{js,mjs}',
  ]
  ```
  > Note: Because we ship source TSX via `exports`, alternatively add `node_modules/@jeeb/react-shared/src/**/*.tsx` and configure Tailwind's `content` to scan it. This avoids purging shared class names.
- [ ] Replace local imports with package imports, e.g.:
  ```ts
  // Before
  import { Button } from '../components/ui/button'
  import { useAuth } from '../lib/useAuth'
  // After
  import { Button } from '@jeeb/react-shared/ui'
  import { useAuth } from '@jeeb/react-shared/auth'
  ```
- [ ] Delete the now-redundant local copies of the extracted files

### Step 8 — Kubernetes / CI for frontend images
- [x] No k8s changes needed — frontends compile the shared lib at build time (it's a build dependency, not a runtime service)
- [ ] Jenkins frontend pipelines already run `npm run build`; they just need Nexus `.npmrc` credentials available during the build stage

---

## Verification

1. **Build the library**: `cd jeeb-react-shared && npm run build` → `dist/` produced with `.mjs`, `.js`, `.d.ts`
2. **Local smoke test**: In `learning-frontend`, temporarily `npm install ../jeeb-react-shared` (file path) and run `npm run dev` — verify shared components render
3. **Publish to Nexus**: Trigger Jenkins job → check Nexus UI shows `@jeeb/react-shared@1.0.0`
4. **Consume from Nexus**: `npm install @jeeb/react-shared` in both frontends → `npm run build` succeeds with no TS errors
5. **Tailwind check**: Verify shared component classes (e.g., `bg-blue-600`, `rounded-lg`) survive Tailwind's purge by inspecting the built CSS

---

## Files to create (new repo)
- `package.json`, `tsconfig.json`, `tsup.config.ts`, `.npmrc`, `Jenkinsfile`, `.gitignore`
- `src/ui/` (6 files + index), `src/charts/` (4 files + index), `src/auth/` (4 files + index), `src/utils/` (3 files + index), `src/index.ts`

## Files to modify (jeeb repo)
- `frontend/package.json` — add `@jeeb/react-shared`
- `frontend/.npmrc` — Nexus registry
- `frontend/tailwind.config.js` — content paths
- `frontend/src/**` — replace local imports, delete extracted files
- `learning-frontend/package.json` — same
- `learning-frontend/.npmrc` — same
- `learning-frontend/tailwind.config.js` — same
- `learning-frontend/src/**` — same
