# Frontend Runtime Notes

## Main frontend

### Data-backed pages

- Dashboard
- Workouts
- Study
- Sleep
- Finance
- Calendar

### Local-only pages

- Goals
- Events
- Settings

These pages currently manage state in the browser and do not call backend APIs.

### Auth and API

- Keycloak login mode: `login-required`
- PKCE method: `S256`
- API base URL: `import.meta.env.VITE_API_URL`
- Token injection happens centrally in `src/lib/api.ts`

## Learning frontend

- Uses the same Keycloak login pattern
- Default local API target is Kong: `http://localhost:30088/learning`
- Topic detail view contains browse, flashcard, recall, and progress modes

## Build and deployment caveat

The frontends are static bundles behind Nginx. Even though Vault renders env files into the frontend pods, the current images do not read them at runtime. To change frontend env values, rebuild and redeploy the image.
