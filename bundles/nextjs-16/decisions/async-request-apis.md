---
type: Decision
title: Mandatory Async Request APIs
tags: [nextjs, react, app-router, async, v16]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/async-request-apis"
---

# Mandatory Async Request APIs

Next.js 16 components, route handlers, and server logic MUST asynchronously await all request-time data structures.

## Invariants
1. **Asynchronous Page & Layout Props:** Page and layout props `params` and `searchParams` are asynchronous Promises in Next.js 16. Components MUST `await params` and `await searchParams` before accessing properties:
   ```typescript
   export default async function Page({
     params,
   }: {
     params: Promise<{ slug: string }>;
   }) {
     const { slug } = await params;
     return <h1>Post: {slug}</h1>;
   }
   ```
2. **Synchronous Property Access Prohibited:** Direct synchronous property extraction (`props.params.slug`) is strictly prohibited and throws runtime compilation errors.
3. **Async Header & Cookie Stores:** Server-side functions reading headers or cookies MUST await the store instances:
   ```typescript
   import { cookies, headers } from 'next/headers';
   const cookieStore = await cookies();
   const headersList = await headers();
   ```
4. **Client Component Separation:** Client components needing search parameters MUST use `useSearchParams()` from `next/navigation` wrapped inside a `<Suspense>` boundary.

## Related
* [Partial Prerendering (PPR) Boundaries](partial-prerendering.md)
* [Turbopack-Native Compilation Invariants](turbopack-invariants.md)
