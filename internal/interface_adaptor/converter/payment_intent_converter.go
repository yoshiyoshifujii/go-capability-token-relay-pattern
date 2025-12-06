package converter

import (
	"fmt"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type PaymentIntentView struct {
	ID                 domain.PaymentIntentID
	SeqNr              uint8
	Status             string
	PaymentMethodTypes domain.PaymentMethodTypes
	PaymentMethodType  domain.PaymentMethodType
	PaymentMethod      domain.PaymentMethod
}

func ToPaymentIntentView(intent domain.PaymentIntent) (PaymentIntentView, error) {
	switch v := intent.(type) {
	case domain.PaymentIntentRequiresPaymentMethodType:
		return PaymentIntentView{
			ID:                 v.ID,
			SeqNr:              v.SeqNr,
			Status:             "requires_payment_method_type",
			PaymentMethodTypes: v.PaymentMethodTypes,
		}, nil
	case domain.PaymentIntentRequiresPaymentMethod:
		return PaymentIntentView{
			ID:                v.ID,
			SeqNr:             v.SeqNr,
			Status:            "requires_payment_method",
			PaymentMethodType: v.PaymentMethodType,
		}, nil
	case domain.PaymentIntentRequiresConfirmation:
		return PaymentIntentView{
			ID:            v.ID,
			SeqNr:         v.SeqNr,
			Status:        "requires_confirmation",
			PaymentMethod: v.PaymentMethod,
		}, nil
	case domain.PaymentIntentRequiresCapture:
		return PaymentIntentView{
			ID:            v.ID,
			SeqNr:         v.SeqNr,
			Status:        "requires_capture",
			PaymentMethod: v.PaymentMethod,
		}, nil
	case domain.PaymentIntentProcessing:
		return PaymentIntentView{
			ID:            v.ID,
			SeqNr:         v.SeqNr,
			Status:        "processing",
			PaymentMethod: v.PaymentMethod,
		}, nil
	default:
		return PaymentIntentView{}, fmt.Errorf("unsupported payment intent state %T", intent)
	}
}
