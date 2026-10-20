---
type: Decision
title: Row-Level Security & Multi-Tenant Isolation
description: Mandates explicit Row Level Security (RLS) policies and tenant context isolation on all multi-tenant or user-scoped tables.
tags: [postgres, database, security, rls, multi-tenancy, isolation]
generated: { by: agent/cli, at: "2026-10-20T13:46:00Z" }
verified: { by: human/sknr, at: "2026-10-20T15:17:15Z" }
status: stable
governance: constraint
id: "decisions/row-level-security"
---

# Row-Level Security & Multi-Tenant Isolation

Production database tables containing tenant-scoped, user-private, or organizational data MUST enforce row-level access control at the database engine level via PostgreSQL Row Level Security (RLS).

## Invariants

1. **Mandatory RLS Activation:** Every table storing multi-tenant or customer-scoped records MUST have Row Level Security explicitly enabled:
   ```sql
   ALTER TABLE tenant_records ENABLE ROW LEVEL SECURITY;
   ```
2. **Force RLS for Table Owners (`FORCE ROW LEVEL SECURITY`):** Application database roles that own tables bypass RLS by default unless forced. Multi-tenant tables MUST enforce RLS even for the table owner role to prevent accidental developer or service bypass:
   ```sql
   ALTER TABLE tenant_records FORCE ROW LEVEL SECURITY;
   ```
3. **Session Variable or JWT Claim Context:** RLS policies MUST filter against scoped session variables set by the connection pool or application layer (e.g. `current_setting('app.current_tenant_id', true)::uuid`) or PostgreSQL user roles, rather than relying on application-level `WHERE tenant_id = ...` clauses alone:
   ```sql
   CREATE POLICY tenant_isolation_policy ON tenant_records
     AS RESTRICTIVE
     USING (tenant_id = current_setting('app.current_tenant_id', true)::uuid)
     WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::uuid);
   ```
4. **Index Coverage for Policy Filters:** Any column referenced in an RLS policy condition (specifically `tenant_id`, `org_id`, or `owner_id`) MUST be indexed as the leading column in composite indexes to prevent full-table sequential scans on filtered queries.
5. **Superuser Prohibition in Application Pools:** Application connection pools (e.g. PgBouncer, backend microservices) MUST connect using unprivileged database roles (`NOSUPERUSER`, `NOBYPASSRLS`). Superusers bypass RLS unconditionally.

## Related
* [Zero-Downtime Schema Migrations](zero-downtime-migrations.md)
* [Connection Pool & PgBouncer Hygiene](connection-hygiene.md)
* [Keyset Cursor-Based Pagination](cursor-pagination.md)
