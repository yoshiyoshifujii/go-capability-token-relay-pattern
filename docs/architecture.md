# Architecture Overview: Capability Token Relay Sample

## Objective & Scope
- Deliver a runnable Go sample that mirrors the Zenn article’s confirmed token relay pattern while staying aligned with repository guidelines in `AGENTS.md`.
- Focus on a CLI-driven workflow: users trigger issuance, relay verification, and consumption via commands rather than standing HTTP services.
- Model AWS KMS signing through a `KMSService` interface with an in-memory mock to keep the sample self-contained.

## Component Breakdown
### Issuer (`cmd/issuer`, `internal/issuer`)
- Validates domain-specific input (coupon/points/payment placeholders) and issues signed JWT capability tokens.
- Depends on `pkg/tokens` for claim structs and signing helpers, and on `pkg/config` for CLI/environment configuration.
- Exposes a CLI that outputs a confirmed token to STDOUT or a file in `testdata/tokens/`.

### Relay (`cmd/relay`, `internal/relay`)
- Accepts confirmed tokens from upstream clients, verifies signatures via `KMSService`, and enforces metadata checks (`aud`, `exp`, `jti`, `order_processing_id`).
- Maintains in-memory `jti` usage tracking per session; persistence is intentionally omitted.
- Emits sanitized audit logs (token IDs, order processing IDs) and forwards validated tokens to the consumer CLI.

### Consumer (`cmd/consumer`, `internal/consumer`)
- Receives forwarded tokens, asserts capability constraints relevant to its domain, and simulates business action (e.g., marking coupon redeemed).
- Stores lightweight state in memory; optional fixtures live under `testdata/consumer/`.

### Shared Packages
- `pkg/tokens`: JWT encoding/decoding, claim definitions, kid-aware verification utilities.
- `pkg/config`: Centralized environment + flag parsing for CLI binaries.
- `internal/infra/kmsmock`: Mock KMS signer/verifier implementation returned by a factory for local runs.

## Data & Control Flow
1. Issuer CLI validates input and calls `tokens.NewConfirmedToken` to sign a JWT using the mock KMS key.
2. Relay CLI receives the token (file path or STDIN), verifies signature and claims, records `jti`, and forwards the token reference.
3. Consumer CLI validates domain-specific constraints and prints simulated processing results.

## Configuration & Secrets
- Environment variables prefixed with `CTR_` configure each binary (e.g., `CTR_KMS_KEY_ID`, `CTR_ORDER_ID`).
- `.env.example` will document defaults; secrets are never committed.
- Mock KMS keys are generated at runtime and never leave the process.

## Testing & Fixtures
- Unit tests: table-driven tests per package under `internal/...` and `pkg/...`.
- Integration script: `scripts/demo.sh` to chain issuer → relay → consumer using temporary token files.
- Golden tokens and sample payloads reside in `testdata/tokens/` and `testdata/requests/`.

## Migration Considerations
- Real KMS or JWKS integration can replace `internal/infra/kmsmock` by implementing the same interface.
- HTTP transport can be layered later by wrapping CLI logic in handlers under `pkg/transport/http`.
