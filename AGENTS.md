# Repository Guidelines

## Project Structure & Module Organization
- The module root is defined in `go.mod` as `yoshiyoshifujii/go-capability-token-relay-pattern`; keep orchestration code at the top level.
- Place runnable entry points under `cmd/<app>/main.go` (e.g., `cmd/tokenrelay/main.go`) and reserve `internal/<domain>` for capability-token logic that should not be imported externally.
- Adapters and fakes live under `internal/interface_adaptor/<kind>` (e.g., repository, service); use them to satisfy `internal/repository` and `internal/service` interfaces from the usecase layer. Tests that wire adapters should go under `internal/interface_adaptor/test`.
- Reusable adapters (HTTP clients, queue relays) that must be imported by multiple binaries belong in `pkg/<component>`; keep domain-facing contracts in `internal/repository` and `internal/service`.
- Check `README.md` and the linked Zenn article for architecture context, and mirror its sections (Issuer, Relay, Consumer) with matching package names.
- Store integration fixtures in `testdata/` beside the packages they verify and add diagrams or flow specs in `docs/` as the design evolves.

## Build, Test, and Development Commands
- `go mod tidy`: sync dependencies with the module graph before committing.
- `go fmt ./... && go vet ./...`: format and statically inspect the codebase; run this pair before opening a PR.
- `go build ./...`: compile every package to catch missing dependencies or build tags.
- `go run ./cmd/tokenrelay`: execute the reference relay binary for manual testing.
- `go test ./...` and `go test -race ./...`: run unit tests and data-race checks; both should be green in CI.

## Coding Style & Naming Conventions
- Rely on `gofmt` defaults (tabs for indentation, standard import grouping) and keep files ASCII unless protocol specs require otherwise.
- Use descriptive, lowercase package names (`issuer`, `relay`, `consumer`) and CamelCase identifiers for exported types such as `CapabilityToken`.
- Keep functions small and composable; prefer returning `(value, error)` and wrap upstream errors with `fmt.Errorf("context: %w", err)`.
- Log with structured key/value pairs and gate verbose output behind a `DEBUG` environment variable.

## Testing Guidelines
- Author table-driven tests in `_test.go` files with functions named `TestComponentBehavior`; encode scenario names in the table for clarity.
- Place golden responses or JWT fixtures under `testdata/` with sanitized payloads.
- Mock external dependencies (authorization servers, messaging layers) via interfaces in `internal/<domain>/ports.go`.
- Maintain ≥80% coverage on core relay flows; ensure `go test -race ./...` is part of pre-merge checks.

## Commit & Pull Request Guidelines
- Write imperative, 72-character max subjects like `Implement relay token signer`; elaborate in the body with motivation and risk notes.
- Reference issues using `Fixes #NNN` or `Refs #NNN` and summarize manual verification steps (commands, curl traces).
- Pull requests must describe protocol impacts, list any new environment variables, and attach sequence diagrams when changing token hand-offs.
- Request review once CI passes; respond to feedback with follow-up commits rather than force-pushing history when possible.

## Security & Configuration Tips
- Never commit live capability tokens; load them from environment variables or a `.env.local` file excluded by `.gitignore`.
- Capture configuration defaults in `configs/.env.example` and document them in `README.md`.
- Validate issuer signatures before relaying tokens, redact PII from logs, and expire cached credentials promptly.
