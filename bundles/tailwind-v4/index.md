---
okf_version: "0.2"
bundle_version: "1.0.0"
title: "Tailwind CSS v4 Architecture & CSS-First Config"
description: "Curated architectural decisions for Tailwind v4 CSS-first theme configuration and import discipline."
license: "MIT"
---

# Tailwind CSS v4 Architecture Seed Bundle

Authoritative configuration patterns and invariants reflecting the architecture shift in Tailwind CSS v4 from JavaScript configurations to CSS-native `@theme` and `@utility` directives.

## Decisions
* [CSS-First Theme Configuration](decisions/css-first-config.md): Replacement of tailwind.config.js with native CSS @theme blocks.
* [Modern Single-Line Imports](decisions/import-invariants.md): Requirement of @import "tailwindcss"; replacing deprecated @tailwind directives.
