---
type: Decision
title: "Lifecycle & Refresh Token Rotation"
description: Establishes short access token lifespans (5-15m) and single-use Refresh Token Rotation with automatic family revocation on replay detection.
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
governance: constraint
---

# Decision: Token Lifecycle & Refresh Token Rotation (RTR)

## Context
Stateless JWT access tokens cannot be revoked individually without maintaining server-side blacklists, which destroys the horizontal scalability advantage of stateless authentication. Long-lived access tokens leave large windows of vulnerability if compromised.

## Decision
1. **Short-Lived Access Tokens:** Access tokens MUST have an expiration time (`exp`) between 5 and 15 minutes.
2. **Refresh Token Rotation (RTR):** Every call to `/api/auth/refresh` MUST issue a new access token AND a new refresh token while immediately invalidating the previous refresh token.
3. **Replay Detection & Family Invalidation:** If an already-used refresh token is submitted again, the server MUST treat this as token theft, invalidate the entire token family (all issued refresh tokens for that session), and force immediate re-authentication.

```mermaid
flowchart TD
    Client[Client Request] --> AuthCheck{Refresh Token Valid?}
    AuthCheck -->|Valid & Unused| Issue[Issue New Access & Refresh Token]
    Issue --> InvalidateOld[Invalidate Old Refresh Token]
    AuthCheck -->|Already Used Replay| RevokeAll[Alert: Revoke All Tokens in Family]
    RevokeAll --> ForceLogin[Force User Re-Authentication]
```

## Related
* [Token Storage & Transport Security](storage-and-transport.md)
* [Claims Validation & Payload Hygiene](claims-hygiene.md)
