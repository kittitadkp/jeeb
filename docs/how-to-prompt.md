# How to Prompt

## Purpose

Use this guide when you need to ask `codex`, Claude Code, Cursor, or ChatGPT for help on an engineering task.

The goal is simple: give enough context for the AI to debug or improve the right thing without wasting time on guesses.

## Good prompt structure

Write prompts in this order:

```text
Task: what you want the AI to do
Context: where the problem lives
Expected: what should happen
Actual: what happens now
Evidence: logs, errors, URLs, status codes, screenshots, or failing commands
Constraints: limits such as minimal safe change, no schema change, docs only, no broad refactor
Output: what kind of answer you want back
```

## Short template

```text
Task: Help debug this issue.
Context: backend/frontend/k8s/jenkins/docs
Expected: [what should happen]
Actual: [what happens instead]
Evidence: [exact error or observed behavior]
Constraints: prefer minimal safe changes
Output: give ranked likely causes, checks to run, minimal fix, and validation steps
```

## Better prompt example

```text
Task: Help debug a Keycloak redirect issue.
Context: local dev at http://jeeb-dev.local/
Expected: app loads and redirects unauthenticated users to Keycloak login
Actual: browser shows 'This site can’t be reached'
Evidence: the failure happens before any login page appears
Constraints: prefer minimal safe changes and check DNS, ingress, proxy, and auth config before suggesting rewrites
Output: give ranked likely causes, exact checks to run, a minimal fix, and how to validate it
```

## Bad prompt example

```text
Fix Keycloak.
```

This is weak because it does not say:

- what is broken
- where it happens
- what the expected behavior is
- what evidence already exists
- what kind of answer is needed

## Prompt writing rules

- Be specific about the failing page, service, endpoint, host, job, or namespace.
- Include the exact error text when you have it.
- Separate expected behavior from actual behavior.
- Add constraints so the AI does not overreach.
- Ask for root cause analysis before code changes.
- Ask for validation steps after the proposed fix.
- Remove secrets, tokens, passwords, and private credentials.

## Good output requests

Ask the AI to return one of these:

- a debugging plan
- ranked likely causes
- a minimal fix
- validation steps
- a copy-paste-ready prompt for another AI tool

## Quick examples

### Go backend bug

```text
Task: Debug a Go API bug.
Context: backend POST /api/workouts
Expected: request saves workout successfully
Actual: API returns 500
Evidence: log mentions nil pointer in usecase
Constraints: keep Clean Architecture boundaries and prefer a small fix
Output: identify root cause, propose minimal fix, and list tests to run
```

### Kubernetes ingress issue

```text
Task: Debug an ingress routing issue.
Context: dev namespace ingress host
Expected: request reaches frontend service
Actual: ingress returns 404
Evidence: service works with port-forward
Constraints: avoid broad Helm changes
Output: ranked likely causes, kubectl checks, minimal manifest fix, and validation steps
```

### MongoDB performance issue

```text
Task: Debug a slow MongoDB query.
Context: task search endpoint
Expected: query completes under 200ms
Actual: takes 4 to 8 seconds
Evidence: explain output shows COLLSCAN
Constraints: try query or index fixes before schema redesign
Output: likely causes, explain-based checks, index or query fix, and validation plan
```

## Repo-specific note

For this repository:

- `backend/` is Go
- `frontend/` is the main UI
- `k8s/` contains cluster manifests
- `jenkins/` contains CI/CD pipelines
- `docs/` contains runbooks and architecture notes

If the task needs stronger structure, use [ai-prompt-agent.md](./ai-prompt-agent.md) and the templates under [prompt-templates](./prompt-templates/).
