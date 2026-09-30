# Decisions
* [Algorithm Enforcement & Insecure Scheme Rejection](algorithm-enforcement.md) - Mandates explicit algorithm allowlisting, asymmetric signature verification across services, and unconditional rejection of alg none.
* [Claims Validation & Payload Hygiene](claims-hygiene.md) - Enforces strict validation of standard claims (iss, aud, exp) and strictly prohibits sensitive PII or secrets in unencrypted JWT payloads.
* [Lifecycle & Refresh Token Rotation](lifecycle-and-rotation.md) - Establishes short access token lifespans (5-15m) and single-use Refresh Token Rotation with automatic family revocation on replay detection.
* [Token Storage & Transport Security](storage-and-transport.md) - Restricts refresh token storage strictly to HttpOnly Secure SameSite cookies and prohibits storing credentials in Web Storage.
