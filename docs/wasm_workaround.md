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
- `[-]` **Custom Interface Binding (`--interface`)**: Binding requests to a specific local interface or IP address is not permitted by the browser's `fetch` API. (Upstream: `third_party/curl/src/tool_getparam.c:1592`)
- `[-]` **Custom DNS Servers (`--dns-servers`)**: DNS resolution is handled by the browser/system and cannot be overridden at the application level in WASM. (Upstream: `third_party/curl/src/tool_getparam.c:1469`)
- `[-]` **UNIX Domain Sockets (`--unix-socket`)**: UNIX sockets require OS-level file system and networking integration unavailable in the WASM sandbox. (Upstream: `third_party/curl/src/tool_getparam.c:1907`)

### Traffic Routing
- `[-]` **Manual Proxy Configuration (`--proxy`)**: In a browser environment, proxy settings are managed by the user's browser or system configuration and are transparent to the `fetch` API. (Upstream: `third_party/curl/src/tool_getparam.c:2790`)

### Filesystem
- `[x]` **In-Memory VFS**: File uploads/downloads (`-T`, `-o`, `-O`) interact with an in-memory filesystem (Afero) rather than the host disk.

### SSL/TLS
- `[-]` **Custom CA Bundles**: Browser environments typically manage their own certificate trust; `-k` (insecure) may be simulated.

---
*Last Updated: 2026-04-19*
