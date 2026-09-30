---
okf_version: "0.2"
bundle_version: "1.0.0"
title: "Next.js 15 Seed Memory Bundle"
description: "Curated architectural decisions, caching invariants, and security patterns for Next.js 15."
license: "MIT"
---

# Next.js 15 Knowledge Bundle

Curated best practices, App Router architectural decisions, and production patterns for Next.js 15 applications.

## Decisions
* [App Router Caching Invariants](decisions/caching.md): Establishes default uncached fetch requests and explicit cache tag opt-ins in Next.js 15.
* [Server Actions Security & Error Handling](decisions/server-actions.md): Requires input validation with Zod and sanitized error boundaries for all Server Actions.
