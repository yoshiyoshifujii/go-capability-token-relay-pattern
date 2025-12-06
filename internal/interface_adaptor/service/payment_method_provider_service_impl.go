package service

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type paymentMethodProviderServiceImpl struct {
	next domain.PaymentConfirmationNext
}

func NewPaymentMethodProviderService(next domain.PaymentConfirmationNext) service.PaymentMethodProviderService {
	if next == "" {
		next = domain.PaymentConfirmationNextProcessing
	}
	return &paymentMethodProviderServiceImpl{next: next}
}

func (p *paymentMethodProviderServiceImpl) ConfirmPaymentMethod(ctx context.Context, request service.PaymentConfirmationRequest) (service.PaymentConfirmationResult, error) {
	return service.PaymentConfirmationResult{NextStatus: p.next}, nil
}

func (p *paymentMethodProviderServiceImpl) CapturePaymentIntent(ctx context.Context, request service.PaymentCaptureRequest) error {
	return nil
}
