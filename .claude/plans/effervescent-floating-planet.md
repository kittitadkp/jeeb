# Claude Code Token & Context Optimization

## Context
RTK hook is already wired up. Now implementing the remaining suggestions to reduce token usage and improve session quality: install ccusage for visibility, strengthen .claudeignore, trim CLAUDE.md, enable auto-compact, and write a comprehensive prompt efficiency guide for the project.

---

## Steps

### 1. Install ccusage
- [x] `npm install -g ccusage`
- Tracks token spend per session/day — baseline visibility before further tuning
- Verify: `ccusage --version`

### 2. Strengthen .claudeignore
**File:** `D:\personal\jeeb\.claudeignore`

Already exists and covers basics. Add missing high-noise entries:
- `go.sum` — large, never read by Claude
- `package-lock.json` / `yarn.lock` / `pnpm-lock.yaml` — lock files
- `charts/*/charts/` — Helm dependency charts (downloaded, not authored)
- `*.tgz` — Helm chart tarballs
- `.helm/` — Helm cache

### 3. Trim project CLAUDE.md
**File:** `D:\personal\jeeb\CLAUDE.md`

Remove or shorten sections inferable from code:
- Collapse verbose NodePort table — keep as minimal lookup reference
- Keep: architecture decisions, hex layout, slash command table, troubleshooting rules
- Target: under 120 lines

### 4. Enable auto-compact in settings
**File:** `C:\Users\kitti\.claude\settings.json`

Add `"autoCompactEnabled": true` so context auto-summarizes before hitting the limit.

### 5. Write prompt efficiency guide
**File:** `D:\personal\jeeb\docs\prompt-efficiency.md`

Comprehensive guide covering:
- **General principles** — task framing, specificity, scope control
- **Context management** — when to `/clear`, `/compact`, use subagents
- **RTK usage** — what it does, `rtk gain`, `rtk discover`
- **Slash commands** — which command for which task (from CLAUDE.md table)
- **File targeting** — how to give Claude exact paths vs. broad searches
- **Anti-patterns** — what makes prompts expensive (vague scope, re-explaining context, asking Claude to "check everything")
- **Token budget rules** — thresholds for when to split tasks
- **Examples** — good vs. bad prompt pairs for this specific project (backend feature, k8s change, frontend component)

---

## Verification
- `ccusage` — shows session stats after next Claude Code session
- `.claudeignore` — confirm new entries present
- CLAUDE.md line count reduced
- `~/.claude/settings.json` — confirm autoCompactEnabled: true
- `docs/prompt-efficiency.md` — readable, actionable, project-specific
