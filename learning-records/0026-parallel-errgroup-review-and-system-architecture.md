# Learning Record 0026: Parallel errgroup Review & End-to-End Architecture

## Date: 2026-09-27

## Milestone & Verified Execution
- The learner implemented `Parallel` processor using `golang.org/x/sync/errgroup` in `internal/processor/parallel.go`.
- Wrote two comprehensive unit tests in `internal/processor/parallel_test.go`:
  1. Concurrency speed verification: two tasks (50ms and 60ms) completed in 60ms total.
  2. First-error cancellation: a 10ms failure immediately cancelled a waiting sibling goroutine via context.
- Verified thread-safe rollback slice assignment with `sync.Mutex`.
- Verified entire test suite passes across all packages with 0 data races.

## Architectural Clarifications Addressed
1. **Sync vs. Async Boundaries in Production**:
   - *Synchronous*: Cart validation, inventory reservation, and payment intent creation (returns Stripe `client_secret` or redirect URL immediately).
   - *Asynchronous*: Post-payment fulfillment (invoices, warehouse dispatch, email notifications) handled by outbox queue.
2. **Client Status Communication**:
   - Short polling (`GET /orders/{id}`) as the industry gold standard.
   - Server-Sent Events (SSE) and Optimistic UI notifications.
3. **Future-Proofing with Generics**:
   - Moving from concrete `Queue` to generic `Queue[T any]` for reusable asynchronous processing of arbitrary types (`Queue[*Order]`, `Queue[*EmailJob]`).
