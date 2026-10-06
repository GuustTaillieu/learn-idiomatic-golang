# Learning Record 0046: Type-Safe Distributed Contracts with Protocol Buffers & Connect-RPC

## Date: 2026-10-05

## Context & Progression
- Evaluated learner's database migration engine (`internal/sqlite/migrations.go`) and initial schema + seed SQL files.
- Resolved unit test in-memory database initialization by ensuring test helpers execute `sqlite.Migrate(db)`. All tests across all 9 packages pass with 0 data races.
- Addressed learner's goal: learning modern frontend and microservice communication patterns (gRPC vs Connect-RPC) with Go.

## Architectural Concepts Introduced
1. **Contract-First API Architecture**:
   - Defining APIs in Protocol Buffers (`.proto`) to eliminate serialization drift and generate type-safe server & client stubs.
   - Dual-protocol versatility: accepting binary Protobuf from microservices and JSON from browsers on the exact same endpoint.
2. **Connect-RPC vs Traditional gRPC**:
   - Traditional gRPC requires HTTP/2 trailers unsupported by browser `fetch()`, necessitating complex Envoy proxies.
   - Connect-RPC (`connectrpc.com/connect`) runs natively over standard HTTP/1.1 and HTTP/2 POST requests with zero proxies.
   - Native integration with standard library `http.Handler` and `http.ServeMux`.

## Artifacts Created:
- Lesson: `lessons/0046-protocol-buffers-and-connect-rpc.html`
- Reference: `reference/0019-protocol-buffers-and-connect-rpc.html`
