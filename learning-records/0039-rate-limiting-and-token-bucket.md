# Learning Record 0039: Rate Limiting, Backpressure & the Token Bucket Pattern

## Date: 2026-10-02

## Milestone & Verified Execution
- All tests across all 8 packages pass with 0 data races.
- Learner implemented `ResponseRecorder` with `Flush()` passthrough and `LoggingMiddleware`.
- Order of middleware execution verified: `RequestIDMiddleware` &rarr; `LoggingMiddleware` ensures `request_id` is present on every access log line.

## Architectural Concepts Introduced
1. **The Principle of Backpressure**:
   - Preventing cascading failures (thread exhaustion, database saturation, tail latency spikes) by shedding load early with HTTP 429.
2. **The Token Bucket Algorithm (`golang.org/x/time/rate`)**:
   - Defining rate limits as a steady-state token inflow (`rate.Limit`) and an instantaneous burst capacity (`burst`).
   - Using non-blocking `limiter.Allow()` for HTTP gateway decisions.
3. **Memory Safety in Stateful Middleware**:
   - Mitigating unbounded map growth for per-IP rate limiters using a sliding TTL window (`lastSeen time.Time`) and periodic ticker-based eviction.

## Artifacts Created:
- Lesson: `lessons/0039-rate-limiting-and-token-bucket.html`
- Reference: `reference/0012-rate-limiting-and-token-bucket.html`
