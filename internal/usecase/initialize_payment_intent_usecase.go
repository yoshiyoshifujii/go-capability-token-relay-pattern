package usecase

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/repository"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type (
	InitializePaymentIntentUseCaseInput struct {
		CartToken          service.SignedToken
		PaymentMethodTypes domain.PaymentMethodTypes
	}

	InitializePaymentIntentUseCaseOutput struct {
		PaymentIntent domain.PaymentIntent
	}

	InitializePaymentIntentUseCase interface {
		Execute(context.Context, InitializePaymentIntentUseCaseInput) (*InitializePaymentIntentUseCaseOutput, error)
	}

	initializePaymentIntentUseCase struct {
		tokenService            service.TokenService
		paymentIntentRepository repository.PaymentIntentRepository
		paymentIntentGenerator  service.PaymentIDGenerator
	}
)

func NewInitializePaymentIntentUseCase(
	tokenService service.TokenService,
	paymentIntentRepository repository.PaymentIntentRepository,
	paymentIntentGenerator service.PaymentIDGenerator,
) InitializePaymentIntentUseCase {
	return &initializePaymentIntentUseCase{
		tokenService:            tokenService,
		paymentIntentRepository: paymentIntentRepository,
		paymentIntentGenerator:  paymentIntentGenerator,
	}
}

func (u *initializePaymentIntentUseCase) Execute(ctx context.Context, input InitializePaymentIntentUseCaseInput) (*InitializePaymentIntentUseCaseOutput, error) {
	if u == nil {
		panic("usecase is nil")
	}
	if u.tokenService == nil {
		return nil, errors.New("tokenService is nil")
	}
	if u.paymentIntentRepository == nil {
		return nil, errors.New("paymentIntentRepository is nil")
	}
	if u.paymentIntentGenerator == nil {
		return nil, errors.New("paymentIntentGenerator is nil")
	}

	if err := u.tokenService.ValidateCartToken(ctx, input.CartToken); err != nil {
		return nil, err
	}

	paymentIntentID, err := u.paymentIntentGenerator.GenerateID(ctx)
	if err != nil {
		return nil, err
	}

	event, aggregate, err := domain.GeneratePaymentIntent(paymentIntentID, input.PaymentMethodTypes)
	if err != nil {
		return nil, err
	}

	if err := u.paymentIntentRepository.Save(ctx, event, aggregate); err != nil {
		return nil, err
	}

	return &InitializePaymentIntentUseCaseOutput{
		PaymentIntent: aggregate,
	}, nil
}
