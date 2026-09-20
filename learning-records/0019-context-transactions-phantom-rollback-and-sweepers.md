# Learning Record 0019: Context Transactions, Phantom Rollback Prevention & Recovery Sweepers

## Date: 2026-09-19

## Milestone & Code Evolution
- The learner successfully built context-based transaction propagation via `internal/lib/db_helper.go` (`DBTX`, `WithTx`, `TxFromContext`).
- The learner connected `PlaceOrderProcessor` to `InventoryStore` and `OrderStore` across a shared database transaction using `txCtx`.
- Verified all tests in `internal/processor`, `internal/store`, `internal/domain`, and `internal/httpapi` pass cleanly with zero data races.
- Identified an edge-case bug: returning `rollbackFunc, tx.Commit()` triggers compensating action (`ReleaseStock`) even when `tx.Commit()` fails, causing phantom stock creation.
- Identified redundant manual UUID string scanning that can now be cleaned up since `ItemID` implements `sql.Scanner`.
- Created Lesson 0019 introducing startup recovery sweepers for stranded `PENDING` orders after process crashes.

## Key Takeaways
1. Compensation closures must only be returned if the underlying forward transaction successfully committed.
2. In Go, `driver.Valuer` and `sql.Scanner` enable clean, direct struct scanning without intermediate string parsing.
3. In-memory queues require persistent stores and recovery sweeps to survive ungraceful process restarts.
