# Learning Record 0051: Standalone Type-Safe Client Architecture

## Date: 2026-10-05

## Context & Progression
- Addressed learner's goal of building an external, standalone client that communicates with the backend in the most optimal and type-safe way without sharing internal code.
- Contrasted internal domain boundaries (`internal/`) with public API contract distribution (`proto/`).

## Architectural Concepts Introduced
1. **Public Contract vs Private Implementation**:
   - Go's compiler strictly forbids importing code from an `internal/` tree across module boundaries.
   - Public contract definition in `proto/order/v1/` and `proto/order/v1/orderv1connect/` enables any external service, CLI, or repo to consume the API without coupling to the server's database or business logic.
2. **Wire Format Optimization**:
   - Connect protocol supports binary Protocol Buffers (`application/proto`), eliminating JSON text serialization overhead, reducing packet size by ~70%, and achieving microsecond serialization speeds.
   - Compile-time type safety: RPC parameters and responses are validated at build time, preventing runtime schema mismatches.
3. **Multi-Language Frontend Consumption**:
   - The same Protocol Buffer schema generates TypeScript clients (`@connectrpc/connect-web`) for React/Next.js/Vite frontends with full autocompletion and zero Envoy proxy requirements.

## Artifacts Created:
- Lesson: `lessons/0051-standalone-typesafe-client-architecture.html`
- Reference: `reference/0024-standalone-typesafe-client-architecture.html`
