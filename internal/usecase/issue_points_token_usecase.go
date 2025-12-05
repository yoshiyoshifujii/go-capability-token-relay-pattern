package usecase

import "context"

type (
	IssuePointsTokenUseCaseInput struct {
		OrderProcessingID string
		UserID            string
		PointsToUse       int64
	}

	IssuePointsTokenUseCaseOutput struct {
		Token string
	}

	IssuePointsTokenUseCase interface {
		Execute(context.Context, IssuePointsTokenUseCaseInput) (*IssuePointsTokenUseCaseOutput, error)
	}

	issuePointsTokenUseCase struct{}
)

func NewIssuePointsTokenUseCase() IssuePointsTokenUseCase {
	return &issuePointsTokenUseCase{}
}

func (u *issuePointsTokenUseCase) Execute(ctx context.Context, input IssuePointsTokenUseCaseInput) (*IssuePointsTokenUseCaseOutput, error) {
	return &IssuePointsTokenUseCaseOutput{}, nil
}
