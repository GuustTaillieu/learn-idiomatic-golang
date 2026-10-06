# Learning Record 0048: OpenTelemetry Distributed Tracing & Span Lifecycle

## Date: 2026-10-05

## Context & Progression
- Diagnosed Go 1.22+ `http.ServeMux` pattern conflict panic when mounting `GET /` alongside `/order.v1.OrderService/`.
- Explained routing precedence: `GET /` specifies a method but has a broader path, whereas `/order.v1.OrderService/` has no method but a narrower path. Resolved by registering root catch-all `/` without method prefix.
- Verified learner's end-to-end integration tests for Connect-RPC in `internal/http/order_rpc_test.go`.
- Advanced to enterprise observability: instrumenting asynchronous distributed pipelines with OpenTelemetry.

## Architectural Concepts Introduced
1. **Go 1.22+ ServeMux Conflict Resolution**:
   - Strict pattern specificity rules: A pattern cannot be simultaneously more specific in method and less specific in path than another pattern.
   - Using root `/` catch-all for single-binary asset servers while preserving namespace prefixes (`/order.v1.OrderService/`).
2. **OpenTelemetry Distributed Tracing**:
   - Tracing request lifecycles across asynchronous boundaries (HTTP -> SQLite Outbox -> Dispatcher -> Worker Pool -> Sagas).
   - Linking child spans to parent spans via `context.Context`.
   - Recording errors and metadata attributes without leaking sensitive data.

## Artifacts Created:
- Lesson: `lessons/0048-opentelemetry-distributed-tracing-and-spans.html`
- Reference: `reference/0021-opentelemetry-distributed-tracing.html`
