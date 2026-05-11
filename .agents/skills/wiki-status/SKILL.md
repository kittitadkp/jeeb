# Wiki Status — Audit Vault Health and Ingestion Delta

Show vault statistics, delta since last sync, and structural insights.

## Trigger

Invoked with `/wiki-status`. Optional: `/wiki-status --insights` for graph analysis.

## Step 1: Resolve Vault Path

Same protocol as wiki-update and wiki-query.

## Step 2: Read the Manifest

Read `<vault>/.manifest.json`. If missing, report "vault not initialized — run /wiki-setup".

Extract:
- Total pages tracked
- Last sync timestamp per project
- Last commit synced per project

## Step 3: Scan Current State

Count:
- Pages per category (concepts/, entities/, skills/, references/, synthesis/, journal/, projects/)
- Orphaned pages (no wikilinks pointing to them)
- Broken wikilinks (reference pages that don't exist)

## Step 4: Classify Delta

For each project in the manifest, check `git log <last_commit_synced>..HEAD --oneline` to count unsyced commits.

Classify:
- **Fresh** — 0 unsynced commits
- **Minor drift** — 1–10 commits (suggest `/wiki-update`)
- **Stale** — 11–50 commits (suggest `/wiki-update`)
- **Very stale** — 50+ commits (suggest `/wiki-update --full`)

## Step 5: Report

Output a status table:

```
VAULT: <path>
Last sync: <timestamp>

PAGES
  concepts/     <N> pages
  projects/     <N> pages (<list of project names>)
  skills/       <N> pages
  references/   <N> pages
  synthesis/    <N> pages
  journal/      <N> pages
  Total:        <N> pages

HEALTH
  Orphaned pages:   <N>
  Broken wikilinks: <N>

SYNC STATUS
  <project-name>: <status> (<N> unsynced commits since <SHA>)

RECOMMENDATION
  <action> — <reason>
```

## Insights Mode (`--insights`)

When `--insights` is passed, analyze the wikilink graph:

- **Anchor pages** — pages with the most inbound wikilinks (top 5)
- **Bridge pages** — pages connecting otherwise-isolated clusters
- **Orphan-adjacent** — pages one hop from a high-traffic page but with no inbound links themselves
- **Surprising connections** — cross-category wikilinks worth reviewing

Write results to `<vault>/_insights.md` and report a summary.
