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
	InitializePaymentIntentUseCaseInput struct {
		CartToken service.SignedToken
	}

	InitializePaymentIntentUseCaseOutput struct {
		PaymentIntentID    domain.PaymentIntentID
		PaymentIntent      domain.PaymentIntent
		PaymentMethodTypes domain.PaymentMethodTypes
	}

	InitializePaymentIntentUseCase interface {
		Execute(context.Context, InitializePaymentIntentUseCaseInput) (*InitializePaymentIntentUseCaseOutput, error)
	}

	initializePaymentIntentUseCase struct {
		tokenService            service.TokenService
		paymentIntentRepository repository.PaymentIntentRepository
		paymentIntentGenerator  service.PaymentIDGenerator
		businessRepository      repository.BusinessRepository
	}
)

func NewInitializePaymentIntentUseCase(
	tokenService service.TokenService,
	paymentIntentRepository repository.PaymentIntentRepository,
	paymentIntentGenerator service.PaymentIDGenerator,
	businessRepository repository.BusinessRepository,
) InitializePaymentIntentUseCase {
	if tokenService == nil {
		panic("tokenService is nil")
	}
	if paymentIntentRepository == nil {
		panic("paymentIntentRepository is nil")
	}
	if paymentIntentGenerator == nil {
		panic("paymentIntentGenerator is nil")
	}
	if businessRepository == nil {
		panic("businessRepository is nil")
	}
	return &initializePaymentIntentUseCase{
		tokenService:            tokenService,
		paymentIntentRepository: paymentIntentRepository,
		paymentIntentGenerator:  paymentIntentGenerator,
		businessRepository:      businessRepository,
	}
}

func (i InitializePaymentIntentUseCaseInput) Validate() error {
	contract.AssertValidatable(i.CartToken)
	return nil
}

func (u *initializePaymentIntentUseCase) Execute(ctx context.Context, input InitializePaymentIntentUseCaseInput) (*InitializePaymentIntentUseCaseOutput, error) {
	contract.AssertValidatable(input)

	cart, err := u.tokenService.ParseCartToken(ctx, input.CartToken)
	if err != nil {
		return nil, err
	}

	business, err := u.businessRepository.FindBy(ctx, cart.BusinessID)
	if err != nil {
		return nil, err
	}
	if business == nil {
		return nil, errors.New("business not found")
	}

	paymentIntentID, err := u.paymentIntentGenerator.GenerateID(ctx)
	if err != nil {
		return nil, err
	}

	event, aggregate, err := domain.GeneratePaymentIntent(paymentIntentID, business.PaymentMethodTypes)
	if err != nil {
		return nil, err
	}

	if err := u.paymentIntentRepository.Save(ctx, event, aggregate); err != nil {
		return nil, err
	}

	return &InitializePaymentIntentUseCaseOutput{
		PaymentIntentID:    paymentIntentID,
		PaymentIntent:      aggregate,
		PaymentMethodTypes: business.PaymentMethodTypes,
	}, nil
}
