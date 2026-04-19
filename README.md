# go-curl

**go-curl** is a Go implementation of the `curl` utility for WebAssembly and CLI environments. It provides a lightweight, sandboxed version of curl targeting high functional parity with the original GNU/curl tool.

## Key Features

- **High Parity**: Aims to support core curl flags and behaviors.
- **WASM First**: Optimized for `js/wasm` and `wasip1/wasm` targets.
- **In-Memory VFS**: Isolated filesystem operations for secure execution.
- **Structured Logging**: Deep observability with `zerolog`.

## Project Structure

- `cmd/go-curl/`: Main entry point for native and WASM builds.
- `cmd/go-bash-wasm/`: Specific WASM entry point for integration.
- `internal/app/`: Core application logic and dependency injection.
- `internal/commands/curl/`: `curl` command implementation.
- `docs/`: Parity tracking and functional gap documentation.

## Prerequisites

- **Go 1.26+**
- **Docker** (for containerized builds)

## Local Development

### 1. Build and Run
```bash
go run ./cmd/go-curl/ https://example.com
```

### 2. Run Tests
```bash
go test -v ./...
```

## Docker Deployment
```bash
docker build -t go-curl .
docker run --rm go-curl https://example.com
```

---
*Developed by the go-curl team.*
