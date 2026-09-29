# Learning Record 0030: Item Domain Polish & Generic Singleflight

## Date: 2026-09-29

## Context & Inquiries
- Learner refactored `internal/domain/item.go` implementing `ItemID`, `ItemOption`, `sql.Scanner`, `driver.Valuer`, and `json.Marshaler`/`Unmarshaler`.
- Learner restructured the entire project following Ben Johnson's Standard Package Layout (`internal/sqlite`, `internal/memory`, `internal/http`, `internal/queue`, `internal/health`).
- Learner requested detailed feedback on `internal/domain/item.go` and to proceed with lessons.

## Detailed Review & Insights
1. **Domain Entity Construction & Functional Options**:
   - Mandatory parameters (like `name`) passed via standard constructor args, while non-mandatory/hydration fields (`WithID`, `WithCreatedAt`) passed via options.
   - For pure data structs with exported fields, direct struct literals (`Item{...}`) are zero-allocation, but options are useful in factories and DB hydration.
2. **Scanner & Valuer Nuance (Binary UUID vs String)**:
   - When scanning `[]byte`, distinguishing between 16-byte raw binary representations (`uuid.FromBytes`) and 36-character ASCII representations (`uuid.Parse`) is necessary for multi-database compatibility.
3. **Singleflight Scope Lifespan**:
   - `singleflight.Group` cannot be scoped inside a standalone function call (`var sf singleflight.Group`), because concurrent goroutines each create their own instance. Singleflight groups must be long-lived struct fields.
4. **Generic Singleflight Design**:
   - Implemented `lib.Singleflight[T any]` to eliminate boilerplate type-assertions at call sites.

## Artifacts Created:
- Lesson: `lessons/0030-generic-singleflight-and-benchmarking.html`
