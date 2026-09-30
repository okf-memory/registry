---
okf_version: "0.2"
bundle_version: "1.0.0"
title: "Next.js 16 Production Architecture & Invariants"
description: "Curated architectural decisions for async request APIs, Partial Prerendering (PPR), and Turbopack builds in Next.js 16."
license: "MIT"
---

# Next.js 16 Production Architecture Seed Bundle

Curated architectural invariants and breaking change governance for Next.js 16 applications, establishing mandatory async request access, Partial Prerendering (PPR) streaming boundaries, and Turbopack compilation discipline.

## Decisions
* [Mandatory Async Request APIs](decisions/async-request-apis.md): Strict requirement to await params, searchParams, cookies, and headers.
* [Partial Prerendering (PPR) Boundaries](decisions/partial-prerendering.md): Isolation of static prerender shells from dynamic streaming chunks via Suspense.
* [Turbopack-Native Compilation Invariants](decisions/turbopack-invariants.md): Strict separation from custom Webpack hooks and plugins.
