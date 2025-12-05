package usecase

import "context"

type (
	IssuePaymentTokenUseCaseInput struct {
		OrderProcessingID string
		UserID            string
		PaymentMethod     string
	}

	IssuePaymentTokenUseCaseOutput struct {
		Token string
	}

	IssuePaymentTokenUseCase interface {
		Execute(context.Context, IssuePaymentTokenUseCaseInput) (*IssuePaymentTokenUseCaseOutput, error)
	}

	issuePaymentTokenUseCase struct{}
)

func NewIssuePaymentTokenUseCase() IssuePaymentTokenUseCase {
	return &issuePaymentTokenUseCase{}
}

func (u *issuePaymentTokenUseCase) Execute(ctx context.Context, input IssuePaymentTokenUseCaseInput) (*IssuePaymentTokenUseCaseOutput, error) {
	return &IssuePaymentTokenUseCaseOutput{}, nil
}
