package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	iarepo "yoshiyoshifujii/go-capability-token-relay-pattern/internal/interface_adaptor/repository"
)

func TestCapturePaymentIntentUseCase_ShouldCancelOnProviderError(t *testing.T) {
	ctx := context.Background()
	repo := iarepo.NewInMemoryPaymentIntentRepository()

	intent := seedPaymentIntentRequiresCapture(t, ctx, repo)
	useCase := NewCapturePaymentIntentUseCase(repo, &fakePaymentMethodProvider{captureErr: errors.New("capture failed")})

	output, err := useCase.Execute(ctx, CapturePaymentIntentUseCaseInput{
		PaymentIntentID: intent.ID,
	})

	require.Error(t, err)
	require.NotNil(t, output)
	assert.Len(t, repo.Events(), 5)

	result, ok := output.PaymentIntent.(domain.PaymentIntentCanceled)
	require.True(t, ok)
	assert.Equal(t, domain.PaymentFailureReasonCaptureFailed, result.FailureReason)
	assert.Equal(t, intent.PaymentMethod, result.PaymentMethod)
}

func seedPaymentIntentRequiresCapture(
	t *testing.T,
	ctx context.Context,
	repo *iarepo.InMemoryPaymentIntentRepository,
) domain.PaymentIntentRequiresCapture {
	t.Helper()

	confirmation := seedPaymentIntentRequiresConfirmation(t, ctx, repo, domain.PaymentCaptureMethodManual)
	event, aggregate, err := confirmation.RequireCapture()
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, event, aggregate))

	return aggregate.(domain.PaymentIntentRequiresCapture)
}
