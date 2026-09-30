---
type: Decision
title: "Layered Error Handling: thiserror for Domains and anyhow for Applications"
description: Mandates strongly typed enum errors in internal crates and context-rich anyhow handling at system boundaries.
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
governance: constraint
---

# Decision: Layered Error Handling

## Context
Stringly-typed or unprincipled errors in Rust lead to brittle pattern matching, unhandled edge cases, and cryptic panic cascades under load.

## Decision
1. **Library & Core Domain Boundaries:** Domain crates MUST use `thiserror` to define exhaustive, strongly-typed error enums. Callers must be able to programmatically inspect failure modes.
2. **Application & Binary Roots:** Executable binaries (`main.rs`, CLI handlers, HTTP adapters) MUST use `anyhow::Result<T>` with chained `.context("action description")` calls for rich operational diagnostics.
3. **Panic Prohibition:** Production execution paths MUST NOT invoke `unwrap()`, `expect()`, or `panic!()` outside of compile-time constants or test fixtures.

```mermaid
flowchart LR
    Domain["Core Domain / Crates<br/>(thiserror typed enums)"] -->|Bubbled with ?| Adapter["Application Boundary / main<br/>(anyhow with .context())"]
    Adapter -->|Structured Log| Output["CLI / JSON Error Output"]
```

## Related
* [Zero-Unsafe Invariant & Formal Audit Gate](zero-unsafe.md)
