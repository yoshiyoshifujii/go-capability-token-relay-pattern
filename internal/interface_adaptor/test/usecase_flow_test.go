package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/interface_adaptor/converter"
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
	paymentIntentView, err := converter.ToPaymentIntentView(paymentIntentOutput.PaymentIntent)
	assert.NoError(t, err)

	// select payment method
	selectPaymentMethod := usecase.NewSelectPaymentMethodUseCase(paymentIntentRepo)
	selectPaymentMethodOutput, err := selectPaymentMethod.Execute(ctx, usecase.SelectPaymentMethodUseCaseInput{
		PaymentIntentID:   paymentIntentOutput.PaymentIntentID,
		PaymentMethodType: paymentIntentView.PaymentMethodTypes[0],
	})
	assert.NoError(t, err)
	assert.NotNil(t, selectPaymentMethodOutput)
	assert.Len(t, paymentIntentRepo.Events(), 2)
	selectedPaymentIntent, err := paymentIntentRepo.FindBy(ctx, paymentIntentOutput.PaymentIntentID)
	assert.NoError(t, err)
	assert.NotNil(t, selectedPaymentIntent)
	selectedView, err := converter.ToPaymentIntentView(*selectedPaymentIntent)
	assert.NoError(t, err)

	providePaymentMethod := usecase.NewProvidePaymentMethodUseCase(paymentIntentRepo)
	providePaymentMethodOutput, err := providePaymentMethod.Execute(ctx, usecase.ProvidePaymentMethodUseCaseInput{
		PaymentIntentID:   selectedView.ID,
		PaymentMethodType: selectedView.PaymentMethodType,
		Card: &domain.PaymentMethodCard{
			Number:   "4242424242424242",
			ExpYear:  25,
			ExpMonth: 12,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, providePaymentMethodOutput)
	assert.Len(t, paymentIntentRepo.Events(), 3)

	latestPaymentIntent, err := paymentIntentRepo.FindBy(ctx, paymentIntentOutput.PaymentIntentID)
	assert.NoError(t, err)
	assert.NotNil(t, latestPaymentIntent)
	latestView, err := converter.ToPaymentIntentView(*latestPaymentIntent)
	assert.NoError(t, err)
	assert.Equal(t, "requires_confirmation", latestView.Status)

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
			relayOutput.VerifiedTokens["cart"].Value,
			relayOutput.VerifiedTokens["payment"].Value,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, completeOutput)
	assert.NotEmpty(t, completeOutput.OrderID)
}
