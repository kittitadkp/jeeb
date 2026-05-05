# Prompt Efficiency Guide

How to get the most out of Claude Code in this project while keeping token usage low.

---

## Core Principles

### 1. Be specific about scope
Claude reads files to understand context. The more precisely you describe what to change, the fewer files it needs to read.

**Expensive:** "fix the bug in the backend"
**Efficient:** "fix the 404 on `GET /workouts/:id` — the handler is in `backend/internal/adapter/in/http/workout_handler.go`"

### 2. Name the layer
This project uses Clean/Hexagonal architecture. Tell Claude which layer you're working in so it doesn't scan the whole backend.

| You want to... | Say... |
|---|---|
| Add an endpoint | "add a handler in `adapter/in/http/`" |
| Add business logic | "add use case in `internal/usecase/`" |
| Add a DB query | "add repository method in `adapter/out/`" |
| Add a domain rule | "update domain struct in `internal/domain/`" |

### 3. Use slash commands
Each command loads only the context relevant to that domain — far cheaper than a blank prompt.

```
/backend  → Go API work
/frontend → React component work
/k8s      → Kubernetes manifests
/jenkins  → CI/CD pipelines
/docs     → Write or update documentation
/status   → Check pod health (read-only)
/logs <svc>   → Tail logs
/deploy <svc> → Restart and watch a deployment
```

---

## Context Management

### When to `/clear`
Start fresh when switching to an unrelated task. Stale context from a previous task costs tokens on every message.

- Finished a backend feature → switching to k8s manifest work → `/clear`
- Fixed a bug → now writing docs → `/clear`
- Don't clear mid-task — you'll lose useful context

### When to `/compact`
Use before starting a long task in the same area. It summarizes the conversation so far instead of carrying every message verbatim.

```
/compact
```

Run it when: conversation is 30+ messages, or Claude starts repeating itself / forgetting earlier changes.

### Auto-compact
`autoCompactEnabled: true` is set in `~/.claude/settings.json` — Claude will automatically compact when the context window fills. You don't need to do this manually unless you want to compact early.

---

## RTK — Token Filtering

RTK (`rtk hook claude`) is wired as a `PreToolUse` hook. It intercepts every Bash command Claude runs and filters verbose output before it reaches the model.

```bash
# See how many tokens RTK has saved this session
rtk gain

# See full history with per-command savings
rtk gain --history

# Find commands you ran manually that RTK could have optimized
rtk discover
```

RTK works transparently — you don't need to prefix commands. `git log` becomes `rtk git log` automatically.

---

## File Targeting

### Give exact paths when you know them
```
# Expensive — Claude will glob and read multiple files
"update the workout stats endpoint"

# Efficient — goes straight to the file
"update the stats logic in backend/internal/usecase/workout_usecase.go"
```

### Use `@` file mentions
In Claude Code, prefix a file path with `@` to attach it directly to your message without Claude needing to search for it:
```
@backend/internal/domain/workout.go add a `calories_burned` field
```

### Anchor k8s changes to the chart
```
# Vague — Claude might read all manifests
"add an env var to the backend"

# Targeted
"add FEATURE_FLAG=true to k8s/charts/jeeb-app/values.yaml under backend.env"
```

---

## Anti-Patterns (What Wastes Tokens)

| Anti-pattern | Why it's expensive | Better approach |
|---|---|---|
| "Check the whole codebase for X" | Reads dozens of files | Name the specific package/dir |
| "What's wrong with the backend?" | Open-ended = max context | Describe the symptom + file |
| Re-explaining context each message | Duplicates tokens | Let Claude remember within session |
| "Update everything to use the new pattern" | Broad scope | One file/package at a time |
| Asking Claude to summarize long logs | Logs are already verbose | Use `/logs <svc>` with grep |
| Attaching entire files when only one function matters | Sends full file content | Mention the function name + line range |

---

## Token Budget Rules

| Task size | Approach |
|---|---|
| Single function / small fix | Direct prompt with file path |
| Feature across 2–3 files | Use a slash command, name the files |
| New feature (handler → usecase → repo) | Use `/backend`, describe feature end-to-end |
| Cross-service change (backend + k8s + CI) | Split into separate tasks, `/clear` between them |
| Infrastructure audit | Use `/k8s` or spawn a subagent via the Agent tool |

---

## Good vs. Bad Prompt Examples

### Backend feature

**Bad:**
> Add a calories feature to the workout system

**Good:**
> `/backend` add `calories_burned int` to `domain/workout.go`, add it to the request struct in `adapter/in/http/workout_handler.go`, persist it in `adapter/out/workout_repo.go` — no migration needed (MongoDB is schemaless)

---

### Kubernetes change

**Bad:**
> The backend is crashing, fix it

**Good:**
> `/k8s` backend pod is OOMKilled — increase memory limit in `k8s/charts/jeeb-app/values.yaml` from 128Mi to 256Mi

---

### Frontend component

**Bad:**
> Add a chart to the dashboard

**Good:**
> `/frontend` add a bar chart to `src/pages/Dashboard.tsx` showing weekly workout count — use data from `useWorkouts` hook, Recharts is already installed

---

### Debugging

**Bad:**
> Why isn't Keycloak working?

**Good:**
> `/logs keycloak` — I'm seeing "invalid redirect_uri" in the logs, check `k8s/charts/jeeb-app/values.yaml` keycloak.redirectUris

---

## ccusage — Track Spend

```bash
ccusage          # today's token usage
ccusage --days 7 # last 7 days
```

Use this to understand where tokens go before trying to optimize further. High spend on a single session usually means broad scope or missing `/clear` calls.
