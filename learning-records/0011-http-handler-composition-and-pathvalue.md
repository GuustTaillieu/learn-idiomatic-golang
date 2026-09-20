# HTTP Handler Composition, PathValue, and E-Commerce Lifecycle

The learner implemented `internal/httpapi/handler.go` with Go 1.22+ ServeMux routing (`POST /tasks`, `GET /tasks/{id}`), JSON stream decoding/encoding, and structured `log/slog` logging.

Key senior refinement patterns identified:
1. Pure Dependency Injection: `NewHandler(queue Queue, store Store) *Handler` and `(h *Handler) Routes() http.Handler` instead of instantiating concrete components inside `NewHandler`.
2. Standard Path Extraction: Using `r.PathValue("id")` instead of manual slice manipulation (`r.URL.Path[len("/tasks/"):]`).
3. State Persistence on Ingestion: Saving `task` (in `StatusPending`) to `Store` at the HTTP ingress point before submitting to the queue, preventing 404 race conditions while tasks wait in the channel buffer.
4. Testing via `net/http/httptest`.
