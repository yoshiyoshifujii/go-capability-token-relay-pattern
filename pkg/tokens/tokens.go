package tokens

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	// TypeJWT denotes the token type value used in headers.
	TypeJWT = "JWT"
	// AlgHMACSHA256 is the signing algorithm used by the mock KMS service.
	AlgHMACSHA256 = "HS256"
)

// ErrInvalidToken indicates the token structure or signature could not be validated.
var ErrInvalidToken = errors.New("tokens: invalid token")

// KeyService abstracts the signing and signature verification operations.
type KeyService interface {
	Sign(ctx context.Context, kid string, data []byte) ([]byte, error)
	Verify(ctx context.Context, kid string, data, signature []byte) error
}

// Header represents the JWS header section.
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Kid string `json:"kid"`
}

// ConfirmedClaims mirrors the capability token contract described in the Zenn article.
type ConfirmedClaims struct {
	Iss               string            `json:"iss"`
	Aud               string            `json:"aud"`
	Sub               string            `json:"sub"`
	IssuedAt          int64             `json:"iat"`
	ExpiresAt         int64             `json:"exp"`
	TokenID           string            `json:"jti"`
	OrderProcessingID string            `json:"order_processing_id"`
	Capabilities      []string          `json:"capabilities"`
	Constraints       map[string]string `json:"constraints,omitempty"`
	RawPayload        json.RawMessage   `json:"-"`
}

// SignedToken bundles the header, claims, and raw representation of a confirmed token.
type SignedToken struct {
	Raw     string
	Header  Header
	Claims  ConfirmedClaims
	Payload []byte
}

// Validate performs inexpensive structural checks against the claims.
func (c ConfirmedClaims) Validate(now time.Time) error {
	switch {
	case c.Iss == "":
		return errors.New("tokens: iss is required")
	case c.Aud == "":
		return errors.New("tokens: aud is required")
	case c.Sub == "":
		return errors.New("tokens: sub is required")
	case c.TokenID == "":
		return errors.New("tokens: jti is required")
	case c.OrderProcessingID == "":
		return errors.New("tokens: order_processing_id is required")
	case c.ExpiresAt == 0:
		return errors.New("tokens: exp is required")
	}

	if c.IssuedAt > 0 && now.Unix() < c.IssuedAt-60 {
		return fmt.Errorf("tokens: iat (%d) is in the future", c.IssuedAt)
	}
	if now.Unix() >= c.ExpiresAt {
		return fmt.Errorf("tokens: token expired at %d", c.ExpiresAt)
	}
	if len(c.Capabilities) == 0 {
		return errors.New("tokens: capabilities must contain at least one entry")
	}
	return nil
}

// Encode marshals and signs the token claims using the supplied key service.
func Encode(ctx context.Context, svc KeyService, header Header, claims ConfirmedClaims) (SignedToken, error) {
	if header.Alg == "" {
		header.Alg = AlgHMACSHA256
	}
	if header.Typ == "" {
		header.Typ = TypeJWT
	}
	if header.Kid == "" {
		return SignedToken{}, errors.New("tokens: kid is required")
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return SignedToken{}, fmt.Errorf("tokens: marshal header: %w", err)
	}
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return SignedToken{}, fmt.Errorf("tokens: marshal claims: %w", err)
	}

	headerSegment := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadSegment := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signingInput := strings.Join([]string{headerSegment, payloadSegment}, ".")

	signature, err := svc.Sign(ctx, header.Kid, []byte(signingInput))
	if err != nil {
		return SignedToken{}, fmt.Errorf("tokens: sign: %w", err)
	}

	signatureSegment := base64.RawURLEncoding.EncodeToString(signature)
	raw := strings.Join([]string{signingInput, signatureSegment}, ".")

	claims.RawPayload = payloadJSON

	return SignedToken{
		Raw:     raw,
		Header:  header,
		Claims:  claims,
		Payload: payloadJSON,
	}, nil
}

// Decode parses and verifies a token string using the provided key service.
func Decode(ctx context.Context, svc KeyService, token string) (SignedToken, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return SignedToken{}, ErrInvalidToken
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return SignedToken{}, fmt.Errorf("tokens: decode header: %w", err)
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return SignedToken{}, fmt.Errorf("tokens: decode payload: %w", err)
	}
	signatureBytes, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return SignedToken{}, fmt.Errorf("tokens: decode signature: %w", err)
	}

	var header Header
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return SignedToken{}, fmt.Errorf("tokens: unmarshal header: %w", err)
	}
	if header.Kid == "" {
		return SignedToken{}, ErrInvalidToken
	}

	if err := svc.Verify(ctx, header.Kid, []byte(strings.Join(parts[:2], ".")), signatureBytes); err != nil {
		return SignedToken{}, fmt.Errorf("tokens: verify signature: %w", err)
	}

	var claims ConfirmedClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return SignedToken{}, fmt.Errorf("tokens: unmarshal claims: %w", err)
	}
	claims.RawPayload = payloadBytes

	return SignedToken{
		Raw:     token,
		Header:  header,
		Claims:  claims,
		Payload: payloadBytes,
	}, nil
}
