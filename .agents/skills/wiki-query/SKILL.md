# Wiki Query — Search and Synthesize from Obsidian Vault

Answer questions using the pre-compiled Obsidian wiki rather than re-reading raw code.

## Trigger

Invoked with `/wiki-query <question>`. Optional modes:
- `/wiki-query --fast <question>` — index-only, answers from summaries only
- `/wiki-query --public <question>` — exclude pages tagged `visibility/internal` or `visibility/pii`

## Step 1: Resolve Vault Path

Same as wiki-update: walk up for `.env` with `OBSIDIAN_VAULT_PATH`, fall back to `~/.obsidian-wiki/config`, then current directory if `.obsidian/` exists.

## Step 2: Prime Context

Read in order (stop if answer is already clear):
1. `<vault>/hot.md` — recent activity context
2. `<vault>/index.md` — full page catalog with summaries

## Step 3: Tiered Retrieval (cheapest first)

**Tier 1 — Index pass:** Search page titles, tags, aliases, and `summary` fields in frontmatter via grep. Cost: minimal. If `--fast` flag, stop here and answer from summaries.

**Tier 2 — Section pass:** Grep for keywords in top 5 candidate pages. Read matching sections and surrounding context (±5 lines).

**Tier 3 — Full page read:** Last resort. Read up to 3 complete pages. Follow one level of `[[wikilinks]]` if they look relevant.

If the answer isn't in the wiki after Tier 3, say so explicitly — do not fall back to reading raw source code unless the user asks.

## Trust Annotations

Apply these annotations to cited pages:
- `lifecycle: archived` → note the successor page
- `lifecycle: disputed` → flag answer as uncertain
- Last modified >90 days ago → warn "may be stale, verify"
- No `lifecycle` field → treat as draft, warn if last modified >60 days

## Answer Format

```
<Synthesized answer with [[wikilinks]] to sources>

**Sources consulted:** [[page1]], [[page2]]
**Retrieved via:** Tier N (index / section / full read)
**Gaps:** <what the wiki doesn't cover about this question, if any>
```

Log the query to `<vault>/log.md`:
`<timestamp> | wiki-query | "<query text>" | tier=<N> | pages=<count>`
