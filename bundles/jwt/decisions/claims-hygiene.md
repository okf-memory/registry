---
type: Decision
title: "Claims Validation & Payload Hygiene"
description: "Enforces strict validation of standard claims (iss, aud, exp) and strictly prohibits sensitive PII or secrets in unencrypted JWT payloads."
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
governance: constraint
---

# Decision: Claims Validation & Payload Hygiene

## Context
Developers frequently confuse signed tokens with encrypted tokens. Standard JSON Web Tokens (JWS) are only signed and base64url-encoded; anyone inspecting the network stream or browser memory can read the payload contents in plain text. Furthermore, omitting audience (`aud`) or issuer (`iss`) checks allows cross-service token substitution attacks.

## Decision
1. **Mandatory Standard Claims:** Every issued token MUST contain:
   * `iss` (Issuer): Fully qualified URI of the authentication service.
   * `aud` (Audience): Identifier of the specific intended resource server.
   * `sub` (Subject): Stable user or service UUID.
   * `exp` (Expiration): Unix timestamp of expiration.
   * `iat` (Issued At): Unix timestamp of issuance.
2. **Strict Audience Checking:** Resource servers MUST reject tokens where `aud` does not match their own service identifier.
3. **Zero Sensitive Data in Claims:** Passwords, API keys, credit card numbers, and Personally Identifiable Information (PII) like national IDs or phone numbers MUST NEVER be stored in JWT claims.

```mermaid
flowchart TD
    RawToken[Incoming Token] --> Decode[Decode Payload]
    Decode --> IssCheck{iss == Expected Issuer?}
    IssCheck -->|No| Reject1[Reject: Invalid Issuer]
    IssCheck -->|Yes| AudCheck{aud == My Service ID?}
    AudCheck -->|No| Reject2[Reject: Audience Mismatch]
    AudCheck -->|Yes| ExpCheck{exp > Current Time?}
    ExpCheck -->|No| Reject3[Reject: Token Expired]
    ExpCheck -->|Yes| Accept[Accept Context]
```

## Related
* [Algorithm Enforcement & Insecure Scheme Rejection](algorithm-enforcement.md)
* [Token Storage & Transport Security](storage-and-transport.md)
