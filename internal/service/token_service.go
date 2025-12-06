package service

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type (
	SignedToken struct {
		Value string
	}

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

	ConfirmCartTokenInput struct {
		Cart domain.Cart
	}

	TokenService interface {
		IssueCouponToken(context.Context, IssueCouponTokenInput) (SignedToken, error)
		IssuePointsToken(context.Context, IssuePointsTokenInput) (SignedToken, error)
		IssuePaymentToken(context.Context, IssuePaymentTokenInput) (SignedToken, error)
		RelayTokens(context.Context, RelayTokensInput) (map[string]SignedToken, error)
		ConfirmCartToken(context.Context, ConfirmCartTokenInput) (SignedToken, error)
	}
)
