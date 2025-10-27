package issuer

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"yoshiyoshifujii/go-capability-token-relay-pattern/pkg/tokens"
)

// Clock abstracts time.Now for easier testing.
type Clock func() time.Time

// Options configures the issuer service.
type Options struct {
	Issuer     string
	Audience   string
	KeyID      string
	DefaultTTL time.Duration
	Clock      Clock
}

// IssueRequest describes the inputs required to mint a confirmed token.
type IssueRequest struct {
	Subject           string
	OrderProcessingID string
	Capabilities      []string
	Constraints       map[string]string
	TTL               time.Duration
	TokenID           string
}

// Service issues confirmed capability tokens based on validated user input.
type Service struct {
	keySvc tokens.KeyService
	opts   Options
	clock  Clock
}

// NewService builds an issuer service with the supplied configuration.
func NewService(keySvc tokens.KeyService, opts Options) (*Service, error) {
	if keySvc == nil {
		return nil, errors.New("issuer: key service is required")
	}
	if opts.Issuer == "" {
		return nil, errors.New("issuer: issuer is required")
	}
	if opts.Audience == "" {
		return nil, errors.New("issuer: audience is required")
	}
	if opts.KeyID == "" {
		return nil, errors.New("issuer: key id is required")
	}
	if opts.DefaultTTL <= 0 {
		opts.DefaultTTL = 5 * time.Minute
	}
	clk := opts.Clock
	if clk == nil {
		clk = time.Now
	}

	return &Service{
		keySvc: keySvc,
		opts:   opts,
		clock:  clk,
	}, nil
}

// IssueConfirmedToken validates the request and signs a confirmed token.
func (s *Service) IssueConfirmedToken(ctx context.Context, req IssueRequest) (tokens.SignedToken, error) {
	if err := validateRequest(req); err != nil {
		return tokens.SignedToken{}, err
	}

	tokenID := req.TokenID
	if tokenID == "" {
		id, err := randomID("cpt")
		if err != nil {
			return tokens.SignedToken{}, fmt.Errorf("issuer: generate token id: %w", err)
		}
		tokenID = id
	}

	ttl := req.TTL
	if ttl <= 0 {
		ttl = s.opts.DefaultTTL
	}

	now := s.clock()
	claims := tokens.ConfirmedClaims{
		Iss:               s.opts.Issuer,
		Aud:               s.opts.Audience,
		Sub:               req.Subject,
		IssuedAt:          now.Unix(),
		ExpiresAt:         now.Add(ttl).Unix(),
		TokenID:           tokenID,
		OrderProcessingID: req.OrderProcessingID,
		Capabilities:      append([]string(nil), req.Capabilities...),
		Constraints:       cloneMap(req.Constraints),
	}

	if err := claims.Validate(now); err != nil {
		return tokens.SignedToken{}, fmt.Errorf("issuer: %w", err)
	}

	header := tokens.Header{
		Alg: tokens.AlgHMACSHA256,
		Typ: tokens.TypeJWT,
		Kid: s.opts.KeyID,
	}

	token, err := tokens.Encode(ctx, s.keySvc, header, claims)
	if err != nil {
		return tokens.SignedToken{}, fmt.Errorf("issuer: encode token: %w", err)
	}
	return token, nil
}

func validateRequest(req IssueRequest) error {
	switch {
	case req.Subject == "":
		return errors.New("issuer: subject is required")
	case req.OrderProcessingID == "":
		return errors.New("issuer: order processing id is required")
	case len(req.Capabilities) == 0:
		return errors.New("issuer: at least one capability is required")
	}
	return nil
}

func randomID(prefix string) (string, error) {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(buf[:])), nil
}

func cloneMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
