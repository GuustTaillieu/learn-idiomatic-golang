# Learning Record 0020: Startup Sweeper Review & Resilient Retry Architecture

## Date: 2026-09-20

## Milestone & Verification
- The learner implemented the boot-time recovery sweeper in `cmd/server/main.go` using `orderStore.GetPendingOrders(ctx)`.
- Verified clean implementation of `GetPendingOrders` in `internal/store/order_sqlite.go` with `defer rows.Close()`, `rows.Next()`, and `rows.Err()`.
- Verified the phantom rollback bug in `internal/processor/place_order.go` was resolved by ensuring `tx.Commit()` succeeds prior to returning the compensating rollback closure.
- Verified redundant scan boilerplate was stripped from `internal/store/inventory_sqlite.go`.
- Verified entire test suite passes with zero race conditions (`go test -count=1 -v -race ./...`).
- Created Lesson 0020 detailing transient vs. non-transient error classification, exponential backoff, and Dead-Letter Queue (DLQ) state transitions.

## Next Target
- Add `Retries` and `MaxRetries` to `Order`.
- Implement worker retry loop with exponential backoff before triggering saga rollback.
- Move permanently failing orders to `DEAD_LETTER`.
