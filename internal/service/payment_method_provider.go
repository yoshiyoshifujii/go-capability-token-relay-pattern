package service

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type (
	PaymentConfirmationRequest struct {
		Intent domain.PaymentIntentRequiresConfirmation
		Amount domain.Money
	}

	PaymentConfirmationResult struct {
		NextStatus domain.PaymentConfirmationNext
	}

	PaymentCaptureRequest struct {
		Intent domain.PaymentIntentRequiresCapture
		Amount domain.Money
	}

	PaymentMethodProviderService interface {
		ConfirmPaymentMethod(context.Context, PaymentConfirmationRequest) (PaymentConfirmationResult, error)
		CapturePaymentIntent(context.Context, PaymentCaptureRequest) error
	}
)
