# Plan: Zenn Article Sample Implementation

## Objective
- Build a runnable Go sample that demonstrates the capability token relay pattern described in the linked Zenn article, aligning code layout with `cmd/`, `internal/`, and `pkg/` conventions captured in `AGENTS.md`.

## Inputs & Dependencies
- Obtain the full text (or a reliable summary) of https://zenn.dev/yoshiyoshifujii/articles/a55815241c3213; request the author if direct access is blocked in the current environment.
- Confirm Go toolchain version (target `go1.25.3` per `go.mod`) is available for local builds.
- Identify any external services or mock endpoints referenced in the article (issuer, relay, consumer, storage).

## Phase 1: Article Research
1. Read the Zenn article and extract its core flow: actors, data contracts, token lifetimes, and security considerations.
2. Capture notes on sequence diagrams or pseudo-code; flag any ambiguous steps that need clarification from the author.
3. Translate article terminology into repository vocabulary (e.g., map “Issuer” to `internal/issuer`).

## Phase 2: Architecture & Scope Definition
1. Draft a lightweight architecture summary (component list, responsibilities, data flow) in `docs/architecture.md`.
2. Decide which parts become runnable binaries (`cmd/issuer`, `cmd/relay`, `cmd/consumer`) versus shared libraries (`pkg/tokens`, `pkg/httpclient`).
3. Define configuration strategy (env vars, `.env` templates) and minimal sample data required to exercise the flow.

## Phase 3: Sample Implementation
1. Scaffold directories and `main.go` entry points for each actor; wire basic CLI/config parsing.
2. Implement token issuance, relay forwarding, and consumer validation according to article guidance, reusing shared helpers.
3. Add interfaces to abstract external dependencies (e.g., authorization servers), plus in-memory or file-backed mocks for the sample.
4. Embed sanitized example tokens/fixtures under `testdata/` for deterministic tests and demonstrations.

## Phase 4: Testing & Validation
1. Write table-driven unit tests for each package (`issuer`, `relay`, `consumer`) covering happy path and error scenarios.
2. Create an integration test or scripted workflow (`make demo` or `scripts/run_demo.sh`) that runs the full issuance→relay→consumption pipeline.
3. Run `go fmt ./...`, `go vet ./...`, `go test ./...`, and `go test -race ./...`; document any platform-specific caveats.

## Phase 5: Documentation & Developer Experience
1. Update `README.md` with setup instructions, command examples, and a link to the new architecture summary.
2. Add troubleshooting tips and security notes (handling secrets, token storage) referencing the article’s recommendations.
3. Capture manual verification steps and demo outputs in `docs/demo-guide.md` (including curl or HTTPie commands).

## Open Questions / Next Steps
- Do we need to integrate with real external systems (e.g., Azure AD, custom issuer), or keep everything mocked?
- Are there performance or scalability targets mentioned in the article that should influence the sample?
- After the initial sample, consider adding observability (log tracing) or CI workflows to keep the demo healthy.
