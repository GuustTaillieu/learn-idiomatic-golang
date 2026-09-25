# Learning Record 0024: Context-Scoped Logging & Request Tracing

## Date: 2026-09-23

## Milestone & Conceptual Evaluation
- The learner implemented composite health checking using `MultiHealthChecker` and `DBHealthChecker` with context timeouts (`2*time.Second`).
- Verified test suite passes cleanly with zero data races.
- Evaluated the 3 architectural approaches to logging in Go:
  1. *Struct injection*: Clean for component identity, but rigid and cannot capture dynamic request IDs.
  2. *Package global `slog.Default()`*: High duplication of attributes on every line.
  3. *Context-scoped logger (`lib.ContextWithLogger` & `lib.Logger(ctx)`)*: Attaches request-scoped fields (`request_id`, `order_id`) once at boundaries; downstream functions inherit context without modifying struct signatures.
- Created Lesson 0024 covering request-scoped logging and `X-Request-ID` correlation middleware.
