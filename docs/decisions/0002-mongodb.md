# ADR 0002: MongoDB for Application Storage

## Status

Accepted

## Context

Both backends persist document-shaped user data, event history, and flexible feature records. The current codebase already uses the MongoDB Go driver and Mongo-backed repositories throughout.

## Decision

Keep MongoDB as the storage layer for both application domains.

## Consequences

- Collections map closely to feature aggregates
- Seed flows are simple document inserts
- Pagination is implemented at the repository layer with shared request metadata
