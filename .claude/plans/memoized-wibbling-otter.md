# Plan: Design hooks and rules for the jeeb project

## Context

Hooks and rules were added ad-hoc in this session. The goal is a coherent, purposeful system:
- **Hooks** (`settings.json`) — automate quality enforcement silently, without friction
- **Rules** (`CLAUDE.md`) — give Claude durable conventions so it doesn't guess
- **Obsidian conventions** (`wiki/CONVENTIONS.md`) — already created ✅

Existing hooks (added this session, need refinement):
- `PostToolUse(Write|Edit)` → gofmt on `.go` ✅
- `PostToolUse(Bash)` → wiki reminder on `git commit` ✅

---

## Design principles

1. **Hooks should be silent and fast** — run automatically, produce no output unless something is wrong
2. **Rules should explain WHY** — a rule without a reason gets ignored when it seems inconvenient
3. **No PreToolUse safety hooks** — personal dev project, friction has no payoff here
4. **No slow hooks** — full `go test` or `tsc --noEmit` on every edit is too slow; save for commit

---

## Step 1 — Refine hooks in `.claude/settings.json`

### 1a. Extend the Go hook: gofmt + go vet together

Replace the current gofmt-only command with a combined command that also runs `go vet ./...` in the package directory. `go vet` catches real bugs (wrong printf args, unreachable code, etc.) that gofmt doesn't.

```python
# PostToolUse Write|Edit — new command
python3 -c "
import json, os, subprocess, pathlib
d = json.loads(os.environ.get('CLAUDE_TOOL_INPUT', '{}'))
f = d.get('file_path', '')
if not f.endswith('.go'):
    exit(0)
subprocess.run(['gofmt', '-w', f])
pkg = str(pathlib.Path(f).parent)
result = subprocess.run(['go', 'vet', './...'], cwd=pkg, capture_output=True, text=True)
if result.returncode != 0:
    print(result.stdout + result.stderr)
"
```

Output only printed when `go vet` finds issues → Claude sees it as context and can fix proactively.

### 1b. Add TypeScript check hook

After writing `.ts` or `.tsx` files, run `tsc --noEmit` scoped to the relevant project. This is fast when targeted (not whole repo).

```python
# PostToolUse Write|Edit — new entry for TS
python3 -c "
import json, os, subprocess, pathlib
d = json.loads(os.environ.get('CLAUDE_TOOL_INPUT', '{}'))
f = d.get('file_path', '')
if not f.endswith(('.ts', '.tsx')):
    exit(0)
# Find nearest tsconfig.json
p = pathlib.Path(f).parent
while p != p.parent:
    if (p / 'tsconfig.json').exists():
        break
    p = p.parent
result = subprocess.run(['npx', 'tsc', '--noEmit'], cwd=str(p), capture_output=True, text=True)
if result.returncode != 0:
    print(result.stdout[-3000:])  # cap output
"
```

### 1c. Remove wiki reminder hook (replace with CLAUDE.md rule)

The `git commit` bash hook is noisy — it fires on every commit. Move this to a CLAUDE.md rule instead: Claude reads the rule at session start and knows to suggest `/wiki-update` when relevant. Hooks are better for automation than reminders.

---

## Step 2 — Add rules to `CLAUDE.md`

Add three sections after the existing `## Wiki` section:

### `## Go conventions`

```markdown
## Go conventions

**New feature = all 5 hexagonal layers.** Adding a domain concept requires:
`domain/` → `port/in/` → `port/out/` → `usecase/` → `adapter/in/http/handler/` + `adapter/out/mongo/`
Never skip a layer by importing across boundaries.

**Error handling:** Wrap with context using `fmt.Errorf("action: %w", err)`. Never `panic` in domain or usecase. Only `log.Fatal` in `cmd/` entry points.

**Tests:** Usecase logic that contains branching must have table-driven tests. Repositories and handlers do not need unit tests (integration tests cover them via the running cluster).

**Avoid:** global state, `init()` functions, embedding structs for behavior reuse (use interfaces).
```

### `## TypeScript / React conventions`

```markdown
## TypeScript / React conventions

**API calls:** Always go through hooks in `src/hooks/`. Components never call `api.*` directly.

**Types:** All interfaces in `src/types/index.ts` must use `snake_case` field names matching the backend JSON exactly. Never add a type that doesn't have a corresponding backend struct.

**No `any`:** Use `unknown` and narrow, or define a proper interface. `as any` is a last resort for third-party interop only.

**Mutations:** Body type is always `Omit<Entity, 'id' | 'user_id' | 'created_at' | 'updated_at'>`. The server assigns those fields.

**Design tokens:** Use `C`, `T`, `W`, `R`, `S` from `src/lib/design.ts` for inline styles. Don't scatter raw hex values or pixel numbers in components.
```

### `## Git commit format`

```markdown
## Git commit format

```
type: short description (imperative, lowercase, no period)
```

Types: `feat` `fix` `chore` `refactor` `docs` `test` `infra`

Examples:
- `feat: add sleep stats endpoint`
- `fix: workout duration not saved on update`
- `chore: bump keycloak-js to v26`
- `infra: add vault policy for learning-backend`
```

### Extend `## Wiki` — add when NOT to run

Add one sentence: "Do **not** run `/wiki-update` for routine bug fixes, UI tweaks, or dependency bumps." (already added — verify it's there).

---

## Step 3 — Update `wiki/CONVENTIONS.md`

Already created. No changes needed unless the lifecycle section needs expanding.

---

## Files to modify

| File | Change |
|------|--------|
| `.claude/settings.json` | Replace gofmt hook with gofmt+vet; add TS hook; remove commit reminder hook |
| `CLAUDE.md` | Add Go conventions, TypeScript/React conventions, Git commit format sections |
| `wiki/CONVENTIONS.md` | Already done ✅ |

---

## Verification

After implementation:
1. Edit a `.go` file — confirm `gofmt` reformats it silently, confirm `go vet` output appears only if there's an issue
2. Introduce a deliberate `go vet` violation (e.g. `fmt.Printf("%d", "string")`) — confirm hook prints the error
3. Edit a `.ts` file with a type error — confirm tsc output appears
4. Edit a `.ts` file with no errors — confirm no output
5. Read `CLAUDE.md` and confirm all three new sections are present and accurate
6. Run `git commit` — confirm no noisy hook output (reminder is gone)
