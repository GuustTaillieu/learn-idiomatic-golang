# Idiomatic Golang: Senior Engineering Guidelines & Roadmap
*Extracted from NotebookLM notebook: "Idiomatic Golang"*

---

## 1. Project Overview & Hands-on Goal

**Target Exercise:** Building a **Production-Grade Concurrent Task Queue Service** in Go.

### High-Level Architecture
- **`cmd/server/main.go`**: Entry point wiring dependencies, reading config, and handling lifecycle.
- **`internal/queue/`**: Core concurrent worker-pool engine using goroutines, channels, and `context.Context` for graceful shutdown.
- **`internal/store/`**: In-memory / persistence storage satisfying small, consumer-defined repository interfaces.
- **`internal/httpapi/`**: REST API handlers with structured logging (`log/slog`) and custom error writing middleware.
- **Testing**: Table-driven unit tests, native mock-free stubs, race detector (`go test -race`).

---

## 2. Senior Go Design Principles & Conventions

### A. Tooling & Automation
- **Formatting**: Strictly format with `gofmt` (or `goimports`). Never debate style in PRs.
- **Linters & Static Analysis**: Run `go vet`, `golangci-lint` / `revive`.
- **Semicolon Insertion Rule**: Opening braces `{` must always stay on the same line as statements/declarations.
- **Modernization**: Use `go fix ./...` to automate updating syntax across Go versions.

### B. Project Layout (Avoiding the 50k-Star Myth)
- **Start Flat**: Begin with files directly in the root directory for single packages / small services.
- **`internal/` for Privacy**: Move business logic into `internal/` only when compiler-enforced package privacy is needed.
- **`cmd/` for Multi-Binary**: Add `cmd/<binary-name>/main.go` only when shipping multiple binaries or both a library and CLI tool.
- **Package by Domain**: Name packages by what they provide (`queue`, `store`, `httpapi`), never technical layers (`models`, `handlers`, `helpers`, `utils`, `common`). Never create `src/`.

### C. Naming & Style
- **Package Names**: Short, lowercase, single-word names (`queue`, `store`, `auth`).
- **Receiver Names**: Short 1–2 letter abbreviations matching the type name (e.g., `q *Queue`, `p Processor`). Never use `self`, `this`, or `me`.
- **Variable Scope**: Short names (`i`, `n`, `err`, `ctx`) for narrow local scopes; clear CamelCase for wider scopes.
- **Control Flow**: Flat control flow with early guard returns. Avoid blank/naked returns in functions with named return values.
- **Context First**: Always pass `ctx context.Context` as the explicit first parameter to functions performing I/O, concurrency, or lifecycle management.

### D. Senior Error Handling
- **Errors as Values**: Errors are regular return values returned as the final tuple element.
- **Wrapping with `%w`**: Wrap errors at boundaries with `fmt.Errorf("operation failed: %w", err)` to preserve causal chains.
- **`errors.Is` vs `errors.As`**:
  - Use `errors.Is` for sentinel error comparisons (`var Err... = errors.New(...)`).
  - Use `errors.As` (or Go 1.26 `errors.AsType`) for custom error structs (must implement `Unwrap() error`).
- **Never Log-and-Return**: Do not log an error and return it. Log once at the outermost boundary (e.g., HTTP handler) with `log/slog`, and wrap/return silently everywhere below.

### E. Interfaces & Composition
- **Accept Interfaces, Return Structs**: Functions take interface parameters to remain testable and decoupled, but constructors return concrete structs.
- **Consumer-Defined Interfaces**: Define small 1–2 method interfaces at the consumption site, not upfront alongside implementations.
- **Composition over Inheritance**: Use embedded struct fields to compose functionality without deep hierarchies.
- **No Boilerplate Getters/Setters**: Access fields directly; reserve methods for state transitions and business logic.

### F. Concurrency & Synchronization
- **Channels**: Unbuffered for direct handoffs; buffered/bounded for worker pools and backpressure.
- **Leak Prevention**: Clean up goroutines via context cancellation or done channels (`GOEXPERIMENT=goroutineleakprofile`).
- **Mutexes**: Use `sync.Mutex` or `sync.RWMutex` for mutable state. Avoid `sync.Map` unless specifically read-heavy or key-disjoint.

---

## 3. Sources in the Notebook

1. **5 Tips for Writing Idiomatic Code in Golang** – Mario Carrion
2. **ardanlabs/gotraining** – William Kennedy (Ardan Labs)
3. **GoBooks (Curated Directory)** – Dariush Abbasi & Community
4. **skills-best-practices** – Minko Gechev
5. **Go 1.26 Release Notes** – The Go Team (Google)
6. **Go 1.26: What's New and Why It Matters** – Travis Media
7. **Go Error Handling in 2026: The Patterns I Actually Ship** – Rayyan / Abrarqasim Blogs
8. **Golang 1.26: Performance, Generics, Error Handling** – Tech Analysis
9. **Grill Me Skill: The Deceptively Simple AI Prompt** – AlphaMatch AI / Matt Pocock
10. **How to Structure a Go Project (2026)** – LevelUpGo
11. **How to Create Good AI Agent Skills** – Thales Assis
12. **Introduction to Agent Skills** – Claude Academy / Anthropic
13. **Learning Go: An Idiomatic Approach to Real-World Go Programming** – Jon Bodner (O'Reilly)
14. **Learn the AI SDLC – Building Agent Skills** – Sarvesh Talele (freeCodeCamp)
15. **Matt Pocock's Skills, Actually Explained** – Kaitlin Z. Albasi
16. **Skill Authoring Best Practices** – Anthropic
17. **Superpowers vs Agent Skills vs Pocock** – Jamil (jamilxt)
18. **The /writing-for-agents Skill** – Matt Pocock (AI Hero)
19. **The Best Golang Books in 2026** – Eric (DEV Community)
20. **Writing Clean and Idiomatic Go** – Talo Oyweka
21. **mattpocock/skills Repository** – Matt Pocock

---

## 4. Current Progress & Next Steps (Phase 1)

- **Completed**:
  - `go.mod` initialized (`idiomatic-go`)
  - Initial sanity check in `main.go`
- **Next Up (Phase 1: Domain Modeling & Internal Privacy)**:
  - Create `internal/queue/task.go`:
    - Sentinel errors: `ErrTaskNotFound`, `ErrQueueClosed`
    - `TaskStatus` enum (`StatusPending`, `StatusRunning`, `StatusCompleted`, `StatusFailed`)
    - `Task` struct (`ID`, `Payload`, `Status`, `CreatedAt`)
    - `Processor` interface (`Process(ctx context.Context, payload string) error`)
  - Create `internal/queue/queue.go`:
    - `Queue` struct (`processor Processor`, `tasks chan Task`, `wg sync.WaitGroup`, `mu sync.RWMutex`, `closed bool`)
    - `New(p Processor, bufferSize int) *Queue`
    - `Submit(ctx context.Context, t Task) error`
