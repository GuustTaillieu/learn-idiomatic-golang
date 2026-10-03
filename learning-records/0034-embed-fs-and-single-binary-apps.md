# Learning Record 0034: Single-Binary Applications with //go:embed & Real-Time Dashboards

## Date: 2026-09-29

## Context & User Achievements
- Learner implemented the Event Hub in `internal/event/hub.go` conforming to `domain.Hub[T any]`.
- Learner added `WithTaskHub` functional option to `queue.GenericQueue[T]` and broadcasted task status transitions.
- Learner implemented the SSE endpoint `GET /events` with `http.Flusher` and `io.Writer` streaming.
- Verified 100% test pass rate across all packages under the race detector.

## Architectural Concepts Introduced
1. **Single-Binary Web Applications (`embed.FS`)**:
   - Compiling static assets (HTML, CSS, JS, images, templates) directly into the Go executable via `//go:embed`.
   - Eliminating external asset dependencies, runtime path failures, and complex Docker multi-stage deployments.
2. **Subtree Filesystem Isolation (`io/fs.Sub`)**:
   - Extracting subdirectories from an `embed.FS` to serve from root `/` cleanly via `http.FileServer`.
3. **End-to-End Reactive Streaming Architecture**:
   - Browser client connects to `/events` using native `EventSource`.
   - Client POSTs to `/orders` &rarr; SQLite outbox persists `PENDING` &rarr; Dispatcher pushes to Queue &rarr; Worker processes in parallel &rarr; Hub broadcasts &rarr; SSE handler writes frame &rarr; Browser DOM updates in real-time.

## Artifacts Created:
- Lesson: `lessons/0034-embed-fs-and-single-binary-apps.html`
- Reference: `reference/0007-go-embed-and-static-file-serving.html`
- Frontend: `web/static/index.html` and `web/web.go`
