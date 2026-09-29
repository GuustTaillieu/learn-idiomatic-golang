# Learning Record 0029: Singleflight Test Determinism & Idiomatic Package Layout

## Date: 2026-09-29

## Context & Inquiries
1. **Singleflight Test Failure**: The learner implemented `singleflight.Group` deduplication in `internal/store/cached_inventory.go`, but `TestCachedInventory_Get` failed (`expected Get to be called once, but got 10`).
2. **Package Structure & Naming Architecture**: The learner asked whether structuring by domain folder (e.g. `order/` with `sql_store` and `model`, yielding `order.SQLStore` instead of `stores.OrderSQLStore`) is more idiomatic in Go.

## Architectural Findings & Analysis
1. **Singleflight Concurrency Semantics in Testing**:
   - `singleflight.Group` deduplicates calls that are concurrently in-flight at the exact same moment. It does not provide historical caching.
   - When a mock function has zero latency (10ns), goroutines in a test loop execute serially. Goroutine 1 completes before Goroutine 2 starts.
   - Deterministic testing of deduplication requires simulating I/O latency (`time.Sleep`) and releasing concurrent goroutines via a synchronized start channel (`close(start)`).

2. **Go Package Layout Consensus**:
   - **Package-by-Layer (`models`, `stores`, `controllers`)**: Anti-pattern. Leads to severe stutter (`stores.OrderStore`), cyclic dependency compilation errors (`import cycle not allowed`), and low cohesion.
   - **Ben Johnson's Standard Package Layout**: Domain types & interfaces in `internal/domain` (no dependencies). Implementations in technology subpackages (`internal/sqlite`, `internal/http`). Result: `sqlite.NewOrderStore(db)` with zero stutter.
   - **Package-by-Feature (Vertical Slices)**: High cohesion per bounded context (`internal/order`, `internal/inventory`). Requires cross-domain workflows to be orchestrated from outside to prevent circular imports.

3. **Created References**:
   - Lesson: `lessons/0029-package-layout-and-singleflight-testing.html`
   - Reference: `reference/0004-idiomatic-go-package-layout.html`
