package usecase

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/repository"
)

type (
	SelectPaymentMethodUseCaseInput struct {
		PaymentIntentID   domain.PaymentIntentID
		PaymentMethodType domain.PaymentMethodType
	}

	SelectPaymentMethodUseCaseOutput struct {
		PaymentIntentID domain.PaymentIntentID
		PaymentIntent   domain.PaymentIntent
	}

	SelectPaymentMethodUseCase interface {
		Execute(context.Context, SelectPaymentMethodUseCaseInput) (*SelectPaymentMethodUseCaseOutput, error)
	}

	selectPaymentMethodUseCase struct {
		paymentIntentRepository repository.PaymentIntentRepository
	}
)

func NewSelectPaymentMethodUseCase(paymentIntentRepository repository.PaymentIntentRepository) SelectPaymentMethodUseCase {
	if paymentIntentRepository == nil {
		panic("paymentIntentRepository is nil")
	}
	return &selectPaymentMethodUseCase{
		paymentIntentRepository: paymentIntentRepository,
	}
}

func (i SelectPaymentMethodUseCaseInput) Validate() error {
	contract.AssertValidatable(i.PaymentIntentID)
	contract.AssertValidatable(i.PaymentMethodType)
	return nil
}

func (u *selectPaymentMethodUseCase) Execute(ctx context.Context, input SelectPaymentMethodUseCaseInput) (*SelectPaymentMethodUseCaseOutput, error) {
	contract.AssertValidatable(input)

	paymentIntent, err := u.paymentIntentRepository.FindBy(ctx, input.PaymentIntentID)
	if err != nil {
		return nil, err
	}
	if paymentIntent == nil {
		return nil, errors.New("payment intent not found")
	}

	event, aggregate, err := (*paymentIntent).RequirePaymentMethod(input.PaymentMethodType)
	if err != nil {
		return nil, err
	}

	if err := u.paymentIntentRepository.Save(ctx, event, aggregate); err != nil {
		return nil, err
	}

	return &SelectPaymentMethodUseCaseOutput{
		PaymentIntentID: input.PaymentIntentID,
		PaymentIntent:   aggregate,
	}, nil
}
