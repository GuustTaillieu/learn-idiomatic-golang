# Learning Record 0031: Go 1.27 Standard Library UUID & Text Marshaling

## Date: 2026-09-29

## Context & Gotchas Diagnosed
1. **Standard Library `uuid` Package (Go 1.27)**:
   - Defined as `type UUID [16]byte`.
   - Does not have `uuid.FromBytes(b)` because Go's built-in type system allows converting slice to array directly: `uuid.UUID(b)` or `uuid.UUID([16]byte(b))`.
2. **`json.Unmarshaler` vs `encoding.TextUnmarshaler` Quote Trap**:
   - `UnmarshalJSON(data []byte)` receives the raw JSON token, including enclosing double quotes (`"f81d4fae-..."`).
   - Passing `data` directly to `uuid.UUID.UnmarshalText(data)` fails with `invalid uuid` because the first character is `"`.
   - Two idiomatic solutions:
     - Use `json.Unmarshal(data, (*uuid.UUID)(id))` to strip quotes.
     - Or implement `MarshalText` / `UnmarshalText` on the type directly without `MarshalJSON` / `UnmarshalJSON`. `encoding/json` natively delegates to `TextMarshaler` and handles quoting/unquoting automatically.
3. **Generic Singleflight Verification**:
   - Learner successfully implemented `lib.Singleflight[T any]` in `internal/lib/db_cache.go`.
   - Verified thread safety and passing unit tests under race detector.
