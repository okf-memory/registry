---
type: Decision
title: "Server Actions Security & Error Handling"
description: Requires input validation with Zod and sanitized error boundaries for all Server Actions.
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
governance: constraint
---

# Decision: Server Actions Security & Error Handling

## Context
Server Actions expose public HTTP POST endpoints under the hood. Unprotected Server Actions are vulnerable to unauthorized execution and sensitive internal stack traces leaking to clients.

## Decision
1. **Schema Validation:** Every Server Action MUST parse input arguments through schema validation before executing business logic.
2. **Sanitized Errors:** Unhandled exceptions must be caught and scrubbed; internal database errors must never be returned directly to the client.

```mermaid
flowchart TD
    ClientAction[Client Invocation] --> AuthCheck{Session Valid?}
    AuthCheck -->|No| RejectUnauthorized[Throw 401 Unauthorized]
    AuthCheck -->|Yes| ParseSchema{Input Schema Valid?}
    ParseSchema -->|No| ValidationFail[Return Field Errors]
    ParseSchema -->|Yes| Exec[Execute Business Logic]
```

## Related
* [App Router Caching Invariants in Next.js 15](caching.md)
