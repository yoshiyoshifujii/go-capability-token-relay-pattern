package usecase

import (
	"context"
	"errors"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/repository"
)

type (
	CreateBusinessUseCaseInput struct {
		BusinessID string
		Name       string
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
	return &createBusinessUseCase{
		businessIDGenerator: businessIDGenerator,
		businessRepository:  businessRepository,
	}
}

func (u *createBusinessUseCase) Execute(ctx context.Context, input CreateBusinessUseCaseInput) (*CreateBusinessUseCaseOutput, error) {
	if u.businessIDGenerator == nil {
		return nil, errors.New("businessIDGenerator is nil")
	}
	if u.businessRepository == nil {
		return nil, errors.New("businessRepository is nil")
	}

	businessID, err := u.businessIDGenerator.GenerateID(ctx)
	if err != nil {
		return nil, err
	}

	businessEvent := domain.NewBusinessInitializedEvent(
		businessID,
		1,
		input.Name,
	)
	business := domain.NewBusiness(
		businessID,
		input.Name,
	)

	err = u.businessRepository.Save(ctx, businessEvent, business)
	if err != nil {
		return nil, err
	}

	return &CreateBusinessUseCaseOutput{
		Business: business,
	}, nil
}
