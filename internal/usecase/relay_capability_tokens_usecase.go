package usecase

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
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
	if tokenService == nil {
		panic("tokenService is nil")
	}
	return &relayCapabilityTokensUseCase{
		tokenService: tokenService,
	}
}

func (i RelayCapabilityTokensUseCaseInput) Validate() error {
	if i.OrderProcessingID == "" {
		return errors.New("orderProcessingID is empty")
	}
	contract.AssertValidatable(i.CartToken)
	contract.AssertValidatable(i.PaymentToken)
	return nil
}

func (u *relayCapabilityTokensUseCase) Execute(ctx context.Context, input RelayCapabilityTokensUseCaseInput) (*RelayCapabilityTokensUseCaseOutput, error) {
	contract.AssertValidatable(input)

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
