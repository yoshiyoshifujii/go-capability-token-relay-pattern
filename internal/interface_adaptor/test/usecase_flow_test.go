package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	iarepo "yoshiyoshifujii/go-capability-token-relay-pattern/internal/interface_adaptor/repository"
	iasvc "yoshiyoshifujii/go-capability-token-relay-pattern/internal/interface_adaptor/service"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/usecase"
)

func TestUseCaseFlow_ShouldPassThroughAllStubs(t *testing.T) {
	ctx := t.Context()
	tokenService := iasvc.NewTokenService()
	businessRepo := iarepo.NewInMemoryBusinessRepository()
	businessIDGenerator := iasvc.NewFakeBusinessIDGenerator(domain.NewBusinessID("biz_123"))

	// create business
	createBusiness := usecase.NewCreateBusinessUseCase(businessIDGenerator, businessRepo)
	businessOutput, err := createBusiness.Execute(ctx, usecase.CreateBusinessUseCaseInput{
		BusinessID: "biz_123",
		Name:       "Test Business",
	})
	assert.NoError(t, err)
	assert.NotNil(t, businessOutput)
	assert.Equal(t, "biz_123", string(businessOutput.Business.ID))
	assert.Len(t, businessRepo.Events, 1)

	// create cart
	createCart := usecase.NewCreateCartUseCase()
	createCartOutput, err := createCart.Execute(ctx, usecase.CreateCartUseCaseInput{})
	assert.NoError(t, err)
	assert.NotNil(t, createCartOutput)

	// issue confirmation tokens from each domain
	couponToken := usecase.NewIssueCouponTokenUseCase(tokenService)
	couponTokenOutput, err := couponToken.Execute(ctx, usecase.IssueCouponTokenUseCaseInput{
		OrderProcessingID: "op_123",
		UserID:            "user_123",
		CouponRef:         "coupon_abc",
	})
	assert.NoError(t, err)
	assert.NotNil(t, couponTokenOutput)
	assert.NotEmpty(t, couponTokenOutput.Token)

	pointsToken := usecase.NewIssuePointsTokenUseCase(tokenService)
	pointsTokenOutput, err := pointsToken.Execute(ctx, usecase.IssuePointsTokenUseCaseInput{
		OrderProcessingID: "op_123",
		UserID:            "user_123",
		PointsToUse:       100,
	})
	assert.NoError(t, err)
	assert.NotNil(t, pointsTokenOutput)
	assert.NotEmpty(t, pointsTokenOutput.Token)

	paymentToken := usecase.NewIssuePaymentTokenUseCase(tokenService)
	paymentTokenOutput, err := paymentToken.Execute(ctx, usecase.IssuePaymentTokenUseCaseInput{
		OrderProcessingID: "op_123",
		UserID:            "user_123",
		PaymentMethod:     "credit-card",
	})
	assert.NoError(t, err)
	assert.NotNil(t, paymentTokenOutput)
	assert.NotEmpty(t, paymentTokenOutput.Token)

	// relay tokens (OrderProcessing)
	relay := usecase.NewRelayCapabilityTokensUseCase(tokenService)
	relayOutput, err := relay.Execute(ctx, usecase.RelayCapabilityTokensUseCaseInput{
		OrderProcessingID: "op_123",
		CouponToken:       couponTokenOutput.Token,
		PointsToken:       pointsTokenOutput.Token,
		PaymentToken:      paymentTokenOutput.Token,
	})
	assert.NoError(t, err)
	assert.NotNil(t, relayOutput)
	assert.NotNil(t, relayOutput.VerifiedTokens)
	assert.Len(t, relayOutput.VerifiedTokens, 3)

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
	assert.NotEmpty(t, completeOutput.OrderID)
}
