# 🐹 GopherOrder — Concurrent Order Processing Engine

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Tests](https://img.shields.io/badge/Tests-100%25%20Passing%20with%20--race-success?style=flat)](file:///home/gustavo/Documents/projects/idiomatic-go)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Container](https://img.shields.io/badge/Docker-Distroless%20%3C25MB-2496ED?style=flat&logo=docker)](file:///home/gustavo/Documents/projects/idiomatic-go/Dockerfile)

> **⚠️ Educational & Reference Project**  
> This repository is a dedicated **hands-on learning lab and reference project** built to refresh, sharpen, and document senior-level Golang engineering patterns. It is designed as an architectural case study exploring concurrent systems, distributed design patterns, and idiomatic Go practices rather than a commercial production deployment.

---

## 📌 Project Overview

**GopherOrder** is an event-driven, concurrent order processing and inventory fulfillment system. Built with a **Standard Library First** philosophy, it demonstrates how to build resilient, fault-tolerant backend services without relying on bloated third-party monolith frameworks.

The project models an e-commerce order lifecycle: from dual-protocol ingestion (REST & Connect-RPC), through SQLite transactional outbox queues, to asynchronous worker pools executing coordinated saga rollbacks, and streaming real-time status updates back to a browser dashboard via Server-Sent Events (SSE).

---

## 🏗️ Architectural Highlights

```
                                      ┌───────────────────────────────────┐
                                      │   Clients (Web / RPC / CLI)       │
                                      └─────────────────┬─────────────────┘
                                                        │
                      ┌─────────────────────────────────┴─────────────────────────────────┐
                      │                                                                   │
               HTTP / REST & SSE                                                     Connect-RPC
         (Go 1.22+ http.ServeMux)                                             (Protocol Buffers / HTTP)
                      │                                                                   │
                      └─────────────────────────────────┬─────────────────────────────────┘
                                                        ▼
                                       ┌────────────────────────────────┐
                                       │    HTTP Middleware Pipeline    │
                                       │  - Token Bucket Rate Limiting  │
                                       │  - X-Request-ID Correlation    │
                                       │  - Contextual slog Access Log  │
                                       └────────────────┬───────────────┘
                                                        │
                                                        ▼
                                       ┌────────────────────────────────┐
                                       │   SQLite Transactional Outbox  │
                                       │   (Atomic Order Persistence)   │
                                       └────────────────┬───────────────┘
                                                        │ (Polled via Claim-Check)
                                                        ▼
                                       ┌────────────────────────────────┐
                                       │     Asynchronous Queue         │
                                       │    (Worker Pool & Generics)    │
                                       └────────────────┬───────────────┘
                                                        │
                                                        ▼
                                       ┌────────────────────────────────┐
                                       │    Composite Saga Processors   │
                                       │ - OrderPlacing (Deduct Stock)  │
                                       │ - Paying (Process Payment)     │
                                       │ - Coordinated Rollbacks        │
                                       └────────────────┬───────────────┘
                                                        │ (Broadcast State)
                                                        ▼
                                       ┌────────────────────────────────┐
                                       │      Generic Pub/Sub Hub       │
                                       │   (Non-blocking SSE Stream)    │
                                       └────────────────────────────────┘
```

### Key Engineering Patterns Implemented:

1. **Idiomatic Concurrency & Primitives**:
   - Generic bounded worker pools with non-blocking channels and graceful drain semantics.
   - Zero-allocation buffer pooling using `sync.Pool` benchmarked to **0 B/op and 0 allocs/op**.
   - Request deduplication using generic `Singleflight[T]` to prevent cache-stampede / thundering-herd issues.
   - Token Bucket rate limiting (`golang.org/x/time/rate`) with background memory cleanup.

2. **Package Structure & Domain Modeling**:
   - Strictly organized using Ben Johnson's **Standard Package Layout**: pure domain interfaces in `internal/domain`, database logic in `internal/sqlite`, and business pipelines in `internal/processor`.
   - Strongly-typed domain IDs implementing `database/sql.Scanner`, `driver.Valuer`, and JSON/Text marshalers.
   - Separation of DDL and DML: pure Go SQLite migrations and idempotent seeders embedded via Go 1.16+ `//go:embed`.

3. **Fault Tolerance & Distributed Sagas**:
   - **Transactional Outbox Pattern**: Decouples API ingestion from asynchronous execution with a stateful Claim-Check pattern (`PENDING` → `RUNNING` → `COMPLETED`/`FAILED`).
   - **Compensating Rollbacks**: Composite processors return compensating rollback closures to cleanly revert stock reservations when subsequent saga steps fail.
   - The **"Ghost Rollback Trap"** resolved via Go 1.21's `context.WithoutCancel` to ensure database rollbacks succeed even when client request contexts expire.

4. **Dual Protocols: REST + Connect-RPC**:
   - Modern Go 1.22+ `http.ServeMux` routing REST and Server-Sent Events (SSE).
   - Single-binary static web dashboard serving via `//go:embed`.
   - **Connect-RPC**: Schema-driven Protocol Buffers (`proto/order/v1/order.proto`) mounted directly onto `net/http` without needing an Envoy proxy.

5. **Production Observability**:
   - Structured logging with `log/slog` enriched with request correlation IDs (`X-Request-ID`).
   - Live runtime diagnostics and memory profiling via `net/http/pprof`.
   - **OpenTelemetry Distributed Tracing**: W3C `traceparent` context serialization across asynchronous outbox database boundaries into background worker pools.

6. **Cloud-Native Containerization**:
   - Multi-stage Docker build producing a statically linked pure Go binary (`CGO_ENABLED=0`).
   - Base image: `gcr.io/distroless/static-debian12:nonroot` running under unprivileged UID 65532 with a total footprint of **<25 MB** and zero CVEs.

---

## 📁 Repository Structure

```text
├── cmd/
│   ├── server/           # Main application composition root
│   └── stress/           # High-concurrency CLI load test generator
├── internal/             # Private server implementations (Go-enforced)
│   ├── config/           # 12-factor configuration with fail-fast startup
│   ├── domain/           # Entities (Order, Item, Stock), IDs, interfaces
│   ├── event/            # Generic Pub/Sub event broker
│   ├── health/           # Composite liveness and readiness probes
│   ├── http/             # HTTP handlers, SSE streaming, and Connect-RPC
│   │   └── middleware/   # Rate limiting, logging, and request correlation
│   ├── lib/              # Buffer pool (sync.Pool), Singleflight, DBTX
│   ├── processor/        # OrderPlacing, Paying, and Parallel saga pipelines
│   ├── queue/            # Generic worker queue engine
│   ├── sqlite/           # Pure Go SQLite storage & embedded SQL migrations
│   └── telemetry/        # OpenTelemetry tracing & W3C context propagation
├── proto/                # Public Protocol Buffer contract (API schema)
├── web/                  # Embedded HTML/CSS real-time dashboard
├── lessons/              # 50 step-by-step interactive lesson guides (HTML)
├── reference/            # 25 architectural reference cheat sheets (HTML)
├── learning-records/     # Chronological learning logs & design decisions (Markdown)
├── Dockerfile            # Multi-stage Google Distroless container
└── .dockerignore         # Docker build context exclusions
```

---

## 🚀 Getting Started

### Prerequisites
- **Go**: 1.24+ (or Go 1.22+)
- Optional: **Docker** / **Podman** for containerized execution.

### 1. Run All Tests
Verify the complete test suite across all 9 packages with Go's race detector:
```bash
go test -v -race ./...
```

### 2. Run Locally
Start the server with embedded database migrations and seeding:
```bash
go run cmd/server/main.go
```
The server will boot on `http://localhost:8080`.

### 3. Open the Dashboard
Navigate to [http://localhost:8080](http://localhost:8080) in your browser.  
- The Server-Sent Events (SSE) status indicator will glow green.
- Click **"⚡ Dispatch to Outbox"** to trigger an order and watch its real-time lifecycle (`PENDING` → `PROCESSING` → `COMPLETED`).

### 4. Run the High-Concurrency Load Tester
Simulate 100 concurrent requests across 15 worker goroutines:
```bash
go run cmd/stress/main.go -c 15 -n 100
```
Observe how the Token Bucket rate limiter protects the backend with `429 Too Many Requests` while the web UI continues streaming events smoothly.

### 5. Inspect Live Diagnostics
Visit the built-in profiling endpoints:
- pprof Index: [http://localhost:8080/debug/pprof/](http://localhost:8080/debug/pprof/)
- Health Check: [http://localhost:8080/healthz](http://localhost:8080/healthz)
- Readiness Probe: [http://localhost:8080/readyz](http://localhost:8080/readyz)

### 6. Build Distroless Docker Image
Build and inspect the minimal container footprint:
```bash
docker build -t gopher-order:latest .
docker images gopher-order:latest  # Total size: ~20-25 MB!

# Run container
docker run --rm -p 8080:8080 gopher-order:latest
```

---

## 🌐 Multi-Language Client Support (Connect-RPC)

Because the API contract is defined in `proto/order/v1/order.proto`, clients can interact with this service in any language:

- **TypeScript (React / Next.js)**: Connect-RPC client using `@connectrpc/connect-web`.
- **Kotlin (Android / Jetpack Compose)**: Mobile client using `com.connectrpc:connect-kotlin`.
- **Rust**: HTTP/JSON or binary protobuf using `reqwest` or `tonic`.
- **cURL / HTTP Clients**: Native JSON POST requests to `/order.v1.OrderService/CreateOrder`.

*(See detailed code snippets in [`reference/0025-multi-language-connect-rpc-clients.html`](reference/0025-multi-language-connect-rpc-clients.html)).*

---

## 📚 Learning Documentation

This repository contains full chronological documentation of the entire curriculum:
- **`lessons/`**: 50 standalone HTML lesson guides with code walkthroughs.
- **`reference/`**: 25 architectural reference sheets covering concurrency, sagas, zero-allocation pooling, OpenTelemetry, and Connect-RPC.
- **`learning-records/`**: 52 detailed Markdown logs capturing every design decision, trade-off, and refactoring milestone.

---

## 📄 License
This project is open source and available under the [MIT License](LICENSE).
