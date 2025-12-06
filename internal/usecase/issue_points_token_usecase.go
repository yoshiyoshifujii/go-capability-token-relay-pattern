package usecase

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

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

	issuePointsTokenUseCase struct {
		tokenService service.TokenService
	}
)

func NewIssuePointsTokenUseCase(tokenService service.TokenService) IssuePointsTokenUseCase {
	return &issuePointsTokenUseCase{
		tokenService: tokenService,
	}
}

func (u *issuePointsTokenUseCase) Execute(ctx context.Context, input IssuePointsTokenUseCaseInput) (*IssuePointsTokenUseCaseOutput, error) {
	token, err := u.tokenService.IssuePointsToken(ctx, service.IssuePointsTokenInput{
		OrderProcessingID: input.OrderProcessingID,
		UserID:            input.UserID,
		PointsToUse:       input.PointsToUse,
	})
	if err != nil {
		return nil, err
	}

	return &IssuePointsTokenUseCaseOutput{
		Token: token,
	}, nil
}
