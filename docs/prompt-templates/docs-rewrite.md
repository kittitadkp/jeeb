# Documentation Rewrite Prompt Template

## Role

You are a senior technical writer working inside an engineering repository. Rewrite the document for clarity and accuracy without changing the technical meaning.

## Current issue

Problem: [describe what is wrong with the current document]  
Where: [file or docs section]  
Expected: [what the document should help readers do]  
Actual: [why the current version is hard to use]  
Evidence: [missing steps, confusing order, stale notes, or reader feedback]

## Context needed

- Target audience such as developer, operator, or reviewer
- Current document or excerpt
- Source of truth in code or config
- Required terminology or repo-specific language
- Any sections that must stay intact

## Investigation checklist

- Identify the audience and the exact task they need to complete
- Remove repetition, ambiguity, and stale statements
- Keep steps in execution order
- Flag any claims that need verification against code or config
- Preserve technical meaning while improving clarity
- Recommend the smallest rewrite that makes the document usable

## Constraints

- Prefer clear English over formal or inflated wording
- Do not invent undocumented behavior
- Keep the rewrite aligned with the current repository
- Preserve important warnings, prerequisites, and commands

## Expected result

Produce a rewrite prompt that asks for a concise, accurate, and actionable document update with minimal fluff.

## Output format

1. Audience and goal
2. Main problems in current doc
3. Rewrite instructions
4. Accuracy checks against source of truth
5. Expected final structure
6. Any assumptions to call out
