# Functional Gap Map: curl

This document tracks known functional gaps, intentional deviations, and implemented workarounds in the `go-curl` implementation.

## Status Definitions

- `[x]` **Workaround Applied**: The feature is addressed using a simulator-specific solution rather than a 1:1 port.
- `[-]` **Deliberately Unsupported**: The feature is explicitly excluded (e.g., hardware-dependent, security-risk, or WASM-incompatible).
- `[ ]` **Unresolved / Decision Pending**: A gap has been identified, but the implementation strategy or priority remains to be decided.

---

## Gap Repository

### Network Sandbox
- `[x]` **Localhost Redirection**: Network requests may be redirected to mock services in the WASM environment.
- `[x]` **CORS Limitations**: Browser-based execution is subject to CORS restrictions unless proxied.

### Filesystem
- `[x]` **In-Memory VFS**: File uploads/downloads (`-T`, `-o`, `-O`) interact with an in-memory filesystem (Afero) rather than the host disk.

### SSL/TLS
- `[-]` **Custom CA Bundles**: Browser environments typically manage their own certificate trust; `-k` (insecure) may be simulated.

---
*Last Updated: 2026-04-19*
