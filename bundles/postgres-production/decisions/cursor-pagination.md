---
type: Decision
title: Keyset Cursor-Based Pagination
tags: [postgres, database, queries, pagination, performance]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/cursor-pagination"
---

# Keyset Cursor-Based Pagination

Public APIs and data-access layers MUST utilize keyset (cursor-based) pagination on high-cardinality tables instead of offset pagination.

## Invariants
1. **Offset Pagination Prohibition on Large Tables:** Queries on tables with potential cardinality > 10,000 rows MUST NOT utilize `OFFSET N`. Offset pagination incurs an O(N) penalty as PostgreSQL must scan and discard N rows on the server before returning the target page.
2. **Deterministic Sort Keys:** Cursor pagination MUST order by a unique, indexed tuple (e.g. `(created_at, id)`). The client cursor MUST encode this composite key (e.g. opaque Base64).
3. **Index Alignment:** The table MUST have a composite index matching the keyset filter direction (e.g. `CREATE INDEX ON orders (created_at DESC, id DESC)`) allowing the planner to execute index range scans (`WHERE (created_at, id) < ($1, $2) LIMIT $3`).

## Related
* [Zero-Downtime Schema Migrations](zero-downtime-migrations.md)
* [Connection Pool & PgBouncer Hygiene](connection-hygiene.md)
