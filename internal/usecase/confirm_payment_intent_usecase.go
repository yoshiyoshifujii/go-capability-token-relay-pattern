package usecase

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/repository"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type (
	ConfirmPaymentIntentUseCaseInput struct {
		PaymentIntentID domain.PaymentIntentID
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
		paymentProvider         service.PaymentMethodProviderService
	}
)

func NewConfirmPaymentIntentUseCase(
	paymentIntentRepository repository.PaymentIntentRepository,
	paymentProvider service.PaymentMethodProviderService,
) ConfirmPaymentIntentUseCase {
	if paymentIntentRepository == nil {
		panic("paymentIntentRepository is nil")
	}
	if paymentProvider == nil {
		panic("paymentProvider is nil")
	}
	return &confirmPaymentIntentUseCase{
		paymentIntentRepository: paymentIntentRepository,
		paymentProvider:         paymentProvider,
	}
}

func (i ConfirmPaymentIntentUseCaseInput) Validate() error {
	contract.AssertValidatable(i.PaymentIntentID)
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

	intent, ok := (*paymentIntent).(domain.PaymentIntentRequiresConfirmation)
	if !ok {
		return nil, errors.New("payment intent not ready for confirmation")
	}

	result, err := u.paymentProvider.ConfirmPaymentMethod(ctx, service.PaymentConfirmationRequest{
		Intent: intent,
	})
	if err != nil {
		return nil, err
	}

	event, aggregate, err := intent.ApplyConfirmationResult(result.NextStatus)
	if err != nil {
		return nil, err
	}

	if err = u.paymentIntentRepository.Save(ctx, event, aggregate); err != nil {
		return nil, err
	}

	return &ConfirmPaymentIntentUseCaseOutput{
		PaymentIntentID: input.PaymentIntentID,
		PaymentIntent:   aggregate,
	}, nil
}
