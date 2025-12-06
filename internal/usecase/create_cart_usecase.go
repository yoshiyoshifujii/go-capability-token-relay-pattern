package usecase

import (
	"context"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type (
	CreateCartUseCaseInput struct {
		BusinessID domain.BusinessID
		Items      domain.CartItems
	}

	CreateCartUseCaseOutput struct {
		Cart domain.Cart
	}

	CreateCartUseCase interface {
		Execute(context.Context, CreateCartUseCaseInput) (*CreateCartUseCaseOutput, error)
	}

	createCartUseCase struct {
		cartIDGenerator service.CartIDGenerator
	}
)

func NewCreateCartUseCase(cartIDGenerator service.CartIDGenerator) CreateCartUseCase {
	if cartIDGenerator == nil {
		panic("cartIDGenerator is nil")
	}
	return &createCartUseCase{
		cartIDGenerator: cartIDGenerator,
	}
}

func (i CreateCartUseCaseInput) Validate() error {
	contract.AssertValidatable(i.BusinessID)
	contract.AssertValidatable(i.Items)
	return nil
}

func (u *createCartUseCase) Execute(ctx context.Context, input CreateCartUseCaseInput) (*CreateCartUseCaseOutput, error) {
	contract.AssertValidatable(input)

	cartID, err := u.cartIDGenerator.GenerateID(ctx)
	if err != nil {
		return nil, err
	}

	cart := domain.NewCart(input.BusinessID, cartID, input.Items)

	return &CreateCartUseCaseOutput{
		Cart: cart,
	}, nil
}
