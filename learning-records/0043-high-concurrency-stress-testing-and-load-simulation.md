# Learning Record 0043: High-Concurrency Stress Testing & Load Simulation

## Date: 2026-10-05

## Context & Progression
- Evaluated the learner's error handling implementation using `errors.Join` and domain error categories.
- Analyzed the trade-offs: category classification is strong, but HTTP terminology in domain layers and redundant `errors.Unwrap` calls are senior pitfalls.
- Progressed to end-to-end load validation: testing the full distributed architecture under real concurrent pressure.

## Architectural Concepts Introduced
1. **Error Hierarchy Evaluation**:
   - Clean taxonomy: `ErrNotFound`, `ErrConflict`, `ErrValidation`, `ErrInternal`.
   - `errors.Is` vs `errors.Unwrap`: `errors.Is` automatically traverses error trees recursively, making manual unwrapping unnecessary and prone to breaking multi-error slices (`errors.Join`).
   - Separation of concerns: Keep HTTP status codes and phrases out of pure domain entities.
2. **Concurrent Load Generation in Go**:
   - Constructing custom stress test harnesses using worker pools (`chan int`, `sync.WaitGroup`, `atomic.Int64`).
   - Transport tuning: Configuring `MaxIdleConns` and `MaxIdleConnsPerHost` to prevent ephemeral TCP port exhaustion under heavy traffic.
   - Live system observation: Verifying rate limit backpressure (429), transactional outbox polling, and Server-Sent Events (SSE) streaming under saturation.

## Artifacts Created:
- Lesson: `lessons/0043-high-concurrency-stress-testing-and-load-simulation.html`
- Reference: `reference/0016-load-testing-and-stress-simulation.html`
