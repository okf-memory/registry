---
type: Decision
title: CSS-First Theme Configuration
tags: [tailwind, css, frontend, styling, v4]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/css-first-config"
---

# CSS-First Theme Configuration

Tailwind CSS v4 projects MUST define theme tokens, custom fonts, and color extensions directly in CSS via `@theme` directives rather than JavaScript configuration files.

## Invariants
1. **Prohibition of `tailwind.config.js`:** Projects standardizing on Tailwind CSS v4 MUST NOT create or maintain a `tailwind.config.js` or `tailwind.config.ts` file. All configuration MUST reside in the primary CSS entrypoint.
2. **`@theme` Block Usage:** Custom design tokens, fonts, spacing, and palettes MUST be declared using standard CSS custom property syntax within an `@theme` block:
   ```css
   @theme {
     --color-primary: #8b5cf6;
     --font-sans: 'Outfit', sans-serif;
   }
   ```
3. **Custom Utilities via `@utility`:** Custom utility classes MUST be authored using `@utility <name>` blocks rather than legacy `@layer utilities` syntax.

## Related
* [Modern Single-Line Import Invariants](import-invariants.md)
