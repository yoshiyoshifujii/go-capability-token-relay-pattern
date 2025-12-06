package service

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type paymentMethodProviderServiceImpl struct {
	next service.PaymentConfirmationNext
}

func NewPaymentMethodProviderService(next service.PaymentConfirmationNext) service.PaymentMethodProviderService {
	if next == "" {
		next = service.PaymentConfirmationNextProcessing
	}
	return &paymentMethodProviderServiceImpl{next: next}
}

func (p *paymentMethodProviderServiceImpl) ConfirmPaymentMethod(ctx context.Context, request service.PaymentConfirmationRequest) (service.PaymentConfirmationResult, error) {
	return service.PaymentConfirmationResult{NextStatus: p.next}, nil
}
