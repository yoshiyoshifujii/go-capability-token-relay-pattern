package domain

type (
	PaymentIntentInitialized struct {
		paymentIntentMeta
	}

	PaymentIntentRequiresPaymentMethodType struct {
		paymentIntentMeta
		PaymentMethodTypes PaymentMethodTypes
	}
)

func (p PaymentIntentInitialized) RequirePaymentMethodTypes(types PaymentMethodTypes) (PaymentIntentEvent, PaymentIntent, error) {
	seqNr := p.SeqNr + 1

	event := PaymentIntentRequiresPaymentMethodTypeEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: p.ID,
			SeqNr:           seqNr,
		},
		PaymentMethodTypes: types,
	}

	aggregate := PaymentIntentRequiresPaymentMethodType{
		paymentIntentMeta: paymentIntentMeta{
			ID:    p.ID,
			SeqNr: seqNr,
		},
		PaymentMethodTypes: types,
	}

	return event, aggregate, nil
}
