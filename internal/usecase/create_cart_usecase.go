package usecase

import "context"

type (
	CreateCartUseCaseInput struct{}

	CreateCartUseCaseOutput struct{}

	CreateCartUseCase interface {
		Execute(context.Context, CreateCartUseCaseInput) (*CreateCartUseCaseOutput, error)
	}

	createCartUseCase struct{}
)

func NewCreateCartUseCase() CreateCartUseCase {
	return &createCartUseCase{}
}

func (c *createCartUseCase) Execute(ctx context.Context, input CreateCartUseCaseInput) (*CreateCartUseCaseOutput, error) {
	return &CreateCartUseCaseOutput{}, nil
}
