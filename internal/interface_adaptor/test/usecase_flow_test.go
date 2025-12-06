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
	businessIDGenerator := iasvc.NewFakeBusinessIDGenerator(domain.NewBusinessID("biz_123"))
	businessRepo := iarepo.NewInMemoryBusinessRepository()
	cartIDGenerator := iasvc.NewFakeCartIDGenerator(domain.NewCartID("cart_123"))
	tokenService := iasvc.NewTokenService()
	paymentIntentRepo := iarepo.NewInMemoryPaymentIntentRepository()
	paymentIntentIDGenerator := iasvc.NewFakePaymentIntentIDGenerator(domain.PaymentIntentID("pi_123"))
	paymentProvider := iasvc.NewPaymentMethodProviderService(domain.PaymentConfirmationNextRequiresAction)

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
		Items: domain.NewCartItems(
			domain.CartItem{
				ItemID: domain.ItemID("item_123"),
				Price:  domain.ItemPrice(120),
			},
		),
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
		PaymentIntentID: selectedView.ID,
		CaptureMethod:   domain.PaymentCaptureMethodManual,
		PaymentMethod: domain.NewPaymentMethod(
			selectedView.PaymentMethodType,
			&domain.PaymentMethodCard{
				Number:   "4242424242424242",
				ExpYear:  25,
				ExpMonth: 12,
			},
			nil,
		),
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

	confirmPaymentIntent := usecase.NewConfirmPaymentIntentUseCase(paymentIntentRepo, paymentProvider)
	confirmPaymentIntentOutput, err := confirmPaymentIntent.Execute(ctx, usecase.ConfirmPaymentIntentUseCaseInput{
		PaymentIntentID: paymentIntentOutput.PaymentIntentID,
	})
	assert.NoError(t, err)
	assert.NotNil(t, confirmPaymentIntentOutput)
	assert.Len(t, paymentIntentRepo.Events(), 4)

	confirmedPaymentIntent, err := paymentIntentRepo.FindBy(ctx, paymentIntentOutput.PaymentIntentID)
	assert.NoError(t, err)
	assert.NotNil(t, confirmedPaymentIntent)
	confirmedView, err := converter.ToPaymentIntentView(*confirmedPaymentIntent)
	assert.NoError(t, err)
	assert.Equal(t, "requires_action", confirmedView.Status)

	// webhook after user completed 3DS action
	handleActionResult := usecase.NewHandlePaymentActionResultUseCase(paymentIntentRepo)
	handleActionResultOutput, err := handleActionResult.Execute(ctx, usecase.HandlePaymentActionResultUseCaseInput{
		PaymentIntentID: paymentIntentOutput.PaymentIntentID,
	})
	assert.NoError(t, err)
	assert.NotNil(t, handleActionResultOutput)
	assert.Len(t, paymentIntentRepo.Events(), 5)

	actionHandledIntent, err := paymentIntentRepo.FindBy(ctx, paymentIntentOutput.PaymentIntentID)
	assert.NoError(t, err)
	assert.NotNil(t, actionHandledIntent)
	actionHandledView, err := converter.ToPaymentIntentView(*actionHandledIntent)
	assert.NoError(t, err)
	assert.Equal(t, "requires_capture", actionHandledView.Status)

	capturePaymentIntent := usecase.NewCapturePaymentIntentUseCase(paymentIntentRepo, paymentProvider)
	capturePaymentIntentOutput, err := capturePaymentIntent.Execute(ctx, usecase.CapturePaymentIntentUseCaseInput{
		PaymentIntentID: paymentIntentOutput.PaymentIntentID,
	})
	assert.NoError(t, err)
	assert.NotNil(t, capturePaymentIntentOutput)
	assert.Len(t, paymentIntentRepo.Events(), 6)

	capturedIntent, err := paymentIntentRepo.FindBy(ctx, paymentIntentOutput.PaymentIntentID)
	assert.NoError(t, err)
	assert.NotNil(t, capturedIntent)
	capturedView, err := converter.ToPaymentIntentView(*capturedIntent)
	assert.NoError(t, err)
	assert.Equal(t, "processing", capturedView.Status)

	handlePaymentSucceeded := usecase.NewHandlePaymentSucceededUseCase(paymentIntentRepo)
	handlePaymentSucceededOutput, err := handlePaymentSucceeded.Execute(ctx, usecase.HandlePaymentSucceededUseCaseInput{
		PaymentIntentID: paymentIntentOutput.PaymentIntentID,
	})
	assert.NoError(t, err)
	assert.NotNil(t, handlePaymentSucceededOutput)
	assert.Len(t, paymentIntentRepo.Events(), 7)

	completedIntent, err := paymentIntentRepo.FindBy(ctx, paymentIntentOutput.PaymentIntentID)
	assert.NoError(t, err)
	assert.NotNil(t, completedIntent)
	completedView, err := converter.ToPaymentIntentView(*completedIntent)
	assert.NoError(t, err)
	assert.Equal(t, "succeeded", completedView.Status)
}
