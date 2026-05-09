# @jeeb/react-shared

Shared React component library for the Jeeb platform. Used by `frontend` and `learning-frontend`.

## Packages

| Import path | Contents |
|-------------|----------|
| `@jeeb/react-shared/ui` | Button, Card, CardHeader, CardContent, Badge, SectionLabel, StatCard, EmptyState, LoadingState, ErrorState |
| `@jeeb/react-shared/charts` | BarChart, DonutChart, Sparkline, PeriodSelector |
| `@jeeb/react-shared/auth` | AuthProvider, AuthContext, useAuth, createKeycloak |
| `@jeeb/react-shared/utils` | C, T, W, R, S, SECTION_COLORS, cn, useDarkMode |

## Installation

Add to your app's `.npmrc`:

```ini
@jeeb:registry=http://localhost:30083/repository/npm-hosted/
//localhost:30083/repository/npm-hosted/:_auth=<base64(admin:password)>
```

Then install:

```bash
npm install @jeeb/react-shared
```

Add the dist path to your `tailwind.config.js` content array so Tailwind doesn't purge shared classes:

```js
content: [
  './src/**/*.{ts,tsx}',
  './node_modules/@jeeb/react-shared/dist/**/*.{js,mjs}',
],
```

## Usage

```tsx
import { Button, Card, StatCard, EmptyState } from '@jeeb/react-shared/ui'
import { BarChart, Sparkline, PeriodSelector } from '@jeeb/react-shared/charts'
import { AuthProvider, useAuth, createKeycloak } from '@jeeb/react-shared/auth'
import { C, T, cn, useDarkMode } from '@jeeb/react-shared/utils'
```

### Auth setup

Each app creates its own Keycloak instance (different `clientId` / default URLs) and passes it to `AuthProvider`:

```tsx
// lib/keycloak.ts
import { createKeycloak } from '@jeeb/react-shared/auth'
import { appConfig } from './app-config'

export const { keycloak, initKeycloak } = createKeycloak({
  url: appConfig.keycloakUrl,
  realm: appConfig.keycloakRealm,
  clientId: appConfig.keycloakClientId,
})

// main.tsx
import { AuthProvider } from '@jeeb/react-shared/auth'
import { setAccessToken } from './lib/api'
import keycloak, { initKeycloak } from './lib/keycloak'

<AuthProvider keycloak={keycloak} initKeycloak={initKeycloak} onTokenChange={setAccessToken}>
  <App />
</AuthProvider>
```

## Development

```bash
npm run build      # compile ESM + CJS + .d.ts into dist/
npm run typecheck  # tsc --noEmit
npm run dev        # watch mode
```

## Publishing

```bash
npm run build
npm publish --registry http://localhost:30083/repository/npm-hosted/
```

Bump `version` in `package.json` before publishing a new release. Jenkins publishes automatically on push to `main` via the `Jenkinsfile`.

## Design tokens

The `C`, `T`, `W`, `R`, `S` objects from `@jeeb/react-shared/utils` map to CSS variables defined in each app's global CSS. Components use `var(--c-bg)` etc. for structural colors (dark/light mode switches automatically) and hardcoded hex only where opacity-suffix arithmetic is needed (e.g. `${C.primary}20`).

Each consuming app is responsible for defining the CSS variables and its own primary color override.
