---
type: Decision
title: Modern Single-Line Import Invariants
tags: [tailwind, css, frontend, imports, v4]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/import-invariants"
---

# Modern Single-Line Import Invariants

Tailwind CSS v4 entrypoints MUST use the modern `@import "tailwindcss";` directive and eliminate deprecated v3 directives.

## Invariants
1. **Single Import Directive:** CSS stylesheets MUST import the framework using the standard CSS `@import` statement:
   ```css
   @import "tailwindcss";
   ```
2. **Deprecated `@tailwind` Directives Prohibited:** The legacy three-part directives (`@tailwind base;`, `@tailwind components;`, `@tailwind utilities;`) MUST NOT be used in v4 projects.
3. **Source Detection Discipline:** The explicit `content: [...]` array from v3 is removed; Tailwind v4 automatically discovers template files. Explicit content paths, when required for monorepos or external node_modules, MUST be specified via `@source` directives in the stylesheet:
   ```css
   @source "../packages/ui";
   ```

## Related
* [CSS-First Theme Configuration](css-first-config.md)
