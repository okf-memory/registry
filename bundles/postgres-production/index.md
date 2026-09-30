---
okf_version: "0.2"
bundle_version: "1.0.0"
title: "PostgreSQL Production Architecture & Invariants"
description: "Curated architectural decisions for zero-downtime migrations, connection pool hygiene, and scalable indexing."
license: "MIT"
---

# PostgreSQL Production Architecture Seed Bundle

Curated database invariants and operational decisions for zero-downtime DDL schema migrations, connection pooling resilience, and scalable query execution.

## Decisions
* [Zero-Downtime Schema Migrations](decisions/zero-downtime-migrations.md): Safe DDL patterns prohibiting blocking table locks on live tables.
* [Connection Pool & PgBouncer Hygiene](decisions/connection-hygiene.md): Bounded connection counts, timeout discipline, and transaction pooling invariants.
* [Keyset Cursor-Based Pagination](decisions/cursor-pagination.md): Prohibition of deep OFFSET/LIMIT scans on high-cardinality tables.
