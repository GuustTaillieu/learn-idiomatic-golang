# Learning Record 0035: Outbox Claim Check, State Machine Lifecycle & Race-Free Persistence

## Date: 2026-10-02

## Context & Diagnosed Architectural Problem
- Learner observed an infinite dispatch loop when running the live application: `OutboxDispatcher` repeatedly re-polled and submitted orders because their database status remained `PENDING`.
- Learner attempted a fix by creating `OrderMarker` and running it in `processor.Parallel`.
- Learner correctly identified that this introduced failure and race condition risks.

## Architectural Findings & Solutions
1. **The Outbox Infinite Loop Mechanism**:
   - In polling-based outboxes, submitting an in-memory pointer to a channel does not mutate database rows.
   - If the database status is not transitioned immediately, subsequent ticker ticks (e.g. 500ms) re-read the same row and re-submit it indefinitely.
2. **Pitfalls of Parallel State Mutation**:
   - Running `OrderMarker` inside `processor.Parallel` with `OrderPlacing` triggers a concurrent write race on the SQLite database row for the same order ID.
   - Separating business domain steps (reserve stock, bill payment) from orchestrator lifecycle transitions (`PENDING` &rarr; `RUNNING` &rarr; `COMPLETED`) is essential.
3. **The 3-Step Lifecycle Architecture**:
   - **Step 1 (Claim at Dispatch)**: `OutboxDispatcher.DispatchOnce` transitions `order.Status` to `RUNNING` and saves to SQLite before queue submission.
   - **Step 2 (Atomic Success)**: `OrderPlacing` sets `order.Status = COMPLETED` and saves via `txCtx` before `tx.Commit()`.
   - **Step 3 (Compensating Rollback)**: On error, rollback sets `order.Status = FAILED` and persists to SQLite.

## Artifacts Created:
- Lesson: `lessons/0035-outbox-claim-check-and-state-machine-lifecycle.html`
- Reference: `reference/0008-outbox-claim-check-and-state-machines.html`
