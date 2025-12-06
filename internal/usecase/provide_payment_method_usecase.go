package usecase

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
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
		PaymentIntentID domain.PaymentIntentID
		PaymentIntent   domain.PaymentIntent
	}

	ProvidePaymentMethodUseCase interface {
		Execute(context.Context, ProvidePaymentMethodUseCaseInput) (*ProvidePaymentMethodUseCaseOutput, error)
	}

	providePaymentMethodUseCase struct {
		paymentIntentRepository repository.PaymentIntentRepository
	}
)

func NewProvidePaymentMethodUseCase(paymentIntentRepository repository.PaymentIntentRepository) ProvidePaymentMethodUseCase {
	if paymentIntentRepository == nil {
		panic("paymentIntentRepository is nil")
	}
	return &providePaymentMethodUseCase{
		paymentIntentRepository: paymentIntentRepository,
	}
}

func (i ProvidePaymentMethodUseCaseInput) Validate() error {
	contract.AssertValidatable(i.PaymentIntentID)
	if i.PaymentMethodType == "" {
		return errors.New("payment method type is empty")
	}
	switch i.PaymentMethodType {
	case domain.PaymentMethodTypeCard:
		if i.Card == nil {
			return errors.New("card details are required for card payment method")
		}
	case domain.PaymentMethodTypePayPay:
		if i.PayPay == nil {
			return errors.New("paypay details are required for paypay payment method")
		}
	default:
		return errors.New("unsupported payment method type")
	}
	return nil
}

func (u *providePaymentMethodUseCase) Execute(ctx context.Context, input ProvidePaymentMethodUseCaseInput) (*ProvidePaymentMethodUseCaseOutput, error) {
	contract.AssertValidatable(input)

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
		PaymentIntentID: input.PaymentIntentID,
		PaymentIntent:   aggregate,
	}, nil
}
