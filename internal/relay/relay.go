package relay

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"yoshiyoshifujii/go-capability-token-relay-pattern/pkg/tokens"
)

// Clock abstracts time.Now for testing.
type Clock func() time.Time

// Service verifies confirmed tokens and relays them to downstream consumers.
type Service struct {
	keySvc tokens.KeyService
	clock  Clock

	mu   sync.Mutex
	seen map[string]struct{}
}

// NewService constructs a relay service instance.
func NewService(keySvc tokens.KeyService, clk Clock) (*Service, error) {
	if keySvc == nil {
		return nil, errors.New("relay: key service is required")
	}
	if clk == nil {
		clk = time.Now
	}
	return &Service{
		keySvc: keySvc,
		clock:  clk,
		seen:   make(map[string]struct{}),
	}, nil
}

// Verify decodes a token, validates its claims, and tracks single-use identifiers.
func (s *Service) Verify(ctx context.Context, raw string) (tokens.SignedToken, error) {
	token, err := tokens.Decode(ctx, s.keySvc, raw)
	if err != nil {
		return tokens.SignedToken{}, fmt.Errorf("relay: decode token: %w", err)
	}

	if err := token.Claims.Validate(s.clock()); err != nil {
		return tokens.SignedToken{}, fmt.Errorf("relay: %w", err)
	}

	if err := s.track(token.Claims.OrderProcessingID, token.Claims.TokenID); err != nil {
		return tokens.SignedToken{}, err
	}

	return token, nil
}

// Reset clears the seen-token cache; primarily useful for tests or CLI demos.
func (s *Service) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen = make(map[string]struct{})
}

func (s *Service) track(orderID, tokenID string) error {
	if orderID == "" || tokenID == "" {
		return errors.New("relay: token missing order id or token id")
	}
	key := orderID + ":" + tokenID

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.seen[key]; exists {
		return fmt.Errorf("relay: token %s already used for order %s", tokenID, orderID)
	}
	s.seen[key] = struct{}{}
	return nil
}
