package usecase

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type (
	IssueCouponTokenUseCaseInput struct {
		OrderProcessingID string
		UserID            string
		CouponRef         string
	}

	IssueCouponTokenUseCaseOutput struct {
		Token service.SignedToken
	}

	IssueCouponTokenUseCase interface {
		Execute(context.Context, IssueCouponTokenUseCaseInput) (*IssueCouponTokenUseCaseOutput, error)
	}

	issueCouponTokenUseCase struct {
		tokenService service.TokenService
	}
)

func NewIssueCouponTokenUseCase(tokenService service.TokenService) IssueCouponTokenUseCase {
	return &issueCouponTokenUseCase{
		tokenService: tokenService,
	}
}

func (u *issueCouponTokenUseCase) Execute(ctx context.Context, input IssueCouponTokenUseCaseInput) (*IssueCouponTokenUseCaseOutput, error) {
	token, err := u.tokenService.IssueCouponToken(ctx, service.IssueCouponTokenInput{
		OrderProcessingID: input.OrderProcessingID,
		UserID:            input.UserID,
		CouponRef:         input.CouponRef,
	})
	if err != nil {
		return nil, err
	}

	return &IssueCouponTokenUseCaseOutput{
		Token: token,
	}, nil
}
