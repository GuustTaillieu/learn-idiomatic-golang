# Learning Record 0028: Generic Queue Review & Singleflight Request Deduplication

## Date: 2026-09-28

## Milestone & Verified Execution
- All tests across all 7 packages pass with 0 data races.
- The learner implemented `Queue[T any]` in `internal/lib/queue.go` with Functional Options (`QueueOptionFn[T]`).
- Identified a subtle concurrency gotcha in `lib.Queue[T]`: `q.retries` was declared on the struct, causing state sharing between concurrent worker goroutines.
- Introduced `golang.org/x/sync/singleflight` to solve the "Thundering Herd" problem in polling microservices.

## Architectural Concepts Introduced
1. **Singleflight Deduplication (`singleflight.Group`)**:
   - Deduplicates identical in-flight requests under high concurrency into a single execution.
   - Broadcasts the result to all waiting goroutines via Go channels.
2. **Decorator Pattern for Stores**:
   - Wrapping `InventoryStore` with `CachedInventoryStore` without modifying the underlying database implementation.
