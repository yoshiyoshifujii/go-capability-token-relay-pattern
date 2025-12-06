package service

import (
	"context"
	"fmt"
)

type (
	IssueCouponTokenInput struct {
		OrderProcessingID string
		UserID            string
		CouponRef         string
	}

	IssuePointsTokenInput struct {
		OrderProcessingID string
		UserID            string
		PointsToUse       int64
	}

	IssuePaymentTokenInput struct {
		OrderProcessingID string
		UserID            string
		PaymentMethod     string
	}

	RelayTokensInput struct {
		OrderProcessingID string
		CouponToken       string
		PointsToken       string
		PaymentToken      string
	}

	TokenService interface {
		IssueCouponToken(context.Context, IssueCouponTokenInput) (string, error)
		IssuePointsToken(context.Context, IssuePointsTokenInput) (string, error)
		IssuePaymentToken(context.Context, IssuePaymentTokenInput) (string, error)
		RelayTokens(context.Context, RelayTokensInput) (map[string]string, error)
	}

	tokenService struct{}
)

func NewTokenService() TokenService {
	return &tokenService{}
}

func (s *tokenService) IssueCouponToken(ctx context.Context, input IssueCouponTokenInput) (string, error) {
	return fmt.Sprintf("coupon-token:%s:%s", input.OrderProcessingID, input.CouponRef), nil
}

func (s *tokenService) IssuePointsToken(ctx context.Context, input IssuePointsTokenInput) (string, error) {
	return fmt.Sprintf("points-token:%s:%d", input.OrderProcessingID, input.PointsToUse), nil
}

func (s *tokenService) IssuePaymentToken(ctx context.Context, input IssuePaymentTokenInput) (string, error) {
	return fmt.Sprintf("payment-token:%s:%s", input.OrderProcessingID, input.PaymentMethod), nil
}

func (s *tokenService) RelayTokens(ctx context.Context, input RelayTokensInput) (map[string]string, error) {
	return map[string]string{
		"coupon":  input.CouponToken,
		"points":  input.PointsToken,
		"payment": input.PaymentToken,
	}, nil
}
