# Learning Record 0049: Asynchronous Trace Propagation & W3C TraceContext

## Date: 2026-10-05

## Context & Progression
- Evaluated learner's implementation of OpenTelemetry instrumentation in `internal/telemetry/tracer.go`, `internal/processor/order_placing.go`, and `internal/http/order_rpc.go`.
- Addressed learner's specific question regarding how task contexts should be handled in asynchronous background queues.
- Introduced W3C `traceparent` context serialization to preserve causal trace linkage between synchronous RPC handlers and asynchronous worker goroutines.

## Architectural Concepts Introduced
1. **The Asynchronous Context Disconnect**:
   - HTTP request contexts are canceled upon client response completion. Passing a request context into an asynchronous worker queue causes immediate `context canceled` failures.
   - Using a decoupled `context.Background()` or worker lifecycle context severs the cancellation signal, but drops trace metadata.
2. **W3C TraceContext Propagation**:
   - Using `propagation.TraceContext{}` to serialize the 16-byte `TraceID` and 8-byte `ParentSpanID` into a compact string (`traceparent`).
   - Persisting `traceparent` with the order metadata across database outbox queues.
   - Extracting and resuming the trace inside worker processors so Jaeger/Datadog displays a unified end-to-end trace tree.

## Artifacts Created:
- Lesson: `lessons/0049-asynchronous-trace-propagation-and-span-links.html`
- Reference: `reference/0022-asynchronous-trace-propagation.html`
