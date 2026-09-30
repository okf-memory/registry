# Decisions
* [Context Lifecycle & Goroutine Leak Avoidance](context-propagation.md) - Mandates explicit context passing as first parameter and deterministic termination of spawned goroutines.
* [Bounded Concurrency via errgroup](errgroup-limits.md) - Requires bounding worker concurrency and unified error collection using errgroup.Group.
