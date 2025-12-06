package service

import (
	"context"
	"fmt"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type tokenServiceImpl struct{}

func NewTokenService() service.TokenService {
	return &tokenServiceImpl{}
}

func (s *tokenServiceImpl) IssuePaymentToken(ctx context.Context, input service.IssuePaymentTokenInput) (service.SignedToken, error) {
	return service.SignedToken{Value: fmt.Sprintf("payment-token:%s:%s", input.OrderProcessingID, input.PaymentMethod)}, nil
}

func (s *tokenServiceImpl) RelayTokens(ctx context.Context, input service.RelayTokensInput) (map[string]service.SignedToken, error) {
	return map[string]service.SignedToken{
		"cart":    {Value: input.CartToken},
		"payment": {Value: input.PaymentToken},
	}, nil
}

func (s *tokenServiceImpl) ConfirmCartToken(ctx context.Context, input service.ConfirmCartTokenInput) (service.SignedToken, error) {
	return service.SignedToken{
		Value: fmt.Sprintf(
			"cart-token:%s:%s:%d",
			input.Cart.BusinessID,
			input.Cart.CartID,
			len(input.Cart.Items),
		),
	}, nil
}
