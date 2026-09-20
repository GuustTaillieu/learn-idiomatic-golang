# Mission: Master Idiomatic Go & Concurrent Systems

## Why
Transition from basic Go syntax knowledge to writing professional, industry-standard senior Go. Prepare for future senior Go engineering roles by understanding the Go philosophy (simplicity, mechanical sympathy, explicit error handling) and building production-grade concurrent software.

## Success looks like
- Design and implement a production-grade Concurrent Task Queue Service from scratch.
- Reason about and confidently apply Go concurrency primitives (`sync.Mutex`/`RWMutex`, `chan`, `sync.WaitGroup`, `select`, `context.Context`, `errgroup`) without race conditions or goroutine leaks.
- Multi-step concurrent pipeline: Fan-out processing (inventory check, billing, receipt generation) with coordinated error cancellation.
- Real database persistence with Go's `database/sql` (SQLite): Atomic transactions (`tx.Commit`/`tx.Rollback`) to prevent half-done orders and outbox event guarantees.
- Adhere to idiomatic conventions naturally: "Accept interfaces, return structs", domain-driven package layout (`internal/queue`, `internal/store`, `internal/httpapi`), short receiver names, explicit error wrapping (`%w`), and zero "log-and-return".
- Write idiomatic table-driven tests and verify concurrency correctness using the Go race detector (`go test -race`).

## Constraints
- **Hands-on first**: Guide the user with handles, mental models, and Socratic hints. Never dump raw completed code solutions; let the user write, test, and debug their own code.
- **Strictly standard library & lightweight tooling**: Master Go's native capabilities (`net/http`, `sync`, `database/sql`, `log/slog`, `context`) before reaching for third-party monolith frameworks.

## Out of scope
- Generic 50k-star GitHub boilerplate layouts (`pkg/`, `api/`, `deployments/`).
- Bloated ORMs (GORM) — master native `database/sql` transactions first.
- Complex distributed cluster deployments (focus is language idioms, architecture, concurrency, and data safety).

