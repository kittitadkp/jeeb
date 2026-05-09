# AI Prompt Agent

## Purpose

This prompt agent turns a short engineering issue into a stronger prompt for `codex`, Claude Code, Cursor, or ChatGPT.

Use it when the raw problem statement is too short, vague, or missing the structure another AI tool needs to debug, refactor, review, or rewrite something effectively.

## When to use it

- A bug report is only one or two lines.
- You want another AI tool to investigate a production-like issue faster.
- The problem crosses multiple layers such as frontend, auth, backend, MongoDB, or Kubernetes.
- You want a prompt that asks for root cause analysis, not random guesses.
- You need a copy-paste prompt with constraints, expected result, and output format.

Do not use it for trivial requests such as "rename this variable" or "explain this function."

## How to write a short problem statement

A good short problem statement is usually 3 to 6 lines and includes:

- What is broken
- Where it happens
- What you expected
- What actually happened
- Any hard evidence already known
- Any constraint that matters

Use this shape:

```text
Problem: [one sentence]
Where: [service, page, URL, job, command, namespace]
Expected: [what should happen]
Actual: [what happens instead]
Evidence: [error text, logs, status code, screenshot note]
Constraints: [safe change, no schema change, no downtime, docs only, etc.]
```

## Basic workflow

1. Pick the closest template from `docs/prompt-templates/`.
2. Write the short issue statement.
3. Tell the AI tool to expand that short issue into a full prompt using the selected template.
4. Review the generated prompt.
5. Add any missing local context before sending it to another AI tool.

## Prompt command examples

### Codex CLI

```text
codex "Read AGENTS.md, docs/ai-prompt-agent.md, and docs/prompt-templates/auth-debug.md. Turn this short issue into a copy-paste-ready debugging prompt for Codex CLI. Keep it practical, ask for missing context only if critical, and output the final prompt only.

Issue:
Problem: Keycloak login redirect is broken.
Where: http://jeeb-dev.local/
Expected: App should open and redirect to Keycloak login if session is missing.
Actual: Browser shows 'This site can’t be reached'.
Evidence: Happens before login page loads.
Constraints: Prefer minimal safe changes."
```

### Claude Code

```text
Use AGENTS.md, docs/ai-prompt-agent.md, and docs/prompt-templates/go-debug.md.
Convert this short issue into a strong debugging prompt for Claude Code.
Return one final prompt I can paste into Claude Code.

Problem: Go API returns 500 on create workout
Where: backend service, POST /api/workouts
Expected: workout saved
Actual: 500 error
Evidence: log mentions nil pointer in usecase
Constraints: prefer minimal safe change, keep Clean Architecture boundaries
```

### Cursor

```text
Read docs/ai-prompt-agent.md and docs/prompt-templates/vue-debug.md.
Expand this into a high-quality prompt for Cursor.
The result should tell Cursor what to inspect, which evidence to request, and how to present a fix.

Problem: Vue settings page freezes after save
Where: settings form in browser
Expected: save succeeds and form stays interactive
Actual: UI freezes and spinner never ends
Evidence: browser console shows promise rejection
Constraints: do not rewrite the whole page
```

### ChatGPT

```text
Use docs/ai-prompt-agent.md and docs/prompt-templates/mongo-debug.md.
Turn the short issue below into a practical debugging prompt for ChatGPT.
The final prompt must include context to gather, ranked likely causes, commands or checks to run, and a minimal-fix mindset.

Problem: MongoDB query is slow on task search
Where: task list endpoint
Expected: under 200ms
Actual: 4 to 8 seconds
Evidence: high CPU during regex search
Constraints: avoid risky schema changes first
```

## Issue examples

### Keycloak redirect issue

Short issue:

```text
Problem: Keycloak login redirect is broken.
Where: http://jeeb-dev.local/
Expected: app opens and sends unauthenticated users to Keycloak.
Actual: browser says 'This site can’t be reached'.
Evidence: failure happens before any login page appears.
Constraints: prefer minimal safe changes.
```

Recommended template: `docs/prompt-templates/auth-debug.md`

### Kubernetes ingress issue

Short issue:

```text
Problem: App is unreachable through ingress.
Where: dev namespace ingress host
Expected: host routes to frontend service
Actual: 404 or timeout from ingress
Evidence: service works with port-forward
Constraints: avoid broad chart rewrite
```

Recommended template: `docs/prompt-templates/k8s-debug.md`

### Go backend bug

Short issue:

```text
Problem: API returns 500 when creating a workout.
Where: backend POST /api/workouts
Expected: request is validated and saved
Actual: 500 response
Evidence: panic or nil pointer in logs
Constraints: keep Clean Architecture boundaries
```

Recommended template: `docs/prompt-templates/go-debug.md`

### Vue frontend bug

Short issue:

```text
Problem: Save button leaves page stuck in loading state.
Where: Vue profile page
Expected: success toast and enabled form
Actual: spinner never stops
Evidence: console shows failed promise
Constraints: avoid rewriting component tree
```

Recommended template: `docs/prompt-templates/vue-debug.md`

### MongoDB query or index issue

Short issue:

```text
Problem: Search endpoint is too slow.
Where: MongoDB-backed task search
Expected: sub-second query time
Actual: several seconds under normal load
Evidence: COLLSCAN in explain output
Constraints: try index or query fixes before data model changes
```

Recommended template: `docs/prompt-templates/mongo-debug.md`

### Centrifugo RPC timeout issue

Short issue:

```text
Problem: Centrifugo RPC calls time out intermittently.
Where: real-time notification flow
Expected: RPC completes within timeout budget
Actual: timeout errors during bursts
Evidence: app logs show RPC timeout and retry noise
Constraints: keep behavior stable for existing clients
```

Recommended template: `docs/prompt-templates/centrifugo-debug.md`

## Good vs bad prompts

Bad:

```text
Fix Keycloak.
```

Why it is bad:

- No symptom
- No target environment
- No expected behavior
- No evidence
- No constraints

Better:

```text
Problem: Users cannot reach the app at http://jeeb-dev.local/.
Expected: App loads and redirects unauthenticated users to Keycloak login.
Actual: Browser shows 'This site can’t be reached' before any Keycloak page appears.
Evidence: Issue happens in local dev through browser, not during API call.
Constraints: Prefer minimal safe changes and investigate DNS, ingress, reverse proxy, and auth config in that order.
```

Why it is better:

- The failure is observable
- The affected URL is explicit
- The expected auth behavior is clear
- The evidence narrows the layer
- The constraint reduces noisy advice

## Checklist before sending a prompt to another AI tool

- Did you choose the right template?
- Did you include the exact error text or URL?
- Did you say what should happen versus what actually happens?
- Did you mention the affected module or layer?
- Did you include constraints such as safe change, no downtime, or docs only?
- Did you remove secrets, tokens, passwords, and internal-only credentials?
- Did you ask for ranked likely causes instead of a blind fix?
- Did you ask for validation steps after the proposed fix?

If those answers are mostly yes, the prompt is ready.
