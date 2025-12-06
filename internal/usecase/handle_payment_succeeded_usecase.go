package usecase

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/repository"
)

type (
	HandlePaymentSucceededUseCaseInput struct {
		PaymentIntentID domain.PaymentIntentID
	}

	HandlePaymentSucceededUseCaseOutput struct {
		PaymentIntentID domain.PaymentIntentID
		PaymentIntent   domain.PaymentIntent
	}

	HandlePaymentSucceededUseCase interface {
		Execute(context.Context, HandlePaymentSucceededUseCaseInput) (*HandlePaymentSucceededUseCaseOutput, error)
	}

	handlePaymentSucceededUseCase struct {
		paymentIntentRepository repository.PaymentIntentRepository
	}
)

func NewHandlePaymentSucceededUseCase(paymentIntentRepository repository.PaymentIntentRepository) HandlePaymentSucceededUseCase {
	if paymentIntentRepository == nil {
		panic("paymentIntentRepository is nil")
	}
	return &handlePaymentSucceededUseCase{
		paymentIntentRepository: paymentIntentRepository,
	}
}

func (i HandlePaymentSucceededUseCaseInput) Validate() error {
	contract.AssertValidatable(i.PaymentIntentID)
	return nil
}

func (u *handlePaymentSucceededUseCase) Execute(ctx context.Context, input HandlePaymentSucceededUseCaseInput) (*HandlePaymentSucceededUseCaseOutput, error) {
	contract.AssertValidatable(input)

	paymentIntent, err := u.paymentIntentRepository.FindBy(ctx, input.PaymentIntentID)
	if err != nil {
		return nil, err
	}
	if paymentIntent == nil {
		// Webhookが先に届いた場合は黙って成功
		return &HandlePaymentSucceededUseCaseOutput{
			PaymentIntentID: input.PaymentIntentID,
		}, nil
	}

	intent, ok := (*paymentIntent).(domain.PaymentIntentProcessing)
	if !ok {
		// すでに進んでいても成功として返す
		return &HandlePaymentSucceededUseCaseOutput{
			PaymentIntentID: input.PaymentIntentID,
			PaymentIntent:   *paymentIntent,
		}, nil
	}

	event, aggregate, err := intent.Complete()
	if err != nil {
		return nil, err
	}

	if err := u.paymentIntentRepository.Save(ctx, event, aggregate); err != nil {
		return nil, err
	}

	return &HandlePaymentSucceededUseCaseOutput{
		PaymentIntentID: input.PaymentIntentID,
		PaymentIntent:   aggregate,
	}, nil
}
