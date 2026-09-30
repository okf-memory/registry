---
type: Decision
title: "Proper PID 1 Signal Handling & Process Reaping"
tags: [docker, lifecycle, reliability, signals]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/signal-handling"
---

# Proper PID 1 Signal Handling & Process Reaping

Container entrypoints MUST correctly receive and forward OS lifecycle signals (`SIGTERM`, `SIGINT`) and reap orphaned zombie child processes.

## Invariants
1. **Exec Form Requirement:** `ENTRYPOINT` and `CMD` instructions MUST use the JSON exec form (`["executable", "param1", "param2"]`). Shell wrapper syntax (`CMD "node index.js"`) is prohibited as it spawns `/bin/sh -c` as PID 1, which ignores `SIGTERM`.
2. **Init Process Requirement for Non-Compliant Runtimes:** Runtimes that do not natively handle PID 1 process reaping or signal forwarding (e.g. Node.js, Python scripts) MUST execute behind a dedicated init supervisor (`tini` or `dumb-init`) or use Docker's `--init` flag.
3. **Graceful Shutdown Deadline:** Applications MUST listen for `SIGTERM`, cease accepting incoming traffic, drain in-flight requests, and exit within the container runtime grace period (default 10s–30s) before `SIGKILL` escalation.

## Related
* [Mandatory Non-Root User Execution](non-root-execution.md)
* [Minimal Multi-Stage Builds](multi-stage-minimal.md)
