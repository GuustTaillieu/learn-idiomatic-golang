# Learning Record 0022: Go Naming Conventions & Production Health Probes

## Date: 2026-09-21

## Milestone & Conceptual Insights
- Completed the full Transactional Outbox pattern with `OutboxDispatcher` running on a `time.Ticker` in `cmd/server/main.go`.
- Successfully resolved test failures in `internal/httpapi/handler_test.go` and verified tests run in 5s (down from 25s) with 0 data races.
- Codified Go naming conventions into a dedicated reference cheat sheet (`reference/0002-idiomatic-go-naming-conventions.html`):
  - Receiver names: 1-2 letters (`s` for Store, `q` for Queue, `h` for Handler).
  - Scope proportionality: Rob Pike's law that identifier length is proportional to scope distance.
  - Avoiding stuttering in package-symbol combinations.
- Introduced cloud-native health checks:
  - Liveness (`GET /healthz`) vs. Readiness (`GET /readyz` via `db.PingContext`).
  - Request-scoped structured logging with `slog.With`.
