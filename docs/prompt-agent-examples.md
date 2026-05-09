# Prompt Agent Examples

These examples show how to ask `codex` to turn a short issue into a stronger prompt for another AI tool.

## Generic pattern

```text
codex "Read AGENTS.md, docs/ai-prompt-agent.md, and docs/prompt-templates/<template>.md. Convert the short issue below into a copy-paste-ready prompt for <target tool>. Keep it practical, prefer minimal safe changes, and output the final prompt only.

Issue:
Problem: ...
Where: ...
Expected: ...
Actual: ...
Evidence: ...
Constraints: ..."
```

## Keycloak redirect example

```text
codex "Read AGENTS.md, docs/ai-prompt-agent.md, and docs/prompt-templates/auth-debug.md. Convert the short issue below into a copy-paste-ready debugging prompt for Claude Code. Output only the final prompt.

Issue:
Problem: Keycloak redirect is broken.
Where: http://jeeb-dev.local/
Expected: App loads and redirects to Keycloak login.
Actual: Browser shows 'This site can’t be reached'.
Evidence: Failure happens before any login page is shown.
Constraints: Prefer minimal safe changes."
```

## Kubernetes ingress example

```text
codex "Read AGENTS.md, docs/ai-prompt-agent.md, and docs/prompt-templates/k8s-debug.md. Turn the issue below into a debugging prompt for Cursor. Output one prompt only.

Issue:
Problem: Ingress host returns 404 in dev.
Where: Kubernetes ingress in dev namespace.
Expected: Request reaches frontend service.
Actual: Ingress returns 404, but service works with port-forward.
Evidence: Host resolves and ingress controller is running.
Constraints: Avoid broad Helm chart rewrite."
```

## Go backend bug example

```text
codex "Read AGENTS.md, docs/ai-prompt-agent.md, and docs/prompt-templates/go-debug.md. Generate a debugging prompt for Codex CLI from the short issue below. Output the prompt only.

Issue:
Problem: Create workout endpoint returns 500.
Where: backend POST /api/workouts.
Expected: Request saves workout and returns success.
Actual: 500 response.
Evidence: Log mentions nil pointer in workout usecase.
Constraints: Keep Clean Architecture boundaries and prefer a small fix."
```

## Vue bug example

```text
codex "Read docs/ai-prompt-agent.md and docs/prompt-templates/vue-debug.md. Turn the issue below into a prompt for ChatGPT. The final prompt should request root cause analysis, a minimal fix, and validation steps. Output the prompt only.

Issue:
Problem: Vue profile page freezes after save.
Where: browser UI on profile settings page.
Expected: Save completes and form is usable.
Actual: Spinner never stops.
Evidence: Console shows unhandled promise rejection.
Constraints: Do not rewrite the whole page."
```

## MongoDB performance example

```text
codex "Read docs/ai-prompt-agent.md and docs/prompt-templates/mongo-debug.md. Expand the short issue below into a practical debugging prompt for Claude Code. Output the prompt only.

Issue:
Problem: MongoDB search query is too slow.
Where: task search endpoint.
Expected: Under 200ms.
Actual: 4 to 8 seconds.
Evidence: explain output shows COLLSCAN.
Constraints: Prefer query or index fixes before schema redesign."
```

## Centrifugo timeout example

```text
codex "Read docs/ai-prompt-agent.md and docs/prompt-templates/centrifugo-debug.md. Create a debugging prompt for Cursor from the issue below. Output only the final prompt.

Issue:
Problem: Centrifugo RPC calls time out during bursts.
Where: real-time notification flow.
Expected: RPC completes inside timeout budget.
Actual: intermittent timeout errors.
Evidence: logs show RPC timeout followed by retries.
Constraints: Keep client-facing behavior stable."
```

## Docs rewrite example

```text
codex "Read AGENTS.md, docs/ai-prompt-agent.md, and docs/prompt-templates/docs-rewrite.md. Turn the short request below into a prompt for ChatGPT that rewrites documentation clearly and accurately. Output only the prompt.

Issue:
Problem: Deployment guide is hard to follow.
Where: internal docs page.
Expected: Short, accurate, ordered setup steps.
Actual: Repetitive wording and missing prerequisites.
Evidence: New developer could not finish setup.
Constraints: Keep technical meaning, remove fluff, and stay aligned with the repo."
```
