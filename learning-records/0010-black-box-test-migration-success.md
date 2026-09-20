# Applied Black-Box Testing Pattern (package queue_test)

The learner migrated `queue_test.go` to `package queue_test`:
1. Successfully eliminated the circular dependency between `internal/queue` and `internal/store` without introducing synthetic testing directories.
2. Verified integration behavior strictly through exported APIs.
3. Repositioned `q.Stop()` outside the task submission loop, validating concurrent batch processing and persistent status transitions (`StatusCompleted`) across all tasks.
ThreadSanitizer passed across the entire workspace (`go test -count=1 -v -race ./...`).
