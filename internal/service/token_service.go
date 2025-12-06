package service

import "context"

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
)
