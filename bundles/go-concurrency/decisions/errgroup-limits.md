---
type: Decision
title: Bounded Concurrency via errgroup
description: Requires bounding worker concurrency and unified error collection using errgroup.Group.
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
governance: constraint
---

# Decision: Bounded Concurrency via errgroup

## Context
Unbounded goroutine spawning (`for item := range items { go process(item) }`) creates denial-of-service pressure on OS threads, network sockets, and file descriptors.

## Decision
1. **Bounded Parallelism:** When executing collections of concurrent tasks, use `golang.org/x/sync/errgroup` with `g.SetLimit(N)` or worker pool channels.
2. **First-Error Fail-Fast:** Concurrency errors must be captured deterministically and short-circuit sibling tasks.

```mermaid
flowchart TD
    Batch[Task Batch] --> Pool[errgroup.SetLimit N]
    Pool --> Worker1[Worker 1]
    Pool --> Worker2[Worker 2]
    Pool --> WorkerN[Worker N]
    Worker1 -->|Error| Cancel[Cancel Group Context]
    Cancel --> SiblingHalt[Halt Sibling Workers]
```

## Related
* [Context Lifecycle & Goroutine Leak Avoidance](context-propagation.md)
