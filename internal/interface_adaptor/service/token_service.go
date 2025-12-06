package service

import (
	"context"
	"fmt"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type tokenService struct{}

func NewTokenService() service.TokenService {
	return &tokenService{}
}

func (s *tokenService) IssueCouponToken(ctx context.Context, input service.IssueCouponTokenInput) (service.SignedToken, error) {
	return service.SignedToken{Value: fmt.Sprintf("coupon-token:%s:%s", input.OrderProcessingID, input.CouponRef)}, nil
}

func (s *tokenService) IssuePointsToken(ctx context.Context, input service.IssuePointsTokenInput) (service.SignedToken, error) {
	return service.SignedToken{Value: fmt.Sprintf("points-token:%s:%d", input.OrderProcessingID, input.PointsToUse)}, nil
}

func (s *tokenService) IssuePaymentToken(ctx context.Context, input service.IssuePaymentTokenInput) (service.SignedToken, error) {
	return service.SignedToken{Value: fmt.Sprintf("payment-token:%s:%s", input.OrderProcessingID, input.PaymentMethod)}, nil
}

func (s *tokenService) RelayTokens(ctx context.Context, input service.RelayTokensInput) (map[string]service.SignedToken, error) {
	return map[string]service.SignedToken{
		"coupon":  {Value: input.CouponToken},
		"points":  {Value: input.PointsToken},
		"payment": {Value: input.PaymentToken},
	}, nil
}

func (s *tokenService) ConfirmCartToken(ctx context.Context, input service.ConfirmCartTokenInput) (service.SignedToken, error) {
	return service.SignedToken{
		Value: fmt.Sprintf(
			"cart-token:%s:%s:%d",
			input.Cart.BusinessID,
			input.Cart.CartID,
			len(input.Cart.Items),
		),
	}, nil
}
