---
okf_version: "0.2"
bundle_version: "1.0.0"
title: "JWT Best Practices & Security Seed Bundle"
description: "Curated architectural decisions, token lifecycle invariants, and security patterns for JSON Web Tokens (JWT)."
license: "MIT"
---

# JSON Web Token (JWT) Security Knowledge Bundle

Curated architectural decisions, token lifecycle governance, and security patterns for production JWT implementations.

## Decisions
* [Algorithm Enforcement & None Rejection](decisions/algorithm-enforcement.md): Mandates explicit algorithm allowlisting, asymmetric signature verification across services, and unconditional rejection of alg none.
* [Token Storage & Transport Security](decisions/storage-and-transport.md): Restricts refresh token storage strictly to HttpOnly Secure SameSite cookies and prohibits storing credentials in Web Storage.
* [Lifecycle & Refresh Token Rotation](decisions/lifecycle-and-rotation.md): Establishes short access token lifespans (5-15m) and single-use Refresh Token Rotation with automatic family revocation on replay detection.
* [Claims Validation & Payload Hygiene](decisions/claims-hygiene.md): Enforces strict validation of standard claims (iss, aud, exp) and strictly prohibits sensitive PII or secrets in unencrypted JWT payloads.
