---
okf_version: "0.2"
bundle_version: "1.0.0"
title: "OAuth 2.1 & OpenID Connect Security Invariants"
description: "Curated architectural decisions for PKCE enforcement, exact redirect URI matching, and token entropy."
license: "MIT"
---

# OAuth 2.1 & OIDC Security Invariants Seed Bundle

Authoritative security invariants and protocol governance based on RFC 7636, RFC 8252, and the OAuth 2.1 draft specification.

## Decisions
* [Mandatory PKCE for All Clients](decisions/pkce-enforcement.md): Proof Key for Code Exchange (S256) requirement across public and confidential clients.
* [Strict Redirect URI Matching](decisions/redirect-uri-strict.md): Exact character-for-character matching and prohibition of wildcards or path traversal.
* [Cryptographic State & Nonce Validation](decisions/nonce-state-entropy.md): High-entropy CSRF and replay prevention invariants for authorization flows.
