package usecase

import (
	"context"
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/repository"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type (
	CapturePaymentIntentUseCaseInput struct {
		PaymentIntentID domain.PaymentIntentID
	}

	CapturePaymentIntentUseCaseOutput struct {
		PaymentIntentID domain.PaymentIntentID
		PaymentIntent   domain.PaymentIntent
	}

	CapturePaymentIntentUseCase interface {
		Execute(context.Context, CapturePaymentIntentUseCaseInput) (*CapturePaymentIntentUseCaseOutput, error)
	}

	capturePaymentIntentUseCase struct {
		paymentIntentRepository repository.PaymentIntentRepository
		paymentProvider         service.PaymentMethodProviderService
	}
)

func NewCapturePaymentIntentUseCase(
	paymentIntentRepository repository.PaymentIntentRepository,
	paymentProvider service.PaymentMethodProviderService,
) CapturePaymentIntentUseCase {
	if paymentIntentRepository == nil {
		panic("paymentIntentRepository is nil")
	}
	if paymentProvider == nil {
		panic("paymentProvider is nil")
	}
	return &capturePaymentIntentUseCase{
		paymentIntentRepository: paymentIntentRepository,
		paymentProvider:         paymentProvider,
	}
}

func (i CapturePaymentIntentUseCaseInput) Validate() error {
	contract.AssertValidatable(i.PaymentIntentID)
	return nil
}

func (u *capturePaymentIntentUseCase) Execute(ctx context.Context, input CapturePaymentIntentUseCaseInput) (*CapturePaymentIntentUseCaseOutput, error) {
	contract.AssertValidatable(input)

	paymentIntent, err := u.paymentIntentRepository.FindBy(ctx, input.PaymentIntentID)
	if err != nil {
		return nil, err
	}
	if paymentIntent == nil {
		return nil, errors.New("payment intent not found")
	}

	intent, ok := (*paymentIntent).(domain.PaymentIntentRequiresCapture)
	if !ok {
		return nil, errors.New("payment intent not ready for capture")
	}

	if err = u.paymentProvider.CapturePaymentIntent(ctx, service.PaymentCaptureRequest{
		Intent: intent,
		Amount: intent.Amount,
	}); err != nil {
		return nil, err
	}

	event, aggregate, err := intent.StartProcessing()
	if err != nil {
		return nil, err
	}

	if err := u.paymentIntentRepository.Save(ctx, event, aggregate); err != nil {
		return nil, err
	}

	return &CapturePaymentIntentUseCaseOutput{
		PaymentIntentID: input.PaymentIntentID,
		PaymentIntent:   aggregate,
	}, nil
}
