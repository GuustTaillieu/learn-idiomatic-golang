# Learning Record 0040: Production Diagnostics with pprof & Goroutine Leak Hunting

## Date: 2026-10-03

## Context & Review of Learner's Rate Limiter
- Learner implemented `IPRateLimiter` with Token Bucket algorithm and ticker-based cleanup loop.
- Learner returned `(rateLimiter *IPRateLimiter, cleanupFunc func())` from the constructor to safely stop the background ticker goroutine.
- Learner wrote comprehensive unit tests verifying burst limits and HTTP 429 status code handling.

## Architectural Findings & Nuances
1. **Mutex Ownership & Encapsulation**:
   - Mutexes should be embedded or declared by value (`sync.RWMutex`), rather than as pointers (`*sync.RWMutex`).
   - Fields guarded by a mutex should be unexported (e.g. `connections` instead of `Connections`) so consumers cannot bypass the lock and induce race conditions.
2. **Middleware Ordering for Full Observability**:
   - Placing `RequestIDMiddleware` and `LoggingMiddleware` before `RateLimitMiddleware` guarantees that rate-limited 429 requests are tracked with correlation IDs and logged in `slog`.
3. **Go Runtime Diagnostics (`net/http/pprof`)**:
   - Standard library diagnostic endpoints (`/debug/pprof/`) provide real-time CPU profiling, heap allocation inspection, and full goroutine stack traces.
   - Diagnosing goroutine leaks and lock contention in live services without external third-party APM overhead.

## Artifacts Created:
- Lesson: `lessons/0040-production-profiling-pprof-and-goroutine-leaks.html`
- Reference: `reference/0013-pprof-profiling-and-diagnostics.html`
