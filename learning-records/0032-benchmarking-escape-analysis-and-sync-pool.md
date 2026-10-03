# Learning Record 0032: Benchmarking, Escape Analysis & Zero-Allocation sync.Pool

## Date: 2026-09-29

## Milestone & Verified Execution
- All tests across all 7 packages pass with 0 data races (`go test -v -race ./...`).
- Packages: `internal/domain`, `internal/http`, `internal/lib`, `internal/memory`, `internal/processor`, `internal/queue`, `internal/sqlite`.
- Ben Johnson's Standard Package Layout successfully verified.

## Concepts Introduced
1. **Mechanical Sympathy & Runtime GC Pressure**:
   - The cost of heap allocations in high-throughput Go services (tail latency, GC pause frequency).
2. **Benchmarking with `testing.B`**:
   - `b.ResetTimer()`, `b.ReportAllocs()`, and the idiomatic `for b.Loop()` pattern.
   - Command-line flags: `-bench=.`, `-benchmem`, `-benchtime=3s`, `-count=5`.
3. **Compiler Escape Analysis**:
   - Inspecting stack vs heap decisions with `go build -gcflags="-m"`.
4. **Object Pooling via `sync.Pool`**:
   - Reusing temporary buffers in hot paths to eliminate heap allocations.
   - The critical requirement: resetting object state before reuse (`buf.Reset()`) to prevent cross-request data leaks.

## Artifacts Created:
- Lesson: `lessons/0032-benchmarking-escape-analysis-and-sync-pool.html`
- Reference: `reference/0005-benchmarking-and-zero-allocation.html`
