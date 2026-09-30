---
type: Decision
title: "Context Lifecycle & Goroutine Leak Avoidance"
description: Mandates explicit context passing as first parameter and deterministic termination of spawned goroutines.
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
governance: constraint
---

# Decision: Context Lifecycle & Goroutine Leak Avoidance

## Context
Goroutines in Go run indefinitely unless they return or terminate. Spawned goroutines without lifecycle bounds or context awareness lead to memory exhaustion and resource leaks.

## Decision
1. **First Parameter Convention:** Functions that perform blocking I/O or background processing MUST accept `ctx context.Context` as their first parameter.
2. **Deterministic Return on Done:** Goroutines listening on channels or performing work loops MUST select on `ctx.Done()` to exit promptly on cancellation.

```mermaid
flowchart TD
    Parent[Parent Context] --> Fork[Spawn Goroutine]
    Fork --> SelectLoop{select on Channels}
    SelectLoop -->|Event Received| Handle[Process Event]
    SelectLoop -->|ctx.Done()| Cleanup[Clean Up & Exit]
```

## Related
* [Bounded Concurrency via errgroup](errgroup-limits.md)
