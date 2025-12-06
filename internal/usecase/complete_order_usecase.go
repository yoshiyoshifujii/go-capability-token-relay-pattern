package usecase

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
)

type (
	CompleteOrderUseCaseInput struct {
		OrderProcessingID string
		CapabilityTokens  []string
	}

	CompleteOrderUseCaseOutput struct {
		OrderID string
	}

	CompleteOrderUseCase interface {
		Execute(context.Context, CompleteOrderUseCaseInput) (*CompleteOrderUseCaseOutput, error)
	}

	completeOrderUseCase struct{}
)

func NewCompleteOrderUseCase() CompleteOrderUseCase {
	return &completeOrderUseCase{}
}

func (i CompleteOrderUseCaseInput) Validate() error {
	if i.OrderProcessingID == "" {
		return errors.New("orderProcessingID is empty")
	}
	if len(i.CapabilityTokens) == 0 {
		return errors.New("capability tokens are empty")
	}
	return nil
}

func (u *completeOrderUseCase) Execute(ctx context.Context, input CompleteOrderUseCaseInput) (*CompleteOrderUseCaseOutput, error) {
	contract.AssertValidatable(input)

	return &CompleteOrderUseCaseOutput{
		OrderID: input.OrderProcessingID,
	}, nil
}
