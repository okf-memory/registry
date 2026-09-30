---
type: Decision
title: Minimal Multi-Stage Container Builds
tags: [docker, performance, security, build]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/multi-stage-minimal"
---

# Minimal Multi-Stage Container Builds

Production container images MUST isolate compilation toolchains from final runtime artifacts.

## Invariants
1. **Multi-Stage Separation:** Dockerfiles MUST define at least two stages: a builder stage containing compilers/package managers (Go SDK, Node.js toolchain, Rust Cargo, Maven) and a minimal runtime stage containing only the compiled binary and essential runtime dependencies.
2. **Runtime Base Image Selection:** Runtime images MUST use minimal base distributions (`gcr.io/distroless/static`, `alpine`, or `scratch` for statically linked binaries). General-purpose development distros (e.g. `ubuntu:latest`, `node:alpine` with package managers left intact) MUST NOT be deployed to production.
3. **Artifact Hygiene:** Package manager caches (`npm cache clean`, `rm -rf /var/cache/apk/*`) and package manifests used only during installation MUST NOT leak into runtime layers.
4. **Secret Leakage Prohibition:** Build arguments (`ARG`) and environment variables containing API tokens, private SSH keys, or registry credentials MUST NOT be written to image layers; `RUN --mount=type=secret` MUST be utilized instead.

## Related
* [Mandatory Non-Root User Execution](non-root-execution.md)
* [Proper PID 1 Signal Handling](signal-handling.md)
