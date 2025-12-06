package usecase

import "context"

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

func (u *completeOrderUseCase) Execute(ctx context.Context, input CompleteOrderUseCaseInput) (*CompleteOrderUseCaseOutput, error) {
	return &CompleteOrderUseCaseOutput{
		OrderID: input.OrderProcessingID,
	}, nil
}
