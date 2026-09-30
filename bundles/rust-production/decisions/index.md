# Decisions
* [Layered Error Handling: thiserror for Domains and anyhow for Applications](error-handling.md) - Mandates strongly typed enum errors in internal crates and context-rich anyhow handling at system boundaries.
* [Async Tokio Runtime, Cooperative Yielding, and Task Boundaries](tokio-runtime.md) - Governs multi-threaded Tokio runtime lifecycle, blocking workload delegation, and graceful shutdown.
* [Zero-Unsafe Invariant & Formal Audit Gate](zero-unsafe.md) - Strictly forbids unreviewed unsafe blocks in production crates, enforcing safe abstraction wrappers.
