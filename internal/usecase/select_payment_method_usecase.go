package usecase

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/repository"
)

type (
	SelectPaymentMethodUseCaseInput struct {
		PaymentIntentID   domain.PaymentIntentID
		PaymentMethodType domain.PaymentMethodType
	}

	SelectPaymentMethodUseCaseOutput struct {
		PaymentIntent domain.PaymentIntent
	}

	SelectPaymentMethodUseCase interface {
		Execute(context.Context, SelectPaymentMethodUseCaseInput) (*SelectPaymentMethodUseCaseOutput, error)
	}

	selectPaymentMethodUseCase struct {
		paymentIntentRepository repository.PaymentIntentRepository
	}
)

func NewSelectPaymentMethodUseCase(paymentIntentRepository repository.PaymentIntentRepository) SelectPaymentMethodUseCase {
	return &selectPaymentMethodUseCase{
		paymentIntentRepository: paymentIntentRepository,
	}
}

func (u *selectPaymentMethodUseCase) Execute(ctx context.Context, input SelectPaymentMethodUseCaseInput) (*SelectPaymentMethodUseCaseOutput, error) {
	if u == nil {
		panic("usecase is nil")
	}
	if u.paymentIntentRepository == nil {
		return nil, errors.New("paymentIntentRepository is nil")
	}

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
		PaymentIntent: aggregate,
	}, nil
}
