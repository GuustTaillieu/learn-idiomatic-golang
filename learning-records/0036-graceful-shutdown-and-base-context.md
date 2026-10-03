# Learning Record 0036: Streaming Teardown & http.Server.BaseContext

## Date: 2026-10-02

## Context & Diagnosed Problem
- Learner tested the application and observed that pressing `Ctrl+C` hung for 5 seconds before logging:
  `ERROR Failed to shutdown server error="context deadline exceeded"`
- Root Cause: `http.Server.Shutdown(ctx)` closes listeners and waits for active connections to finish naturally. It does NOT cancel `r.Context()` on active HTTP connections.
- The open SSE stream (`GET /events`) in the browser had an infinite `for { select { ... } }` loop that never exited because `r.Context()` was never cancelled.

## Architectural Solutions
1. **`http.Server.BaseContext` Hook**:
   - Setting `BaseContext: func(l net.Listener) context.Context { return ctx }` links the application's root signal context to all incoming HTTP request contexts.
   - When `Ctrl+C` fires, `ctx` cancels, immediately triggering `<-r.Context().Done()` in the SSE handler.
   - Streaming handlers clean up in <1ms, enabling `srv.Shutdown()` to complete instantly with 0 errors.
2. **Defensive Termination Fallback**:
   - Calling `srv.Close()` if `srv.Shutdown()` times out ensures no rogue connections hold process exit hostage.

## Artifacts Created:
- Lesson: `lessons/0036-graceful-shutdown-and-base-context.html`
- Reference: `reference/0009-graceful-shutdown-and-base-context.html`
