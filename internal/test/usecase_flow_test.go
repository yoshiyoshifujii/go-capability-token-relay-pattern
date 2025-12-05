package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/usecase"
)

func TestUseCaseFlow_ShouldPassThroughAllStubs(t *testing.T) {
	ctx := t.Context()

	// create cart
	createCart := usecase.NewCreateCartUseCase()
	createCartOutput, err := createCart.Execute(ctx, usecase.CreateCartUseCaseInput{})
	assert.NoError(t, err)
	assert.NotNil(t, createCartOutput)

	// issue confirmation tokens from each domain
	couponToken := usecase.NewIssueCouponTokenUseCase()
	couponTokenOutput, err := couponToken.Execute(ctx, usecase.IssueCouponTokenUseCaseInput{
		OrderProcessingID: "op_123",
		UserID:            "user_123",
		CouponRef:         "coupon_abc",
	})
	assert.NoError(t, err)
	assert.NotNil(t, couponTokenOutput)

	pointsToken := usecase.NewIssuePointsTokenUseCase()
	pointsTokenOutput, err := pointsToken.Execute(ctx, usecase.IssuePointsTokenUseCaseInput{
		OrderProcessingID: "op_123",
		UserID:            "user_123",
		PointsToUse:       100,
	})
	assert.NoError(t, err)
	assert.NotNil(t, pointsTokenOutput)

	paymentToken := usecase.NewIssuePaymentTokenUseCase()
	paymentTokenOutput, err := paymentToken.Execute(ctx, usecase.IssuePaymentTokenUseCaseInput{
		OrderProcessingID: "op_123",
		UserID:            "user_123",
		PaymentMethod:     "credit-card",
	})
	assert.NoError(t, err)
	assert.NotNil(t, paymentTokenOutput)

	// relay tokens (OrderProcessing)
	relay := usecase.NewRelayCapabilityTokensUseCase()
	relayOutput, err := relay.Execute(ctx, usecase.RelayCapabilityTokensUseCaseInput{
		OrderProcessingID: "op_123",
		CouponToken:       couponTokenOutput.Token,
		PointsToken:       pointsTokenOutput.Token,
		PaymentToken:      paymentTokenOutput.Token,
	})
	assert.NoError(t, err)
	assert.NotNil(t, relayOutput)
	assert.NotNil(t, relayOutput.VerifiedTokens)

	// complete order
	completeOrder := usecase.NewCompleteOrderUseCase()
	completeOutput, err := completeOrder.Execute(ctx, usecase.CompleteOrderUseCaseInput{
		OrderProcessingID: "op_123",
		CapabilityTokens: []string{
			couponTokenOutput.Token,
			pointsTokenOutput.Token,
			paymentTokenOutput.Token,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, completeOutput)
}
