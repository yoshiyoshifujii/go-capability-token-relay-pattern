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
	paymentIntentRepo := iarepo.NewInMemoryPaymentIntentRepository()
	businessIDGenerator := iasvc.NewFakeBusinessIDGenerator(domain.NewBusinessID("biz_123"))
	cartIDGenerator := iasvc.NewFakeCartIDGenerator(domain.NewCartID("cart_123"))
	paymentIntentIDGenerator := iasvc.NewFakePaymentIntentIDGenerator(domain.PaymentIntentID("pi_123"))

	// create business
	createBusiness := usecase.NewCreateBusinessUseCase(businessIDGenerator, businessRepo)
	businessOutput, err := createBusiness.Execute(ctx, usecase.CreateBusinessUseCaseInput{
		BusinessID:         "biz_123",
		Name:               "Test Business",
		PaymentMethodTypes: domain.PaymentMethodTypes{domain.PaymentMethodTypeCard},
	})
	assert.NoError(t, err)
	assert.NotNil(t, businessOutput)
	assert.Equal(t, "biz_123", string(businessOutput.Business.ID))
	assert.Len(t, businessRepo.Events(), 1)

	// create cart
	createCart := usecase.NewCreateCartUseCase(cartIDGenerator)
	createCartOutput, err := createCart.Execute(ctx, usecase.CreateCartUseCaseInput{
		BusinessID: businessOutput.Business.ID,
		Items:      domain.NewCartItems(domain.ItemID("item_123")),
	})
	assert.NoError(t, err)
	assert.NotNil(t, createCartOutput)
	assert.Equal(t, domain.CartID("cart_123"), createCartOutput.Cart.CartID)
	assert.Equal(t, businessOutput.Business.ID, createCartOutput.Cart.BusinessID)

	// confirm cart token
	confirmCart := usecase.NewConfirmCartUseCase(tokenService)
	confirmCartOutput, err := confirmCart.Execute(ctx, usecase.ConfirmCartUseCaseInput{
		Cart: createCartOutput.Cart,
	})
	assert.NoError(t, err)
	assert.NotNil(t, confirmCartOutput)
	assert.NotEmpty(t, confirmCartOutput.Token.Value)

	// initialize payment intent
	initializePaymentIntent := usecase.NewInitializePaymentIntentUseCase(tokenService, paymentIntentRepo, paymentIntentIDGenerator, businessRepo)
	paymentIntentOutput, err := initializePaymentIntent.Execute(ctx, usecase.InitializePaymentIntentUseCaseInput{
		CartToken: confirmCartOutput.Token,
	})
	assert.NoError(t, err)
	assert.NotNil(t, paymentIntentOutput)
	assert.Len(t, paymentIntentRepo.Events(), 1)

	// issue confirmation tokens from each domain
	paymentToken := usecase.NewIssuePaymentTokenUseCase(tokenService)
	paymentTokenOutput, err := paymentToken.Execute(ctx, usecase.IssuePaymentTokenUseCaseInput{
		OrderProcessingID: "op_123",
		UserID:            "user_123",
		PaymentMethod:     "credit-card",
	})
	assert.NoError(t, err)
	assert.NotNil(t, paymentTokenOutput)
	assert.NotEmpty(t, paymentTokenOutput.Token.Value)

	// relay tokens (OrderProcessing)
	relay := usecase.NewRelayCapabilityTokensUseCase(tokenService)
	relayOutput, err := relay.Execute(ctx, usecase.RelayCapabilityTokensUseCaseInput{
		OrderProcessingID: "op_123",
		CartToken:         confirmCartOutput.Token,
		PaymentToken:      paymentTokenOutput.Token,
	})
	assert.NoError(t, err)
	assert.NotNil(t, relayOutput)
	assert.NotNil(t, relayOutput.VerifiedTokens)
	assert.Len(t, relayOutput.VerifiedTokens, 2)
	assert.Equal(t, confirmCartOutput.Token.Value, relayOutput.VerifiedTokens["cart"].Value)
	assert.Equal(t, paymentTokenOutput.Token.Value, relayOutput.VerifiedTokens["payment"].Value)

	// complete order
	completeOrder := usecase.NewCompleteOrderUseCase()
	completeOutput, err := completeOrder.Execute(ctx, usecase.CompleteOrderUseCaseInput{
		OrderProcessingID: "op_123",
		CapabilityTokens: []string{
			confirmCartOutput.Token.Value,
			paymentTokenOutput.Token.Value,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, completeOutput)
	assert.NotEmpty(t, completeOutput.OrderID)
}
