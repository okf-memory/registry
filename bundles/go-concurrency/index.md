---
okf_version: "0.2"
bundle_version: "1.0.0"
title: "Go Concurrency Best Practices Seed Bundle"
description: "Curated architectural patterns for leak-free, bounded Go concurrency and context lifecycle."
license: "MIT"
---

# Go Concurrency Knowledge Bundle

Curated decisions and patterns for production Go concurrency, goroutine leak prevention, and structured context cancellation.

## Decisions
* [Context Lifecycle & Propagation](decisions/context-propagation.md): Mandates explicit context passing as first parameter and deterministic termination of spawned goroutines.
* [Bounded Concurrency via errgroup](decisions/errgroup-limits.md): Requires bounding worker concurrency and unified error collection using errgroup.Group.
