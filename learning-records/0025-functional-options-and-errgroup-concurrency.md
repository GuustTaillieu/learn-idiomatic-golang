# Learning Record 0025: Functional Options & errgroup Parallel Fan-Out

## Date: 2026-09-25

## Milestone & Breakthroughs
- The learner independently implemented the Functional Options Pattern in `internal/domain/item.go` (`optsFunc`, `WithID`, `WithCreatedAt`).
- Verified implementation of request-scoped structured logging (`internal/lib/logger.go`) and request correlation middleware (`internal/httpapi/middleware/request_id.go`).
- Reviewed Functional Options and provided the canonical senior simplification: mutating `*Item` directly without an intermediary `Opts` struct.
- Identified and pointed out a default value omission in `item.go` (`Name: ""` instead of `Name: name`).
- Installed `golang.org/x/sync` (upgraded to v0.23.0).
- Created Reference Sheet `reference/0003-functional-options-pattern.html`.
- Created Lesson 0025 introducing `errgroup.WithContext(ctx)` for parallel fan-out concurrency with coordinated error cancellation.
