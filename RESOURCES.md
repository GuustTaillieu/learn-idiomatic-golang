# Idiomatic Go Resources

## Knowledge

- [Book: _Learning Go: An Idiomatic Approach to Real-World Go Programming_ (2nd Ed) by Jon Bodner](https://www.oreilly.com/library/view/learning-go-2nd/9781098139285/)
  Authoritative text on idiomatic patterns, pointers, heap vs stack, implicit interfaces, and concurrency. Use for: foundational idioms, tooling, and error handling.
- [Training: _Go Training Class Material (gotraining)_ by William Kennedy (Ardan Labs)](https://github.com/ardanlabs/gotraining)
  Exemplary syllabus on mechanical sympathy, data-oriented design, memory profiling, and concurrency semantics. Use for: goroutine mechanics and performance decisions.
- [Official: _Go 1.26 Release Notes & Standard Documentation_ by The Go Team](https://go.dev/doc/)
  Primary source for language spec, Green Tea GC, expression-based `new()`, and standard package docs (`sync`, `context`, `log/slog`). Use for: language accuracy and modern idioms.
- [Guide: _How to Structure a Go Project_ by LevelUpGo](https://levelup.gitconnected.com/)
  Practical explanation debunking the 50k-star layout myth. Use for: directory structure evolution (`internal/`, `cmd/`).
- [Article: _Go Error Handling in 2026: The Patterns I Actually Ship_ by Rayyan](https://abrarqasim.com/)
  Real-world boundary wrapping, sentinel errors vs typed structs, and eliminating "log-and-return". Use for: production error design.

## Wisdom (Communities)

- [r/golang](https://reddit.com/r/golang)
  High-signal Reddit community focused on idiomatic code reviews, architecture discussions, and new Go developments. Use for: community consensus and code critique.
- [Gophers Slack](https://gophers.slack.com/)
  The primary global community for Go developers. Channels like `#newbie`, `#performance`, and `#architecture` provide direct interaction with core contributors and senior practitioners.
- [Go Forum](https://forum.golangbridge.org/)
  Official discussion forum for architecture and in-depth Go design questions.
