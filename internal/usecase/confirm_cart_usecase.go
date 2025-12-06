package usecase

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type (
	ConfirmCartUseCaseInput struct {
		Cart domain.Cart
	}

	ConfirmCartUseCaseOutput struct {
		Token service.SignedToken
	}

	ConfirmCartUseCase interface {
		Execute(context.Context, ConfirmCartUseCaseInput) (*ConfirmCartUseCaseOutput, error)
	}

	confirmCartUseCase struct {
		tokenService service.TokenService
	}
)

func NewConfirmCartUseCase(tokenService service.TokenService) ConfirmCartUseCase {
	if tokenService == nil {
		panic("tokenService is nil")
	}
	return &confirmCartUseCase{
		tokenService: tokenService,
	}
}

func (u *confirmCartUseCase) Execute(ctx context.Context, input ConfirmCartUseCaseInput) (*ConfirmCartUseCaseOutput, error) {
	contract.AssertValidatable(input.Cart)

	token, err := u.tokenService.ConfirmCartToken(ctx, service.ConfirmCartTokenInput{
		Cart: input.Cart,
	})
	if err != nil {
		return nil, err
	}

	return &ConfirmCartUseCaseOutput{
		Token: token,
	}, nil
}
