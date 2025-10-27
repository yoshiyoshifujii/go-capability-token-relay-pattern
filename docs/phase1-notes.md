# Phase 1 Notes: Confirmed Token Relay Pattern

## Article Highlights
- Microservice orchestration often leads to central services accumulating business rules, creating tight coupling and scaling risks.
- The article proposes relaying confirmed capability tokens so the orchestrator validates signatures only and delegates rule enforcement to domain services.
- Tokens are short-lived JWTs issued per domain, carrying capability claims and constraints tied to a specific order-processing session.

## Actors & Responsibilities
- **Client/User**: Initiates checkout, selects coupon/points/payment combinations.
- **OrderProcessing (Orchestrator)**: Relays domain-issued tokens, verifies signatures and metadata (`aud`, `exp`, `jti`, `order_processing_id`), and forwards requests without interpreting business semantics.
- **Domain Services** (`Coupons`, `Points`, `Payments`): Validate user input, issue signed confirmation tokens encapsulating domain-specific constraints, and later redeem tokens when orchestrator relays them.

## Token Contract
- Format: JWT (likely JWS) with issuer-specific signing keys; consider `kid` rotation.
- Mandatory claims: `iss`, `aud`, `sub`, `iat`, `exp`, `jti`, `order_processing_id`, domain-specific `capabilities` and `constraints`.
- Replay protection hinges on single-use `jti`; orchestrator must track seen IDs per session.
- Sensitive data should remain outside the JWT or use JWE for encrypted payload sections.

## Security & Operational Notes
- Tokens expire within minutes; orchestrator should reject stale tokens and avoid persisting secrets.
- Each domain governs issuance and revocation; orchestrator has no authority to mint or invalidate tokens.
- Logging must avoid exposing raw token contents; log identifiers (`jti`, `order_processing_id`) for traceability.

## Mapping to Repository Structure
- Model domain services under `internal/issuer`, `internal/relay`, `internal/consumer` (names to align with AGENTS.md guidance and article roles).
- Place shared token helpers and JWT utilities in `pkg/tokens` (signing, verification, claim structs).
- Create binaries in `cmd/issuer`, `cmd/relay`, and `cmd/consumer` that mirror the article’s flow: issue token → relay verification → consume capability.
- Store sample fixtures and golden tokens under `testdata/` to support integration tests.

## Open Questions
- How will we simulate or mock key management (KMS, JWKS) for local demos?
    - We assume the use of AWS KMS, but will handle it via a mock implementation. Define a KMSService interface and create a mock that consistently returns the same values.
- Should we provide both JSON/HTTP APIs and CLI demos for issuing and redeeming tokens?
    - Implement it as a CLI.
- What persistence (if any) does the orchestrator need to track `jti` usage in the sample?
    - In-memory tracking is sufficient.
