# Wiki Update — Sync Project Knowledge to Obsidian Vault

Sync the current project's architectural decisions and patterns into the Obsidian wiki vault.

## Trigger

Invoked with `/wiki-update`. May include optional scope: `/wiki-update backend` or `/wiki-update --full` (ignore manifest, reprocess everything).

## Step 1: Resolve Vault Path

Walk up from current directory looking for `.env` containing `OBSIDIAN_VAULT_PATH`. Fall back to `~/.obsidian-wiki/config`. If neither exists, use the current working directory if it contains `.obsidian/`.

## Step 2: Understand the Project

Scan:
- `README.md`, `CLAUDE.md` for purpose and conventions
- `git log --oneline -50` for recent work
- Directory structure (one level deep) for architecture shape
- Existing wiki pages under `projects/<name>/` to avoid duplication

Derive project name from the directory name.

## Step 3: Compute the Delta

Read `<vault>/.manifest.json`. If absent or `--full` flag given, process everything.

Otherwise only process commits and files changed since `last_commit_synced`.

## Step 4: Decide What to Distill

**Worth capturing:**
- Architectural decisions and the reasoning behind them
- Patterns that repeat across the codebase
- Non-obvious constraints (why X wasn't used, why Y was chosen)
- Integration details that took effort to figure out
- Lessons learned from bugs or incidents
- Key abstractions and what problem they solve

**Skip:**
- File listings or routine CRUD boilerplate
- Version numbers and dependency details
- Anything the code already explains clearly
- Routine bug fixes without architectural significance

Guiding question: *"What would I need to know returning to this project in 3 months with zero context?"*

## Step 5: Write Wiki Pages

Organize under `<vault>/projects/<project-name>/`:
- `overview.md` — purpose, stack, architecture summary (always create/update this)
- One page per major concept or decision (e.g. `auth-flow.md`, `hexagonal-architecture.md`)

Global patterns or concepts that apply beyond this project go in `<vault>/concepts/` or `<vault>/skills/`.

Every page requires this YAML frontmatter:
```yaml
---
title: <Title>
category: projects | concepts | skills | references | synthesis
tags: [<relevant-tags>]
summary: >
  One or two sentences. Max 200 characters.
sources: [<file-paths-or-commits>]
provenance:
  extracted: 70   # visible in code/docs
  inferred: 25    # why decisions were made
  ambiguous: 5    # uncertain
lifecycle: draft  # draft | reviewed | verified | disputed | archived
updated: <YYYY-MM-DD>
---
```

Use `[[wikilink]]` notation to reference related pages.

## Step 6: Cross-link

After writing pages, scan for concepts mentioned in new pages and add wikilinks to existing pages that reference the same concept.

## Step 7: Update Tracking Files

Update these files in the vault root:

**`.manifest.json`** — append or update entry:
```json
{
  "last_synced": "<ISO timestamp>",
  "last_commit_synced": "<git SHA>",
  "projects": {
    "<project-name>": {
      "pages": ["<list of pages created/updated>"],
      "synced_at": "<ISO timestamp>"
    }
  }
}
```

**`index.md`** — add/update the project entry with one-line summary and links to its pages.

**`log.md`** — append one line: `<timestamp> | wiki-update | <project-name> | <N pages created> created, <M pages updated> updated`

**`hot.md`** — rewrite as a ~300-word semantic summary of recent wiki activity (last 5 syncs from log.md). This is used by wiki-query as a fast-path context primer.

## Output to User

Report:
- Pages created (list with paths)
- Pages updated (list with paths)
- Pages merged into existing (list)
- What was skipped and why
