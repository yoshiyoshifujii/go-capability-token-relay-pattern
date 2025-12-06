package usecase

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type (
	RelayCapabilityTokensUseCaseInput struct {
		OrderProcessingID string
		CartToken         service.SignedToken
		PaymentToken      service.SignedToken
	}

	RelayCapabilityTokensUseCaseOutput struct {
		VerifiedTokens map[string]service.SignedToken
	}

	RelayCapabilityTokensUseCase interface {
		Execute(context.Context, RelayCapabilityTokensUseCaseInput) (*RelayCapabilityTokensUseCaseOutput, error)
	}

	relayCapabilityTokensUseCase struct {
		tokenService service.TokenService
	}
)

func NewRelayCapabilityTokensUseCase(tokenService service.TokenService) RelayCapabilityTokensUseCase {
	return &relayCapabilityTokensUseCase{
		tokenService: tokenService,
	}
}

func (u *relayCapabilityTokensUseCase) Execute(ctx context.Context, input RelayCapabilityTokensUseCaseInput) (*RelayCapabilityTokensUseCaseOutput, error) {
	verified, err := u.tokenService.RelayTokens(ctx, service.RelayTokensInput{
		OrderProcessingID: input.OrderProcessingID,
		CartToken:         input.CartToken.Value,
		PaymentToken:      input.PaymentToken.Value,
	})
	if err != nil {
		return nil, err
	}

	return &RelayCapabilityTokensUseCaseOutput{
		VerifiedTokens: verified,
	}, nil
}
