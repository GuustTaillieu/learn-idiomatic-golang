# Learning Record 0021: Exponential Backoff Review & The Transactional Outbox

## Date: 2026-09-20

## Milestone & Verification
- The learner implemented exponential retry backoff and Dead-Letter Queue handling in `internal/domain/queue.go`.
- Added `Retries` and `MaxRetries` to `Order` struct.
- Created `IsRetryable` error check with `ErrTransient` sentinel error in `internal/domain/errors.go`.
- Verified test suite passes under race detector (`go test -count=1 -v -race ./...`).
- Tested two scenarios:
  1. Temporary failures recovering to `StatusCompleted`.
  2. Persistent failures exhausting retries and ending in `StatusDeadLetter`.

## Feedback Provided
1. **Loop vs Recursion**: Iterative loops (`for`) in Go are preferred over recursive function calls for retry state machines to avoid growing stack frames.
2. **Configurable Base Delay**: Moving `baseDelay` from a package global variable to a struct field on `Queue` allows tests to set a fast delay (e.g. 5ms) instead of waiting seconds.
3. **Dead Letter Cleanup**: When moving to `DEAD_LETTER`, ensuring compensation closures are evaluated if partial work occurred.

## Next Stage
- Implement the continuous Transactional Outbox Dispatcher.
- Decouple HTTP ingestion from in-memory channel queuing.
