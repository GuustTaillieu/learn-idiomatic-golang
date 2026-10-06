# Learning Record 0047: Type-Safe Connect-RPC Client & In-Memory Integration Testing

## Date: 2026-10-05

## Context & Progression
- Evaluated learner's Protocol Buffer schema (`proto/order/v1/order.proto`), code generation (`protoc-gen-go`, `protoc-gen-connect-go`), and RPC implementation (`internal/http/order_rpc.go`).
- Verified Connect-RPC handler registration on `http.ServeMux` alongside existing REST, SSE, and pprof endpoints.
- Advanced to client-side RPC consumption and in-memory black-box integration testing using `httptest.Server`.

## Architectural Concepts Introduced
1. **Type-Safe RPC Client Stubs**:
   - Using generated `orderv1connect.NewOrderServiceClient` to eliminate manual URL building, HTTP headers, and JSON serialization.
   - Context propagation and automatic error translation into Connect status codes (`connect.Code`).
2. **In-Memory Black-Box Integration Testing**:
   - Bootstrapping lightweight `httptest.NewServer` instances with in-memory stores to test the full HTTP/RPC networking pipeline.
   - Exercising `CreateOrder` and `GetOrder` RPC flows without external network dependencies.
3. **Prefix Routing in Go 1.22+**:
   - Verifying how `http.ServeMux` routes exact paths (`POST /orders`), prefix subtrees (`/order.v1.OrderService/*`), and root fallbacks (`/` for `web.FileServer()`).

## Artifacts Created:
- Lesson: `lessons/0047-type-safe-connect-rpc-client-and-integration-testing.html`
- Reference: `reference/0020-connect-rpc-client-and-integration-testing.html`
