# Learning Record 0045: Embedded Database Migrations & Idempotent Seeding

## Date: 2026-10-05

## Context & Progression
- Evaluated learner's production `Dockerfile` with multi-stage build (`golang:1.27.1-alpine` to `gcr.io/distroless/static-debian12:nonroot`) and `.dockerignore`.
- Learner identified the operational bottleneck during stress testing: fresh databases lack items and initial stock, causing saga stock deductions to fail.
- Advanced to database lifecycle management: replacing scattered DDL with embedded versioned migrations and idempotent catalog seeders.

## Architectural Concepts Introduced
1. **Separation of DDL (Schema) and DML (Queries)**:
   - Moving table definitions out of individual repository constructors (`NewOrderStore`, `NewInventoryStore`).
   - Preventing race conditions and ordering bugs on foreign key creation during service bootstrap.
2. **Go `embed.FS` Migration Engine**:
   - Embedding versioned SQL files directly into the compiled Go binary.
   - Deterministic execution sorted by filename (`00001_init.sql`, `00002_seed.sql`).
3. **Idempotent Data Seeding**:
   - Using `INSERT OR IGNORE` to safely initialize catalog items and inventory without failing on subsequent container restarts.
   - Enabling end-to-end stress tests and development environments to work immediately out-of-the-box.

## Artifacts Created:
- Lesson: `lessons/0045-embedded-database-migrations-and-seeding.html`
- Reference: `reference/0018-embedded-database-migrations-and-seeding.html`
