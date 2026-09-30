---
type: Decision
title: Strict Redirect URI Matching
tags: [oauth, security, redirect-uri, validation]
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
status: stable
governance: constraint
id: "decisions/redirect-uri-strict"
---

# Strict Redirect URI Matching

Authorization servers MUST validate client redirect URIs using exact string matching against pre-registered endpoints.

## Invariants
1. **Exact String Match Required:** The authorization server MUST perform exact, full-string comparisons between the requested `redirect_uri` and the client's pre-registered URIs.
2. **Wildcards & Subdomains Prohibited:** Wildcards (`*`), regular expressions, partial path matches, and open subdomains in registered redirect URIs are strictly prohibited.
3. **Localhost & Loopback Security:** Native apps redirecting to loopback interfaces MUST use explicit loopback IP literals (`http://127.0.0.1:<port>/callback` or `http://[::1]:<port>/callback`) rather than `http://localhost`. The port number may be dynamic, but the host and path MUST remain fixed.
4. **Scheme Security:** Production redirect URIs MUST enforce the `https` scheme, except for private-use custom URI schemes in native mobile applications.

## Related
* [Mandatory PKCE for All Clients](pkce-enforcement.md)
* [Cryptographic State & Nonce Validation](nonce-state-entropy.md)
