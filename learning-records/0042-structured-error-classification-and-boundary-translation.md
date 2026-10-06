# Learning Record 0042: Structured Error Classification & Boundary Translation

## Date: 2026-10-03

## Context & Progression
- Fixed the flag redefinition panic in `internal/config/config.go` by replacing global `flag.StringVar` calls with an isolated `flag.NewFlagSet("server", flag.ContinueOnError)`.
- Verified all unit tests across all packages pass with 0 race conditions (`go test -count=1 ./...`).
- Advanced to error architecture: diagnosing leaky abstractions where storage failure modes (database down, timeout) were mistakenly reported as client errors (404 Not Found).

## Architectural Concepts Introduced
1. **The Global FlagSet Gotcha**:
   - `flag.StringVar` modifies package-level global singleton `flag.CommandLine`. Calling it across multiple unit tests triggers `panic: flag redefined`.
   - Solution: Use `flag.NewFlagSet` scoped to each parse call, allowing clean unit tests and flexible CLI argument overrides.
2. **Layered Error Architecture**:
   - Domain errors (`internal/domain`): pure sentinel errors and business invariants (`ErrOrderNotFound`, `ErrInvalidOrder`, `ErrInsufficientStock`).
   - Repository errors (`internal/sqlite`): wraps low-level driver errors with domain sentinels (`fmt.Errorf("...: %w", domain.ErrOrderNotFound)`).
   - API boundary translation (`internal/http`): translates errors via `errors.Is` into HTTP status codes (404, 400, 409, 500, 504) without leaking SQL details or sensitive storage metadata to HTTP callers.

## Artifacts Created:
- Lesson: `lessons/0042-structured-error-classification-and-boundary-translation.html`
- Reference: `reference/0015-structured-errors-and-api-translation.html`
