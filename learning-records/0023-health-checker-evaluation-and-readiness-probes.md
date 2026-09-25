# Learning Record 0023: Health Checker Evaluation & Deep Readiness Probes

## Date: 2026-09-23

## Milestone & Conceptual Evaluation
- The learner implemented `HealthChecker` interface in `internal/httpapi/handler.go` with `/healthz` and `/readyz` routes.
- The learner attached `Ping(ctx)` to `Queue` (`queue.Ping`) and wired it into `main.go`.
- Evaluated the architectural merits and blind spots of queue-only health checks:
  - *Strength*: Detects queue shutdown state (`q.closed`).
  - *Blind spot*: Misses database failures. If SQLite is locked or disk full, `/readyz` returns 200 while all `POST /orders` fail with 500.
- Introduced the Composite Health Checker pattern to inspect both `db.PingContext(ctx)` and `queue.Ping(ctx)`.
- Highlighted the necessity of `context.WithTimeout` on health check probes to prevent hanging HTTP sockets.

## Next Steps
- Implement composite health checking (DB + Queue).
- Add 503 Service Unavailable negative assertion test in `handler_test.go`.
