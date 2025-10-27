package kmsmock

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"

	"yoshiyoshifujii/go-capability-token-relay-pattern/pkg/tokens"
)

var errUnknownKey = errors.New("kmsmock: unknown key id")

// Service emulates a subset of KMS signing/verification using HMAC-SHA256.
type Service struct {
	mu   sync.RWMutex
	keys map[string][]byte
}

// NewService constructs a mock KMS service pre-loaded with a key.
func NewService(keyID string, secret []byte) *Service {
	svc := &Service{keys: make(map[string][]byte)}
	if keyID != "" && len(secret) > 0 {
		svc.keys[keyID] = append([]byte(nil), secret...)
	}
	return svc
}

// RegisterKey adds or replaces a signing key for the supplied key identifier.
func (s *Service) RegisterKey(keyID string, secret []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.keys == nil {
		s.keys = make(map[string][]byte)
	}
	s.keys[keyID] = append([]byte(nil), secret...)
}

// Sign satisfies tokens.KeyService using HMAC-SHA256.
func (s *Service) Sign(_ context.Context, kid string, data []byte) ([]byte, error) {
	key, err := s.lookupKey(kid)
	if err != nil {
		return nil, err
	}

	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write(data); err != nil {
		return nil, fmt.Errorf("kmsmock: sign: %w", err)
	}
	return mac.Sum(nil), nil
}

// Verify satisfies tokens.KeyService by checking the HMAC signature.
func (s *Service) Verify(_ context.Context, kid string, data, signature []byte) error {
	key, err := s.lookupKey(kid)
	if err != nil {
		return err
	}

	expectedMAC := hmac.New(sha256.New, key)
	if _, err := expectedMAC.Write(data); err != nil {
		return fmt.Errorf("kmsmock: verify: %w", err)
	}

	if !hmac.Equal(expectedMAC.Sum(nil), signature) {
		return errors.New("kmsmock: signature mismatch")
	}
	return nil
}

func (s *Service) lookupKey(kid string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.keys) == 0 {
		return nil, errUnknownKey
	}
	key, ok := s.keys[kid]
	if !ok {
		return nil, errUnknownKey
	}
	return key, nil
}

var _ tokens.KeyService = (*Service)(nil)
