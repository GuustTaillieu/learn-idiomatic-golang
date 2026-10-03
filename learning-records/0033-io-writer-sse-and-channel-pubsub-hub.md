# Learning Record 0033: Real-Time Streaming, io.Writer & Channel Pub-Sub Hub

## Date: 2026-09-29

## Benchmark Results Verified
- `BenchmarkUnpooledBuffer-16`: 26.26 ns/op, 64 B/op, 1 allocs/op
- `BenchmarkPooledBuffer-16`: 15.24 ns/op, 0 B/op, 0 allocs/op
- Achieved zero heap allocations in the buffer hot path.

## Architectural Concepts Introduced
1. **Frontend-Backend Integration Spectrum**:
   - Evaluated gRPC (backend-to-backend RPC) vs WebSockets (bidirectional full duplex) vs Server-Sent Events (SSE).
   - Selected SSE as the optimal production pattern for real-time order tracking and browser dashboards.
2. **`io.Writer` and `http.Flusher`**:
   - `io.Writer` as Go's unified streaming abstraction.
   - `http.Flusher` interface type assertion to bypass HTTP response buffering and stream packets across TCP immediately.
3. **Channel Pub-Sub Hub Pattern**:
   - Managing dynamic subscriber channels (`map[chan T]struct{}`).
   - Non-blocking channel broadcast using `select ... default` to prevent slow network consumers from deadlocking backend workers.
   - Graceful lifecycle teardown via `r.Context().Done()`.

## Artifacts Created:
- Lesson: `lessons/0033-io-writer-sse-and-channel-pubsub-hub.html`
- Reference: `reference/0006-io-writer-and-streaming.html`
