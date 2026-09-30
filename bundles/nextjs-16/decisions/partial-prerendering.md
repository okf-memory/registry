---
type: Decision
title: Partial Prerendering (PPR) Boundaries
tags: [nextjs, react, ppr, streaming, performance, v16]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/partial-prerendering"
---

# Partial Prerendering (PPR) Boundaries

Next.js 16 pages leveraging Partial Prerendering MUST isolate dynamic, personalized data fetching inside explicit `<Suspense>` boundaries.

## Invariants
1. **Static Shell Preservation:** Top-level layouts and static page frames MUST NOT directly invoke dynamic request functions (`cookies()`, uncached database calls) outside of `<Suspense>` boundaries. The static shell must be prerenderable at build time.
2. **Mandatory Suspense Enclosure:** Any component performing runtime dynamic fetching (e.g. user profile cart, personalized recommendations) MUST be wrapped in a React `<Suspense fallback={<Skeleton />}>` boundary to allow instant edge shell streaming.
3. **PPR Configuration:** Routes utilizing partial prerendering MUST declare experimental or stable PPR enablement in `next.config.js` or the route configuration:
   ```typescript
   export const experimental_ppr = true;
   ```
4. **Fallback Determinism:** Suspense fallbacks MUST render deterministic placeholder HTML/CSS and MUST NOT initiate secondary client-side request cascades.

## Related
* [Mandatory Async Request APIs](async-request-apis.md)
* [Turbopack-Native Compilation Invariants](turbopack-invariants.md)
