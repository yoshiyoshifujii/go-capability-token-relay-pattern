package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	iarepo "yoshiyoshifujii/go-capability-token-relay-pattern/internal/interface_adaptor/repository"
)

func seedPaymentIntentRequiresConfirmation(
	t *testing.T,
	ctx context.Context,
	repo *iarepo.InMemoryPaymentIntentRepository,
	captureMethod domain.PaymentCaptureMethod,
) domain.PaymentIntentRequiresConfirmation {
	t.Helper()

	paymentIntentID := domain.PaymentIntentID("pi_fail_test")
	paymentMethodType := domain.PaymentMethodTypeCard
	amount := domain.NewMoney(120)
	paymentMethod := domain.NewPaymentMethod(
		paymentMethodType,
		&domain.PaymentMethodCard{
			Number:   "4242424242424242",
			ExpYear:  25,
			ExpMonth: 12,
		},
		nil,
	)

	event, aggregate, err := domain.GeneratePaymentIntent(paymentIntentID, domain.PaymentMethodTypes{paymentMethodType}, amount)
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, event, aggregate))

	event, aggregate, err = aggregate.(domain.PaymentIntentRequiresPaymentMethodType).RequirePaymentMethod(paymentMethodType)
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, event, aggregate))

	event, aggregate, err = aggregate.(domain.PaymentIntentRequiresPaymentMethod).RequireConfirmation(paymentMethod, captureMethod)
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, event, aggregate))

	return aggregate.(domain.PaymentIntentRequiresConfirmation)
}
