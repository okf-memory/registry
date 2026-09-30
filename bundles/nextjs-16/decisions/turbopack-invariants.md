---
type: Decision
title: Turbopack-Native Compilation Invariants
tags: [nextjs, turbopack, build, compilation, v16]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/turbopack-invariants"
---

# Turbopack-Native Compilation Invariants

Next.js 16 build configurations MUST adhere to Turbopack-native asset processing rules and eliminate legacy Webpack hooks.

## Invariants
1. **Turbopack as Default Compiler:** Production builds and development servers in Next.js 16 run on Turbopack by default (`next build --turbo` is the default pipeline). Configurations MUST NOT attempt to force webpack unless an incompatible legacy plugin is explicitly documented.
2. **Prohibition of Custom Webpack Plugins:** Custom `webpack: (config) => { ... }` blocks in `next.config.js` MUST NOT be used for CSS, SVG, or module aliasing. Native Turbopack options (e.g. `turbo: { rules: { ... } }`) or CSS module standards MUST be used instead.
3. **Rust-Based SWC Transform Invariants:** Custom Babel configurations (`.babelrc`, `babel.config.js`) are strictly prohibited as they disable Turbopack's native Rust compilation engine.

## Related
* [Mandatory Async Request APIs](async-request-apis.md)
* [Partial Prerendering (PPR) Boundaries](partial-prerendering.md)
