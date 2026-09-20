# Solving Test Import Cycles: The Black-Box Testing Pattern (package <name>_test)

The learner encountered a classic Go compiler constraint: `import cycle not allowed` when importing `internal/store` inside `internal/queue/queue_test.go` (since `store` already imports `queue`).

The idiomatic Go solution:
1. Go explicitly supports `package <name>_test` inside the same directory as the package under test.
2. `package queue_test` compiles as a separate, external consumer package.
3. This eliminates cyclic dependencies, allowing the test to import both `internal/queue` and `internal/store` simultaneously while testing the public API as an outside consumer (black-box testing).
