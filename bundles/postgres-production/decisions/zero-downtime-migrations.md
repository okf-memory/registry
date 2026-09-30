---
type: Decision
title: Zero-Downtime Schema Migrations
tags: [postgres, database, ddl, migrations, reliability]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/zero-downtime-migrations"
---

# Zero-Downtime Schema Migrations

Database DDL operations MUST NOT hold exclusive table locks (`AccessExclusiveLock`) on active production tables.

## Invariants
1. **Index Creation Governance:** All new indexes on production tables MUST be created using `CREATE INDEX CONCURRENTLY`. Standard non-concurrent indexing is strictly prohibited as it blocks concurrent writes.
2. **Column Addition Rules:** Adding a non-nullable column MUST be executed in two steps: first add the column as nullable or with a runtime default (`ALTER TABLE ... ADD COLUMN ... DEFAULT ...`), backfill historical data out-of-band in batches, and finally add the `NOT NULL` constraint with a verified check.
3. **Lock Timeout Mandatory:** Every migration script or transaction MUST set an explicit, aggressive lock timeout (e.g. `SET lock_timeout = '2s';`). If the required lock cannot be acquired within 2 seconds, the transaction MUST fail and retry rather than queue behind long-running queries and block the entire connection pool.
4. **Column & Table Deletion (Expand/Contract):** Dropping columns or tables MUST follow the expand/contract pattern across at least two consecutive application releases (deprecate and ignore in code first, drop DDL second).

## Related
* [Connection Pool & PgBouncer Hygiene](connection-hygiene.md)
* [Keyset Cursor-Based Pagination](cursor-pagination.md)
