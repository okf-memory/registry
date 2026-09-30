---
type: Decision
title: "Connection Pool & PgBouncer Hygiene"
tags: [postgres, database, pool, pgbouncer, performance]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/connection-hygiene"
---

# Connection Pool & PgBouncer Hygiene

Applications interacting with PostgreSQL MUST enforce strict connection bounds and compatibility with transaction-level pooling proxies.

## Invariants
1. **Bounded Client Pools:** Application services MUST configure strict maximum connection pool sizes calculated against the database instance's capacity (`max_connections`). Unbounded connection pools are prohibited.
2. **Transaction Pooling Compatibility:** Applications connecting through PgBouncer in transaction mode MUST NOT execute session-level state alterations (`SET SESSION`, `LISTEN/NOTIFY`, unpinned prepared statements, or temporary tables).
3. **Mandatory Statement Timeouts:** Application connection strings or session initialization hooks MUST declare a bounded `statement_timeout` (e.g. 15s to 30s) to prevent runaway unindexed queries from consuming backend worker processes indefinitely.
4. **Health Check Discipline:** Pool health checks MUST use lightweight, non-locking probes (`SELECT 1`) rather than querying system catalogs.

## Related
* [Zero-Downtime Schema Migrations](zero-downtime-migrations.md)
* [Keyset Cursor-Based Pagination](cursor-pagination.md)
