#!/usr/bin/env bash
set -euo pipefail

# Demonstrates issuing a confirmed token, relaying it, and consuming it locally.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

export CTR_KMS_KEY_ID="${CTR_KMS_KEY_ID:-mock-key}"
export CTR_KMS_SECRET="${CTR_KMS_SECRET:-local-dev-secret}"

token="$(go run "$ROOT_DIR/cmd/issuer" \
  --subject "user_123" \
  --order-id "op_demo" \
  --capabilities "coupons:redeem" \
  --constraints "coupon_ref=ABC123")"

printf '%s\n' "$token" | go run "$ROOT_DIR/cmd/relay" | \
  go run "$ROOT_DIR/cmd/consumer" \
    --capability "coupons:redeem" \
    --require "coupon_ref"
