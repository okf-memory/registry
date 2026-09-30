---
type: Decision
title: "Token Storage & Transport Security"
description: Restricts refresh token storage strictly to HttpOnly Secure SameSite cookies and prohibits storing credentials in Web Storage.
generated: { by: agent/cli, at: "2026-09-30T11:43:22Z" }
governance: constraint
---

# Decision: Token Storage & Transport Security

## Context
Browser applications that store tokens in `window.localStorage` or `window.sessionStorage` are completely vulnerable to Cross-Site Scripting (XSS) attacks. Any injected JavaScript script (via third-party npm packages, CDN compromises, or unescaped user input) can exfiltrate credentials silently.

## Decision
1. **No Web Storage for Long-Lived Tokens:** Neither access tokens nor refresh tokens may be placed in `localStorage` or `sessionStorage`.
2. **HttpOnly Cookies for Refresh Tokens:** Refresh tokens MUST be transported exclusively via cookies configured with:
   * `HttpOnly`: Prevents client-side JavaScript access.
   * `Secure`: Requires HTTPS transport.
   * `SameSite=Strict` (or `SameSite=Lax` for necessary cross-site top-level navigations): Protects against CSRF attacks.
   * `Path=/api/auth`: Restricts cookie transmission to the dedicated token rotation endpoint.
3. **In-Memory Access Tokens:** Short-lived access tokens should reside exclusively in client runtime memory (e.g. state management or closure variables) and cleared upon page reload or tab closure.

```mermaid
flowchart LR
    subgraph Browser["Browser Client"]
        Mem["In-Memory State<br/>(Short-Lived Access Token)"]
        CookieJar["Secure Cookie Jar<br/>(HttpOnly Refresh Token)"]
    end

    subgraph Backend["API Server"]
        AuthRoute["/api/auth/refresh"]
        ResourceRoute["/api/v1/resources"]
    end

    Mem -->|Authorization: Bearer ...| ResourceRoute
    CookieJar -->|Automatic HTTPS Cookie| AuthRoute
```

## Related
* [Lifecycle & Refresh Token Rotation](lifecycle-and-rotation.md)
* [Claims Validation & Payload Hygiene](claims-hygiene.md)
