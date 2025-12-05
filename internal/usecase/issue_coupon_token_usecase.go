package usecase

import "context"

type (
	IssueCouponTokenUseCaseInput struct {
		OrderProcessingID string
		UserID            string
		CouponRef         string
	}

	IssueCouponTokenUseCaseOutput struct {
		Token string
	}

	IssueCouponTokenUseCase interface {
		Execute(context.Context, IssueCouponTokenUseCaseInput) (*IssueCouponTokenUseCaseOutput, error)
	}

	issueCouponTokenUseCase struct{}
)

func NewIssueCouponTokenUseCase() IssueCouponTokenUseCase {
	return &issueCouponTokenUseCase{}
}

func (u *issueCouponTokenUseCase) Execute(ctx context.Context, input IssueCouponTokenUseCaseInput) (*IssueCouponTokenUseCaseOutput, error) {
	return &IssueCouponTokenUseCaseOutput{}, nil
}
