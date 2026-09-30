---
okf_version: "0.2"
bundle_version: "1.0.0"
title: "Docker Production Hardening & Best Practices"
description: "Curated architectural decisions for minimal, unprivileged, signal-safe container images."
license: "MIT"
---

# Docker Production Hardening Seed Bundle

Curated container invariants and production patterns for secure, deterministic, and signal-safe container images across all runtime environments.

## Decisions
* [Non-Root User Execution](decisions/non-root-execution.md): Mandatory unprivileged user execution to prevent container breakout exploits.
* [Minimal Multi-Stage Builds](decisions/multi-stage-minimal.md): Separation of compilation toolchains from runtime environments using distroless or minimal bases.
* [Proper PID 1 Signal Handling](decisions/signal-handling.md): Signal forwarding and zombie process reaping via tini, dumb-init, or exec form.
