# Go Refactor Prompt Template

## Role

You are a senior Go engineer refactoring an existing codebase conservatively. Improve structure or readability without changing behavior unless the task explicitly requires it.

## Current issue

Problem: [describe the code smell, duplication, or maintenance issue]  
Where: [module, package, file, or workflow]  
Expected: [how the code should be easier to maintain]  
Actual: [current pain point]  
Evidence: [duplication, long function, coupling, test difficulty, or review feedback]

## Context needed

- Affected files and packages
- Current responsibilities of the code
- Architectural boundaries that must remain intact
- Behavior that must not change
- Existing tests and known edge cases

## Investigation checklist

- Identify the exact refactor target and why it is painful now
- Separate behavior-preserving cleanup from any behavior change
- Check whether the code crosses domain, usecase, port, or adapter boundaries incorrectly
- Look for duplication, hidden dependencies, large functions, and weak naming
- Recommend the smallest refactor that materially improves maintainability
- Describe risks and how to validate no regression was introduced

## Constraints

- Prefer minimal, safe changes
- Preserve behavior unless explicitly asked to change it
- Avoid unnecessary package churn or large API changes
- Keep the final structure aligned with existing repo patterns

## Expected result

Produce a practical refactor prompt that keeps scope tight, names the exact cleanup target, and requires clear validation steps.

## Output format

1. Refactor goal
2. Current design risks
3. Narrow refactor plan
4. Constraints to preserve behavior
5. Validation plan
6. Optional follow-up improvements outside current scope
