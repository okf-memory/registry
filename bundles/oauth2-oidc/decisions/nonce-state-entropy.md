---
type: Decision
title: "Cryptographic State & Nonce Validation"
tags: [oauth, oidc, security, csrf, replay]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/nonce-state-entropy"
---

# Cryptographic State & Nonce Validation

Clients initiating OAuth 2.0 authorization or OpenID Connect authentication flows MUST generate and verify high-entropy `state` and `nonce` parameters.

## Invariants
1. **Cryptographic Entropy Requirement:** The `state` and `nonce` parameters MUST be generated using a cryptographically secure pseudo-random number generator (CSPRNG) with at least 128 bits (16 bytes) of entropy, formatted as URL-safe Base64 or hex.
2. **Client-Side Verification:** The client application MUST bind the `state` parameter to the user's browser session (e.g. encrypted or signed HttpOnly session cookie). Upon receiving the callback, the client MUST reject the request if the returned `state` does not match the session-bound value.
3. **OpenID Connect Nonce Binding:** When requesting an OpenID Connect ID Token, the client MUST include a `nonce`. The client MUST verify that the `nonce` claim in the decoded ID token matches the exact nonce generated for that authorization request before accepting user identity.

## Related
* [Mandatory PKCE for All Clients](pkce-enforcement.md)
* [Strict Redirect URI Matching](redirect-uri-strict.md)
