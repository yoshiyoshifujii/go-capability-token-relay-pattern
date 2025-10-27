package consumer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"yoshiyoshifujii/go-capability-token-relay-pattern/pkg/tokens"
)

// Options configures the consumer service.
type Options struct {
	Domain string
}

// ConsumeInput represents the parameters required to process a token.
type ConsumeInput struct {
	Token               tokens.SignedToken
	Capability          string
	RequiredConstraints []string
}

// ConsumeResult reports the simulated processing outcome.
type ConsumeResult struct {
	Domain            string
	OrderProcessingID string
	TokenID           string
	ProcessedAt       time.Time
	Constraints       map[string]string
}

// Service consumes confirmed tokens and performs domain-specific actions.
type Service struct {
	domain string
}

// NewService constructs a domain consumer with the supplied options.
func NewService(opts Options) *Service {
	domain := opts.Domain
	if domain == "" {
		domain = "consumer"
	}
	return &Service{domain: domain}
}

// Consume validates the token against expected capabilities and constraints.
func (s *Service) Consume(_ context.Context, input ConsumeInput) (ConsumeResult, error) {
	if input.Capability == "" {
		return ConsumeResult{}, errors.New("consumer: capability is required")
	}
	if len(input.Token.Claims.Capabilities) == 0 {
		return ConsumeResult{}, errors.New("consumer: token has no capabilities")
	}
	if !hasCapability(input.Token.Claims.Capabilities, input.Capability) {
		return ConsumeResult{}, fmt.Errorf("consumer: capability %s not present in token", input.Capability)
	}

	for _, key := range input.RequiredConstraints {
		if key == "" {
			continue
		}
		if val := input.Token.Claims.Constraints[key]; val == "" {
			return ConsumeResult{}, fmt.Errorf("consumer: constraint %s missing or empty", key)
		}
	}

	return ConsumeResult{
		Domain:            s.domain,
		OrderProcessingID: input.Token.Claims.OrderProcessingID,
		TokenID:           input.Token.Claims.TokenID,
		ProcessedAt:       time.Now(),
		Constraints:       cloneConstraints(input.Token.Claims.Constraints),
	}, nil
}

func hasCapability(list []string, capability string) bool {
	for _, item := range list {
		if item == capability {
			return true
		}
	}
	return false
}

func cloneConstraints(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
