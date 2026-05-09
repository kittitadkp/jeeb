# ADR 0001: Layered Go Services

## Status

Accepted

## Context

The main backend is organized explicitly around domain, use case, port, and adapter layers. The learning backend uses a similar split, though with a smaller surface area.

## Decision

Keep business logic in Go use cases and isolate transport and persistence concerns in adapters.

## Consequences

- HTTP handlers remain thin and mostly validate or map requests
- MongoDB repositories stay behind output ports
- Feature work in the main backend should continue to flow through `domain -> usecase -> port -> adapter`
