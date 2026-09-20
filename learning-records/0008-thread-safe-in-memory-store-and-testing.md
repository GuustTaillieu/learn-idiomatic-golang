# Implemented Thread-Safe In-Memory Store with RWMutex

The learner authored `internal/store/memory.go` and `internal/store/memory_test.go`:
1. Guarded an in-memory `map[string]*queue.Task` using `sync.RWMutex`, preventing `fatal error: concurrent map writes`.
2. Differentiated write locks (`s.mu.Lock()` in `Save`) from read locks (`s.mu.RLock()` in `GetById`).
3. Returned sentinel errors (`queue.ErrTaskNotFound`) when lookups fail.
4. Validated with table-free native unit tests using `errors.Is` and verified race safety via `go test -count=1 -v -race ./internal/store`.
