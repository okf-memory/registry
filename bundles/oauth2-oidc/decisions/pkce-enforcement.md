---
type: Decision
title: Mandatory PKCE for All Clients
tags: [oauth, security, pkce, auth, protocol]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/pkce-enforcement"
---

# Mandatory PKCE for All Clients

Authorization Code grants MUST require Proof Key for Code Exchange (PKCE, RFC 7636) across all client architectures.

## Invariants
1. **Universal PKCE Requirement:** Both public clients (SPAs, mobile apps) and confidential server-side web applications MUST send `code_challenge` and `code_challenge_method=S256` during authorization requests.
2. **Prohibition of Plain Transformation:** Authorization servers MUST reject requests with `code_challenge_method=plain`. The transform method MUST be `S256` (SHA-256 base64url-encoded).
3. **Implicit & Password Grants Deprecated:** The OAuth 2.0 Implicit Grant (`response_type=token`) and Resource Owner Password Credentials Grant MUST NOT be implemented or enabled.
4. **Code Verifier Single-Use:** The authorization server MUST ensure that each authorization code and associated `code_verifier` can be redeemed exactly once. Any replay attempt MUST immediately revoke all tokens issued from that code.

## Related
* [Strict Redirect URI Matching](redirect-uri-strict.md)
* [Cryptographic State & Nonce Validation](nonce-state-entropy.md)
