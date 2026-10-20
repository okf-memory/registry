---
okf_version: "0.2"
bundle_version: "1.1.0"
title: "PostgreSQL Production Architecture & Invariants"
description: "Curated architectural decisions for zero-downtime migrations, row-level security, connection pool hygiene, and scalable indexing."
license: "MIT"
verified: { by: human/sknr, at: "2026-10-20T10:22:35Z" }
---

# PostgreSQL Production Architecture Seed Bundle

Curated database invariants and operational decisions for zero-downtime DDL schema migrations, row-level security isolation, connection pooling resilience, and scalable query execution.

## Decisions
* [Zero-Downtime Schema Migrations](decisions/zero-downtime-migrations.md): Safe DDL patterns prohibiting blocking table locks on live tables.
* [Row-Level Security & Multi-Tenant Isolation](decisions/row-level-security.md): Mandatory RLS policies, tenant session context, and table owner isolation for multi-tenant data.
* [Connection Pool & PgBouncer Hygiene](decisions/connection-hygiene.md): Bounded connection counts, timeout discipline, and transaction pooling invariants.
* [Keyset Cursor-Based Pagination](decisions/cursor-pagination.md): Prohibition of deep OFFSET/LIMIT scans on high-cardinality tables.
