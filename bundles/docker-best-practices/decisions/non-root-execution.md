---
type: Decision
title: Mandatory Non-Root User Execution
tags: [docker, security, container, hardening]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/non-root-execution"
---

# Mandatory Non-Root User Execution

Containers MUST NOT run as root (UID 0) in production runtime environments.

## Invariants
1. **Explicit Unprivileged User:** Dockerfiles MUST create and switch to a dedicated non-root user and group (e.g. `appuser:appgroup` with explicit UID/GID >= 10001) prior to defining the `ENTRYPOINT` or `CMD`.
2. **File Ownership Discipline:** Application binaries, configuration files, and assets MUST have ownership explicitly assigned using `COPY --chown=appuser:appgroup`.
3. **No Sudo or Privileged Escalation:** Production runtime images MUST NOT contain `sudo`, setuid binaries, or write access to `/etc` or system library paths.
4. **Port Binding Governance:** Applications MUST listen on unprivileged ports (>= 1024, e.g. 3000, 8080) rather than standard privileged ports (80, 443).

## Related
* [Minimal Multi-Stage Builds](multi-stage-minimal.md)
* [Proper PID 1 Signal Handling](signal-handling.md)
