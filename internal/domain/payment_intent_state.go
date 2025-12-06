package domain

type (
	PaymentIntentRequiresPaymentMethodType struct {
		paymentIntentMeta
		PaymentMethodTypes PaymentMethodTypes
	}
)

func GeneratePaymentIntent(id PaymentIntentID, types PaymentMethodTypes) (PaymentIntentEvent, PaymentIntent, error) {
	seqNr := uint8(1)

	event := PaymentIntentRequiresPaymentMethodTypeEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: id,
			SeqNr:           seqNr,
		},
		PaymentMethodTypes: types,
	}

	aggregate := PaymentIntentRequiresPaymentMethodType{
		paymentIntentMeta: paymentIntentMeta{
			ID:    id,
			SeqNr: seqNr,
		},
		PaymentMethodTypes: types,
	}

	return event, aggregate, nil
}
