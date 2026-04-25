# ADR-0002: MongoDB as Primary Database

## Status
Accepted

## Context
Need flexible schema for varied feature data (workouts, sleep, finance).

## Decision
Use MongoDB for all persistent data.

## Consequences
- (+) Flexible document schema per feature
- (+) Easy to evolve data models
- (+) Good Go driver support
- (-) No ACID transactions across collections (acceptable for this use case)
