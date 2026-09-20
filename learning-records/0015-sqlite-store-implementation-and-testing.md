# Implemented SQLite Store with database/sql and Upsert Semantics

The learner implemented `internal/store/task_sqlite.go` using standard `database/sql` with pure-Go `modernc.org/sqlite`:
1. Automatic table migration in constructor (`CREATE TABLE IF NOT EXISTS tasks`).
2. Translating `sql.ErrNoRows` into domain sentinel `queue.ErrTaskNotFound`.
3. Writing black-box tests in `package store_test` using `:memory:` SQLite and `t.Helper()`.
4. Identified the need for Upsert (`ON CONFLICT(id) DO UPDATE SET status = excluded.status`) to support lifecycle state updates on the same task ID.
All tests passed under ThreadSanitizer (`go test -count=1 -v -race ./...`).
