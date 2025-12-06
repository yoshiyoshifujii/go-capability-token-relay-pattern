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

func (s *tokenService) IssueCouponToken(ctx context.Context, input service.IssueCouponTokenInput) (string, error) {
	return fmt.Sprintf("coupon-token:%s:%s", input.OrderProcessingID, input.CouponRef), nil
}

func (s *tokenService) IssuePointsToken(ctx context.Context, input service.IssuePointsTokenInput) (string, error) {
	return fmt.Sprintf("points-token:%s:%d", input.OrderProcessingID, input.PointsToUse), nil
}

func (s *tokenService) IssuePaymentToken(ctx context.Context, input service.IssuePaymentTokenInput) (string, error) {
	return fmt.Sprintf("payment-token:%s:%s", input.OrderProcessingID, input.PaymentMethod), nil
}

func (s *tokenService) RelayTokens(ctx context.Context, input service.RelayTokensInput) (map[string]string, error) {
	return map[string]string{
		"coupon":  input.CouponToken,
		"points":  input.PointsToken,
		"payment": input.PaymentToken,
	}, nil
}
