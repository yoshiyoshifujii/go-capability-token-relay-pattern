package usecase

import "context"

type (
	RelayCapabilityTokensUseCaseInput struct {
		OrderProcessingID string
		CouponToken       string
		PointsToken       string
		PaymentToken      string
	}

	RelayCapabilityTokensUseCaseOutput struct {
		VerifiedTokens map[string]string
	}

	RelayCapabilityTokensUseCase interface {
		Execute(context.Context, RelayCapabilityTokensUseCaseInput) (*RelayCapabilityTokensUseCaseOutput, error)
	}

	relayCapabilityTokensUseCase struct{}
)

func NewRelayCapabilityTokensUseCase() RelayCapabilityTokensUseCase {
	return &relayCapabilityTokensUseCase{}
}

func (u *relayCapabilityTokensUseCase) Execute(ctx context.Context, input RelayCapabilityTokensUseCaseInput) (*RelayCapabilityTokensUseCaseOutput, error) {
	return &RelayCapabilityTokensUseCaseOutput{
		VerifiedTokens: make(map[string]string),
	}, nil
}
