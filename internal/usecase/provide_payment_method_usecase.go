package usecase

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/repository"
)

type (
	ProvidePaymentMethodUseCaseInput struct {
		PaymentIntentID   domain.PaymentIntentID
		PaymentMethodType domain.PaymentMethodType
		Card              *domain.PaymentMethodCard
		PayPay            *domain.PaymentMethodPayPay
	}

	ProvidePaymentMethodUseCaseOutput struct {
		PaymentIntent domain.PaymentIntent
	}

	ProvidePaymentMethodUseCase interface {
		Execute(context.Context, ProvidePaymentMethodUseCaseInput) (*ProvidePaymentMethodUseCaseOutput, error)
	}

	providePaymentMethodUseCase struct {
		paymentIntentRepository repository.PaymentIntentRepository
	}
)

func NewProvidePaymentMethodUseCase(paymentIntentRepository repository.PaymentIntentRepository) ProvidePaymentMethodUseCase {
	return &providePaymentMethodUseCase{
		paymentIntentRepository: paymentIntentRepository,
	}
}

func (u *providePaymentMethodUseCase) Execute(ctx context.Context, input ProvidePaymentMethodUseCaseInput) (*ProvidePaymentMethodUseCaseOutput, error) {
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

	method := domain.PaymentMethod{
		PaymentMethodType: input.PaymentMethodType,
		Card:              input.Card,
		PayPay:            input.PayPay,
	}

	event, aggregate, err := (*paymentIntent).RequireConfirmation(method)
	if err != nil {
		return nil, err
	}

	if err := u.paymentIntentRepository.Save(ctx, event, aggregate); err != nil {
		return nil, err
	}

	return &ProvidePaymentMethodUseCaseOutput{
		PaymentIntent: aggregate,
	}, nil
}
