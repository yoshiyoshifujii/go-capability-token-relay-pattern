package usecase

import (
	"context"
	"errors"
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
	return &createCartUseCase{
		cartIDGenerator: cartIDGenerator,
	}
}

func (u *createCartUseCase) Execute(ctx context.Context, input CreateCartUseCaseInput) (*CreateCartUseCaseOutput, error) {
	if u == nil {
		panic("usecase is nil")
	}
	if u.cartIDGenerator == nil {
		return nil, errors.New("cartIDGenerator is nil")
	}
	contract.AssertValidatable(input.BusinessID)
	contract.AssertValidatable(input.Items)

	cartID, err := u.cartIDGenerator.GenerateID(ctx)
	if err != nil {
		return nil, err
	}

	cart := domain.NewCart(input.BusinessID, cartID, input.Items)

	return &CreateCartUseCaseOutput{
		Cart: cart,
	}, nil
}
