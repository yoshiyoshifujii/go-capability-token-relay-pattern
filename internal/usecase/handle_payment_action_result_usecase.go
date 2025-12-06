package usecase

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/repository"
)

type (
	HandlePaymentActionResultUseCaseInput struct {
		PaymentIntentID domain.PaymentIntentID
	}

	HandlePaymentActionResultUseCaseOutput struct {
		PaymentIntentID domain.PaymentIntentID
		PaymentIntent   domain.PaymentIntent
	}

	HandlePaymentActionResultUseCase interface {
		Execute(context.Context, HandlePaymentActionResultUseCaseInput) (*HandlePaymentActionResultUseCaseOutput, error)
	}

	handlePaymentActionResultUseCase struct {
		paymentIntentRepository repository.PaymentIntentRepository
	}
)

func NewHandlePaymentActionResultUseCase(paymentIntentRepository repository.PaymentIntentRepository) HandlePaymentActionResultUseCase {
	if paymentIntentRepository == nil {
		panic("paymentIntentRepository is nil")
	}
	return &handlePaymentActionResultUseCase{
		paymentIntentRepository: paymentIntentRepository,
	}
}

func (i HandlePaymentActionResultUseCaseInput) Validate() error {
	contract.AssertValidatable(i.PaymentIntentID)
	return nil
}

func (u *handlePaymentActionResultUseCase) Execute(ctx context.Context, input HandlePaymentActionResultUseCaseInput) (*HandlePaymentActionResultUseCaseOutput, error) {
	contract.AssertValidatable(input)

	paymentIntent, err := u.paymentIntentRepository.FindBy(ctx, input.PaymentIntentID)
	if err != nil {
		return nil, err
	}
	if paymentIntent == nil {
		// Webhook が先に届いた場合などは黙って終了
		return &HandlePaymentActionResultUseCaseOutput{
			PaymentIntentID: input.PaymentIntentID,
		}, nil
	}

	switch intent := (*paymentIntent).(type) {
	case domain.PaymentIntentRequiresAction:
		event, aggregate, err := intent.StartProcessing()
		if err != nil {
			return nil, err
		}

		if err := u.paymentIntentRepository.Save(ctx, event, aggregate); err != nil {
			return nil, err
		}

		return &HandlePaymentActionResultUseCaseOutput{
			PaymentIntentID: input.PaymentIntentID,
			PaymentIntent:   aggregate,
		}, nil
	default:
		// すでに進んでいる場合も成功として返す
		return &HandlePaymentActionResultUseCaseOutput{
			PaymentIntentID: input.PaymentIntentID,
			PaymentIntent:   intent,
		}, nil
	}
}
