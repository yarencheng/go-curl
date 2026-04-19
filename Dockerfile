# Stage 1: Build and Test
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Copy go.mod and go.sum (if present) for caching dependencies
COPY go.mod go.sum* ./
RUN go mod download

# Copy project files
COPY . .

# 1. Run all unit tests
RUN go test -v -cover ./...

# 2. Build web assembly output
RUN GOOS=wasip1 GOARCH=wasm go build -o /out/main.wasm ./cmd/go-curl/

# Stage 2: Wasmtime Runner
FROM debian:bullseye-slim

# Install Wasmtime
RUN apt-get update && \
    apt-get install -y --no-install-recommends curl xz-utils ca-certificates && \
    curl https://wasmtime.dev/install.sh -sSf | bash && \
    mv /root/.wasmtime/bin/wasmtime /usr/local/bin/wasmtime && \
    apt-get remove -y curl xz-utils && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /root/.wasmtime

# Non-Root Execution
RUN useradd -m -u 1001 appuser
USER appuser

WORKDIR /app

# Copy built artifacts
COPY --from=builder /out/main.wasm ./main.wasm

# 3. Run web assembly with Wasmtime
ENTRYPOINT ["wasmtime", "run", "main.wasm"]
