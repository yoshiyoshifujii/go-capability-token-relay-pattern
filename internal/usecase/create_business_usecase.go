package usecase

import (
	"context"
	"errors"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/repository"
)

type (
	CreateBusinessUseCaseInput struct {
		BusinessID         string
		Name               string
		PaymentMethodTypes domain.PaymentMethodTypes
	}

	CreateBusinessUseCaseOutput struct {
		Business domain.Business
	}

	CreateBusinessUseCase interface {
		Execute(context.Context, CreateBusinessUseCaseInput) (*CreateBusinessUseCaseOutput, error)
	}

	createBusinessUseCase struct {
		businessIDGenerator service.BusinessIDGenerator
		businessRepository  repository.BusinessRepository
	}
)

func NewCreateBusinessUseCase(
	businessIDGenerator service.BusinessIDGenerator,
	businessRepository repository.BusinessRepository,
) CreateBusinessUseCase {
	if businessIDGenerator == nil {
		panic("businessIDGenerator is nil")
	}
	if businessRepository == nil {
		panic("businessRepository is nil")
	}
	return &createBusinessUseCase{
		businessIDGenerator: businessIDGenerator,
		businessRepository:  businessRepository,
	}
}

func (i CreateBusinessUseCaseInput) Validate() error {
	if len(i.BusinessID) == 0 {
		return errors.New("business id is empty")
	}
	if len(i.Name) == 0 {
		return errors.New("business name is empty")
	}
	contract.AssertValidatable(i.PaymentMethodTypes)
	return nil
}

func (u *createBusinessUseCase) Execute(ctx context.Context, input CreateBusinessUseCaseInput) (*CreateBusinessUseCaseOutput, error) {
	contract.AssertValidatable(input)

	businessID, err := u.businessIDGenerator.GenerateID(ctx)
	if err != nil {
		return nil, err
	}

	businessEvent := domain.NewBusinessInitializedEvent(
		businessID,
		1,
		input.Name,
		input.PaymentMethodTypes,
	)
	business := domain.NewBusiness(
		businessID,
		input.Name,
		input.PaymentMethodTypes,
	)

	err = u.businessRepository.Save(ctx, businessEvent, business)
	if err != nil {
		return nil, err
	}

	return &CreateBusinessUseCaseOutput{
		Business: business,
	}, nil
}
