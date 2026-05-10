# Plan: User Profile Feature

## Context
The Settings page currently stores all preferences (theme, accent color, currency, week start) in local component state only — they reset on refresh and don't sync across devices. The `/me` endpoint returns a user from MongoDB but has no profile fields. This plan wires the full stack: backend stores preferences in MongoDB, a new `PUT /me` endpoint persists them, and the frontend Settings page reads/writes through TanStack Query. The shared library is updated to expose the `UserPreferences` type so both frontends benefit.

---

## Scope

Settings persisted to DB: **theme, accent color, currency, week start, display name (editable)**

---

## Step-by-step Plan

### Step 1 — Backend: Extend domain model
**File:** `backend/internal/domain/user.go`

Add `UserPreferences` struct and embed it in `User`:
```go
type UserPreferences struct {
    Theme       string `bson:"theme"        json:"theme"`        // "light"|"dark"|"system"
    AccentColor string `bson:"accent_color" json:"accent_color"` // "blue"|"purple"|"green"|"rose"|"amber"
    Currency    string `bson:"currency"     json:"currency"`     // "THB"|"USD"|"EUR"…
    WeekStart   string `bson:"week_start"   json:"week_start"`   // "monday"|"sunday"
}

// Add to User struct:
Preferences UserPreferences `bson:"preferences" json:"preferences"`
```

---

### Step 2 — Backend: Add UpdateProfile to port/in
**File:** `backend/internal/port/in/user.go`

Add request DTO and method to interface:
```go
type UpdateProfileRequest struct {
    DisplayName string                `json:"display_name" validate:"required,min=1,max=100"`
    Preferences domain.UserPreferences `json:"preferences"  validate:"required"`
}

// Add to UserUseCase interface:
UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*domain.User, error)
```

---

### Step 3 — Backend: Add Update to port/out
**File:** `backend/internal/port/out/repositories.go`

Add `Update` to `UserRepository` interface:
```go
Update(ctx context.Context, user *domain.User) (*domain.User, error)
```

---

### Step 4 — Backend: Implement Update in mongo repository
**File:** `backend/internal/adapter/out/mongo/user_repository.go`

Add `Update` method using `$set` on `display_name`, `preferences`, `updated_at`. Pattern matches existing `Upsert`.

---

### Step 5 — Backend: Implement UpdateProfile in usecase
**File:** `backend/internal/usecase/user_usecase.go`

- Fetch user by ID, validate ownership
- Apply `DisplayName` and `Preferences` from request
- Set `UpdatedAt = time.Now()`
- Call `repo.Update()`
- Return updated `*domain.User`

Default preferences on new users (applied in `GetOrCreate`):
```go
Preferences: domain.UserPreferences{
    Theme:       "system",
    AccentColor: "blue",
    Currency:    "THB",
    WeekStart:   "monday",
}
```

---

### Step 6 — Backend: Add UpdateProfile handler
**File:** `backend/internal/adapter/in/http/handler/user_handler.go`

Add `UpdateProfile(w http.ResponseWriter, r *http.Request)`:
- Decode JSON body into `port.UpdateProfileRequest`
- Validate with `validator/v10`
- Extract `userID` from context via `userIDFromCtx(r)` (same pattern as all other handlers)
- Call `h.useCase.UpdateProfile()`
- Return `200 OK` with updated user JSON

---

### Step 7 — Backend: Register PUT /me route
**File:** `backend/internal/adapter/in/http/router.go`

```go
r.Put("/me", h.User.UpdateProfile)
```

---

### Step 8 — Shared lib: Add UserPreferences type + accent color util
**File:** `jeeb-react-shared/src/utils/design.ts` (already modified per git status)

Add `UserPreferences` TypeScript interface:
```ts
export interface UserPreferences {
    theme: 'light' | 'dark' | 'system'
    accent_color: 'blue' | 'purple' | 'green' | 'rose' | 'amber'
    currency: string
    week_start: 'monday' | 'sunday'
}

export const DEFAULT_PREFERENCES: UserPreferences = {
    theme: 'system',
    accent_color: 'blue',
    currency: 'THB',
    week_start: 'monday',
}
```

Add `getAccentColorClass(accent: string): string` helper that maps accent name to Tailwind class.

**File:** `jeeb-react-shared/src/index.ts` (or appropriate barrel)  
Re-export `UserPreferences` and `DEFAULT_PREFERENCES`.

**After changes:** bump `package.json` version, `npm run build`, publish to Nexus.

---

### Step 9 — Frontend: Update User type
**File:** `frontend/src/types/index.ts`

```ts
import { UserPreferences } from '@jeeb/react-shared'

export interface User {
    id: string
    keycloak_id: string
    email: string
    display_name: string
    preferences: UserPreferences
    created_at: string
    updated_at: string
}
```

---

### Step 10 — Frontend: Add useUpdateProfile mutation
**File:** `frontend/src/hooks/useUser.ts`

```ts
export function useUpdateProfile() {
    return useMutation({
        mutationFn: (data: { display_name: string; preferences: UserPreferences }) =>
            api.put<User>('/me', data),
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['user', 'me'] }),
    })
}
```

---

### Step 11 — Frontend: Add ProfileProvider to sync preferences on load
**File:** `frontend/src/providers/ProfileProvider.tsx` (new file)

On mount, reads `useMe()` data and applies:
- Calls `setTheme(profile.preferences.theme)` into the existing `useDarkMode` hook
- Applies accent color as a CSS variable or Tailwind class on `<body>`

Wrap in `App.tsx` inside `<QueryClientProvider>`.

---

### Step 12 — Frontend: Wire Settings page
**File:** `frontend/src/pages/Settings.tsx`

- Replace all local `useState` for theme, accent color, currency, week start with values from `useMe()` data
- Replace display name read-only text with an editable `<input>`
- On each setting change, call `useUpdateProfile()` mutation (debounced or on blur/save button)
- Show loading/saving state feedback using existing `StatCard`/`Badge` patterns from `@jeeb/react-shared`

---

## Files Modified

| File | Change |
|------|--------|
| `backend/internal/domain/user.go` | Add `UserPreferences` struct + field on `User` |
| `backend/internal/port/in/user.go` | Add `UpdateProfileRequest` DTO + `UpdateProfile` to interface |
| `backend/internal/port/out/repositories.go` | Add `Update` to `UserRepository` interface |
| `backend/internal/adapter/out/mongo/user_repository.go` | Implement `Update` |
| `backend/internal/usecase/user_usecase.go` | Implement `UpdateProfile`, set defaults in `GetOrCreate` |
| `backend/internal/adapter/in/http/handler/user_handler.go` | Add `UpdateProfile` handler |
| `backend/internal/adapter/in/http/router.go` | Register `PUT /me` |
| `jeeb-react-shared/src/utils/design.ts` | Add `UserPreferences`, `DEFAULT_PREFERENCES`, `getAccentColorClass` |
| `jeeb-react-shared/package.json` | Bump version |
| `frontend/src/types/index.ts` | Add `preferences` field to `User` |
| `frontend/src/hooks/useUser.ts` | Add `useUpdateProfile` mutation |
| `frontend/src/providers/ProfileProvider.tsx` | New — sync profile to theme on load |
| `frontend/src/App.tsx` | Wrap with `ProfileProvider` |
| `frontend/src/pages/Settings.tsx` | Wire all settings to backend via `useUpdateProfile` |

---

## Verification

1. **Backend unit test:** `go test ./...` in `backend/`
2. **Backend manual:** `PUT /me` with `{"display_name":"Test","preferences":{"theme":"dark","accent_color":"purple","currency":"USD","week_start":"sunday"}}` → returns updated user; `GET /me` returns same values
3. **Frontend:** Open Settings page → change theme → refresh → theme persists (no flicker). Change currency → finance stats page reflects new currency symbol.
4. **Cross-device:** Log in from second browser → same theme/currency/preferences apply immediately.
