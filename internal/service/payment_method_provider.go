package service

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type (
	PaymentConfirmationNext string

	PaymentConfirmationRequest struct {
		Intent domain.PaymentIntentRequiresConfirmation
	}

	PaymentConfirmationResult struct {
		NextStatus PaymentConfirmationNext
	}

	PaymentMethodProviderService interface {
		ConfirmPaymentMethod(context.Context, PaymentConfirmationRequest) (PaymentConfirmationResult, error)
	}
)

const (
	PaymentConfirmationNextProcessing      PaymentConfirmationNext = "processing"
	PaymentConfirmationNextRequiresAction  PaymentConfirmationNext = "requires_action"
	PaymentConfirmationNextRequiresCapture PaymentConfirmationNext = "requires_capture"
)
