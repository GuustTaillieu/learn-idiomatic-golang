# Learning Record 0017: Domain Modeling vs Infrastructure & Transaction Boundaries

## Date: 2026-09-17

## Milestone & Conceptual Breakthrough
- The learner independently transitioned from generic infrastructure (`Task`) to a concrete business domain model (`Order`, `Stock`, `Item`).
- The learner attempted to build `SafeDBMultiProcessor` wrapping multiple sub-processors in a database transaction (`BeginTx` / `Commit`), discovering the fundamental transaction propagation dilemma (`db.ExecContext` vs `tx.ExecContext`).
- Clarified the two architectural levels:
  1. Generic infrastructure job queue (`Task` / `Job`).
  2. Business domain entity (`Order`).
- Introduced the standard Go `DBTX` interface pattern used by `sqlc` for transaction-agnostic stores.
- Highlighted the boundary between local DB ACID transactions vs distributed Sagas (e.g. external HTTP payment gateways).
- Identified a common SQL gotcha: `ExecContext` returns `err == nil` on zero-row updates; checking `RowsAffected() == 0` is required to detect failure.

## Next Steps
- Implement `RowsAffected()` check in `ReserveStock`.
- Wire atomic order creation and stock reservation in SQLite.
- Run background pipeline on orders committed to the database.
