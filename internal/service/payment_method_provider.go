package service

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type (
	PaymentConfirmationRequest struct {
		Intent domain.PaymentIntentRequiresConfirmation
	}

	PaymentConfirmationResult struct {
		NextStatus domain.PaymentConfirmationNext
	}

	PaymentMethodProviderService interface {
		ConfirmPaymentMethod(context.Context, PaymentConfirmationRequest) (PaymentConfirmationResult, error)
	}
)
