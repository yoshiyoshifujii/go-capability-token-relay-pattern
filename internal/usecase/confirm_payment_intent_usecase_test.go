package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	iarepo "yoshiyoshifujii/go-capability-token-relay-pattern/internal/interface_adaptor/repository"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type fakePaymentMethodProvider struct {
	confirmErr error
	captureErr error
}

func (f *fakePaymentMethodProvider) ConfirmPaymentMethod(context.Context, service.PaymentConfirmationRequest) (service.PaymentConfirmationResult, error) {
	if f.confirmErr != nil {
		return service.PaymentConfirmationResult{}, f.confirmErr
	}
	return service.PaymentConfirmationResult{NextStatus: domain.PaymentConfirmationNextProcessing}, nil
}

func (f *fakePaymentMethodProvider) CapturePaymentIntent(context.Context, service.PaymentCaptureRequest) error {
	return f.captureErr
}

func TestConfirmPaymentIntentUseCase_ShouldFailOnProviderError(t *testing.T) {
	ctx := context.Background()
	repo := iarepo.NewInMemoryPaymentIntentRepository()

	intent := seedPaymentIntentRequiresConfirmation(t, ctx, repo, domain.PaymentCaptureMethodAutomatic)
	useCase := NewConfirmPaymentIntentUseCase(repo, &fakePaymentMethodProvider{confirmErr: errors.New("provider down")})

	output, err := useCase.Execute(ctx, ConfirmPaymentIntentUseCaseInput{
		PaymentIntentID: intent.ID,
	})

	require.Error(t, err)
	require.NotNil(t, output)
	assert.Len(t, repo.Events(), 4)

	result, ok := output.PaymentIntent.(domain.PaymentIntentRequiresPaymentMethod)
	require.True(t, ok)
	assert.Equal(t, domain.PaymentFailureReasonConfirmationFailed, result.FailureReason)
	assert.Equal(t, domain.PaymentMethodTypeCard, result.PaymentMethodType)
}
