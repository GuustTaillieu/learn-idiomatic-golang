# Learning Record 0041: 12-Factor Configuration & Fail-Fast Startup

## Date: 2026-10-03

## Context & Progression
- Learner successfully registered `pprof` diagnostic routes (`/debug/pprof/`) in `internal/http/handler.go`.
- All unit tests across all 8 packages pass with 0 data races.
- The server is functionally complete with embedded web dashboard, SSE streaming, rate limiting, and access logging.

## Architectural Concepts Introduced
1. **The 12-Factor Config Principle**:
   - Eliminating hardcoded configuration constants (`":8080"`, database paths, worker counts) in favor of environment variable injection.
2. **Standard Library Configuration vs Third-Party Frameworks**:
   - Avoiding heavy dependencies (like Viper) by relying on Go standard library primitives: `os.LookupEnv`, `strconv`, and `flag`.
3. **Fail-Fast Validation at Bootstrap**:
   - Validating all configuration invariants (e.g. `WorkerCount > 0`, non-empty database URLs) *before* binding network ports or establishing database pools.
   - Preventing services from booting into half-broken or corrupted runtime states.

## Artifacts Created:
- Lesson: `lessons/0041-twelve-factor-configuration-and-fail-fast-startup.html`
- Reference: `reference/0014-configuration-and-fail-fast-startup.html`
