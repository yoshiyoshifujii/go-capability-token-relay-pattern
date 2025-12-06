package usecase

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
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
	if tokenService == nil {
		panic("tokenService is nil")
	}
	return &issuePaymentTokenUseCase{
		tokenService: tokenService,
	}
}

func (i IssuePaymentTokenUseCaseInput) Validate() error {
	if i.OrderProcessingID == "" {
		return errors.New("orderProcessingID is empty")
	}
	if i.UserID == "" {
		return errors.New("userID is empty")
	}
	if i.PaymentMethod == "" {
		return errors.New("paymentMethod is empty")
	}
	return nil
}

func (u *issuePaymentTokenUseCase) Execute(ctx context.Context, input IssuePaymentTokenUseCaseInput) (*IssuePaymentTokenUseCaseOutput, error) {
	contract.AssertValidatable(input)

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
