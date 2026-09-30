---
type: Decision
title: App Router Caching Invariants in Next.js 15
description: Establishes default uncached fetch requests and explicit cache tag opt-ins in Next.js 15.
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
governance: constraint
---

# Decision: Next.js 15 App Router Fetch Caching Defaults

## Context
In Next.js 14 and earlier, `fetch` requests inside Server Components were cached by default (`force-cache`). Next.js 15 changes this behavior: all `fetch` requests now default to `no-store` (uncached) unless explicitly configured otherwise.

## Decision
1. **Explicit Caching:** Whenever persistent HTTP responses are desired across requests, developers MUST explicitly pass `cache: 'force-cache'` or declare `next: { revalidate: <seconds>, tags: [...] }`.
2. **Deterministic Purging:** Dynamic data must be accompanied by tags for cache revalidation via `revalidateTag()`.

```mermaid
flowchart TD
    Req[Incoming Fetch Request] --> Check{Explicit Cache Opt-in?}
    Check -->|Yes: force-cache or tags| CacheLayer[Next.js Data Cache]
    Check -->|No: Default in Next.js 15| DirectOrigin[Direct Origin Request]
```

## Related
* [Server Actions Security & Error Handling](server-actions.md)
