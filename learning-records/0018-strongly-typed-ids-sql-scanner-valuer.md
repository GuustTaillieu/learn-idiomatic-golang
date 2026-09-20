# Learning Record 0018: Strongly-Typed IDs & SQL Serialization (Valuer & Scanner)

## Date: 2026-09-19

## Milestone & Breakthrough
- The learner implemented strongly-typed domain IDs (`type ItemID uuid.UUID`, `type OrderID uuid.UUID`) with JSON marshaling and stringer methods.
- The learner implemented a complete Saga compensation mechanism in `internal/processor/multi.go` executing compensating closures in reverse order upon failure.
- Identified the root cause of test failures in `internal/store`: `database/sql` requires `driver.Valuer` and `sql.Scanner` to serialize custom Go types into SQL primitives.
- Created Lesson 0018 explaining how `driver.Valuer` and `sql.Scanner` bridge domain types to `database/sql`.

## Feedback Provided
1. **Strengths**:
   - Building domain NewType IDs (`OrderID`, `ItemID`) eliminates argument swapping errors at compile time.
   - Clean Saga pattern implementation with reverse-order rollback in `MultiProcessor`.
   - `RowsAffected()` check correctly implemented in `inventory_sqlite.go`.
2. **Corrections & Guidance**:
   - Explain that SQL can handle custom types natively once `driver.Valuer` and `sql.Scanner` are satisfied.
   - Warn about uninitialized `db *sql.DB` in `PlaceOrderProcessor`.
