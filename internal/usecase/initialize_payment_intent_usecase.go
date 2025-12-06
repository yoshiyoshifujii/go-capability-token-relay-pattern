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
		CartToken service.SignedToken
	}

	InitializePaymentIntentUseCaseOutput struct {
		PaymentIntentID domain.PaymentIntentID
		PaymentIntent   domain.PaymentIntent
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
	return &initializePaymentIntentUseCase{
		tokenService:            tokenService,
		paymentIntentRepository: paymentIntentRepository,
		paymentIntentGenerator:  paymentIntentGenerator,
		businessRepository:      businessRepository,
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
	if u.businessRepository == nil {
		return nil, errors.New("businessRepository is nil")
	}

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
		PaymentIntentID: paymentIntentID,
		PaymentIntent:   aggregate,
	}, nil
}
