# Learning Record 0037: The Ghost Rollback Trap & context.WithoutCancel

## Date: 2026-10-02

## Context & User Innovations
- Learner successfully resolved the 5-second graceful shutdown delay using `http.Server.BaseContext`.
- Learner upgraded `domain.Processor[T]` so that compensating rollback functions accept a context (`func(context.Context) error`), recognizing that database rollbacks require a context for cancellations and deadlines.

## Architectural Findings & Solutions
1. **The Ghost Rollback Trap in Distributed Sagas**:
   - When a processing pipeline times out (`context.DeadlineExceeded`), the incoming `ctx` is already cancelled (`ctx.Err() != nil`).
   - If that cancelled context is passed directly to the compensating rollback closure, all database operations within the rollback (e.g. `ReleaseStock`) immediately abort without running.
2. **Go 1.21 `context.WithoutCancel`**:
   - Strips the cancellation signal from a parent context while preserving all values (Request IDs, correlation IDs, `slog` metadata).
   - Allows spawning a clean, bounded child timeout context:
     `rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)`
3. **Mock Signature Alignment**:
   - Coordinated test mock updates across `internal/domain/outbox_test.go`, `internal/queue/order_queue_test.go`, and `internal/queue/queue_test.go` to match the updated `cleanup func(context.Context) error` contract.

## Artifacts Created:
- Lesson: `lessons/0037-the-ghost-rollback-trap-and-context-without-cancel.html`
- Reference: `reference/0010-context-without-cancel-and-compensating-transactions.html`
