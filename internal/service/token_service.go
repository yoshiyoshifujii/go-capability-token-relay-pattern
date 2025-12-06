package service

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type (
	SignedToken struct {
		Value string
	}

	IssuePaymentTokenInput struct {
		OrderProcessingID string
		UserID            string
		PaymentMethod     string
	}

	RelayTokensInput struct {
		OrderProcessingID string
		CartToken         string
		PaymentToken      string
	}

	ConfirmCartTokenInput struct {
		Cart domain.Cart
	}

	TokenService interface {
		IssuePaymentToken(context.Context, IssuePaymentTokenInput) (SignedToken, error)
		ConfirmCartToken(context.Context, ConfirmCartTokenInput) (SignedToken, error)
		ParseCartToken(context.Context, SignedToken) (domain.Cart, error)
		RelayTokens(context.Context, RelayTokensInput) (map[string]SignedToken, error)
	}
)

func (s SignedToken) Validate() error {
	if s.Value == "" {
		return errors.New("signed token is empty")
	}
	return nil
}
