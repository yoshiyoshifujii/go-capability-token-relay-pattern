package usecase

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type (
	IssuePaymentTokenUseCaseInput struct {
		OrderProcessingID string
		UserID            string
		PaymentMethod     string
	}

	IssuePaymentTokenUseCaseOutput struct {
		Token service.SignedToken
	}

	IssuePaymentTokenUseCase interface {
		Execute(context.Context, IssuePaymentTokenUseCaseInput) (*IssuePaymentTokenUseCaseOutput, error)
	}

	issuePaymentTokenUseCase struct {
		tokenService service.TokenService
	}
)

func NewIssuePaymentTokenUseCase(tokenService service.TokenService) IssuePaymentTokenUseCase {
	return &issuePaymentTokenUseCase{
		tokenService: tokenService,
	}
}

func (u *issuePaymentTokenUseCase) Execute(ctx context.Context, input IssuePaymentTokenUseCaseInput) (*IssuePaymentTokenUseCaseOutput, error) {
	token, err := u.tokenService.IssuePaymentToken(ctx, service.IssuePaymentTokenInput{
		OrderProcessingID: input.OrderProcessingID,
		UserID:            input.UserID,
		PaymentMethod:     input.PaymentMethod,
	})
	if err != nil {
		return nil, err
	}

	return &IssuePaymentTokenUseCaseOutput{
		Token: token,
	}, nil
}
