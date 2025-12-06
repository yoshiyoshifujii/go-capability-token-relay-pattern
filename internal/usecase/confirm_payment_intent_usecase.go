package usecase

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/repository"
)

type (
	ConfirmPaymentIntentUseCaseInput struct {
		PaymentIntentID domain.PaymentIntentID
		CaptureMethod   domain.PaymentCaptureMethod
	}

	ConfirmPaymentIntentUseCaseOutput struct {
		PaymentIntentID domain.PaymentIntentID
		PaymentIntent   domain.PaymentIntent
	}

	ConfirmPaymentIntentUseCase interface {
		Execute(context.Context, ConfirmPaymentIntentUseCaseInput) (*ConfirmPaymentIntentUseCaseOutput, error)
	}

	confirmPaymentIntentUseCase struct {
		paymentIntentRepository repository.PaymentIntentRepository
	}
)

func NewConfirmPaymentIntentUseCase(paymentIntentRepository repository.PaymentIntentRepository) ConfirmPaymentIntentUseCase {
	if paymentIntentRepository == nil {
		panic("paymentIntentRepository is nil")
	}
	return &confirmPaymentIntentUseCase{
		paymentIntentRepository: paymentIntentRepository,
	}
}

func (i ConfirmPaymentIntentUseCaseInput) Validate() error {
	contract.AssertValidatable(i.PaymentIntentID)
	contract.AssertValidatable(i.CaptureMethod)
	return nil
}

func (u *confirmPaymentIntentUseCase) Execute(ctx context.Context, input ConfirmPaymentIntentUseCaseInput) (*ConfirmPaymentIntentUseCaseOutput, error) {
	contract.AssertValidatable(input)

	paymentIntent, err := u.paymentIntentRepository.FindBy(ctx, input.PaymentIntentID)
	if err != nil {
		return nil, err
	}
	if paymentIntent == nil {
		return nil, errors.New("payment intent not found")
	}

	var (
		event     domain.PaymentIntentEvent
		aggregate domain.PaymentIntent
	)

	switch input.CaptureMethod {
	case domain.PaymentCaptureMethodAutomatic:
		event, aggregate, err = (*paymentIntent).StartProcessing()
	case domain.PaymentCaptureMethodManual:
		event, aggregate, err = (*paymentIntent).RequireCapture()
	default:
		panic("invalid capture method")
	}

	if err != nil {
		return nil, err
	}

	if err := u.paymentIntentRepository.Save(ctx, event, aggregate); err != nil {
		return nil, err
	}

	return &ConfirmPaymentIntentUseCaseOutput{
		PaymentIntentID: input.PaymentIntentID,
		PaymentIntent:   aggregate,
	}, nil
}
