# ADR-0001: Clean Architecture

## Status
Accepted

## Context
Need maintainable, testable codebase with clear boundaries.

## Decision
Use Clean Architecture with hexagonal/ports-and-adapters pattern.

## Consequences
- (+) Easy to test (mock ports)
- (+) Swap implementations without changing business logic
- (+) Clear dependency direction
- (-) More boilerplate initially
