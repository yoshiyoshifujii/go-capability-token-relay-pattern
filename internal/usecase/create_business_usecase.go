package usecase

import (
	"context"
	"errors"

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
		businessRepository repository.BusinessRepository
	}
)

func NewCreateBusinessUseCase(businessRepository repository.BusinessRepository) CreateBusinessUseCase {
	return &createBusinessUseCase{
		businessRepository: businessRepository,
	}
}

func (u *createBusinessUseCase) Execute(ctx context.Context, input CreateBusinessUseCaseInput) (*CreateBusinessUseCaseOutput, error) {
	if u.businessRepository == nil {
		return nil, errors.New("businessRepository is nil")
	}

	businessID := domain.NewBusinessID(input.BusinessID)

	businessEvent := domain.NewBusinessInitializedEvent(
		businessID,
		1,
		input.Name,
	)
	business := domain.NewBusiness(
		businessID,
		input.Name,
	)

	err := u.businessRepository.Save(ctx, businessEvent, business)
	if err != nil {
		return nil, err
	}

	return &CreateBusinessUseCaseOutput{
		Business: business,
	}, nil
}
