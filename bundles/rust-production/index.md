---
okf_version: "0.2"
bundle_version: "1.0.0"
title: "Rust Production Safety & Systems Best Practices"
description: "Authoritative architectural decisions for memory safety, zero-unsafe invariants, and structured async concurrency."
license: "MIT"
---

# Rust Production Knowledge Bundle

Curated architectural invariants and memory safety rules for high-throughput, production-grade Rust services and CLI tools.

## Decisions
* [Zero-Unsafe Invariant & Formal Audit Gate](decisions/zero-unsafe.md): Strictly forbids unreviewed unsafe blocks in production crates, enforcing safe abstraction wrappers.
* [Layered Error Handling: thiserror for Domains and anyhow for Applications](decisions/error-handling.md): Mandates strongly typed enum errors in internal crates and context-rich anyhow handling at system boundaries.
* [Async Tokio Runtime, Cooperative Yielding, and Task Boundaries](decisions/tokio-runtime.md): Governs multi-threaded Tokio runtime lifecycle, blocking workload delegation, and graceful shutdown.
