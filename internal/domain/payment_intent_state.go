package domain

import "errors"

type (
	PaymentIntentRequiresPaymentMethodType struct {
		paymentIntentMeta
		PaymentMethodTypes PaymentMethodTypes
	}

	PaymentIntentRequiresPaymentMethod struct {
		paymentIntentMeta
		PaymentMethodType PaymentMethodType
	}

	PaymentIntentRequiresConfirmation struct {
		paymentIntentMeta
		PaymentMethod PaymentMethod
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

func (p PaymentIntentRequiresPaymentMethodType) RequirePaymentMethod(methodType PaymentMethodType) (PaymentIntentEvent, PaymentIntent, error) {
	if !p.PaymentMethodTypes.Contains(methodType) {
		return nil, nil, errors.New("payment method type is not allowed")
	}

	seqNr := p.SeqNr + 1

	event := PaymentIntentRequiresPaymentMethodEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: p.ID,
			SeqNr:           seqNr,
		},
		PaymentMethodType: methodType,
	}

	aggregate := PaymentIntentRequiresPaymentMethod{
		paymentIntentMeta: paymentIntentMeta{
			ID:    p.ID,
			SeqNr: seqNr,
		},
		PaymentMethodType: methodType,
	}

	return event, aggregate, nil
}

func (p PaymentIntentRequiresPaymentMethod) RequireConfirmation(method PaymentMethod) (PaymentIntentEvent, PaymentIntent, error) {
	if method.PaymentMethodType != p.PaymentMethodType {
		return nil, nil, errors.New("payment method type mismatch")
	}

	switch method.PaymentMethodType {
	case PaymentMethodTypeCard:
		if method.Card == nil {
			return nil, nil, errors.New("card payment method requires card details")
		}
	case PaymentMethodTypePayPay:
		if method.PayPay == nil {
			return nil, nil, errors.New("paypay payment method requires paypay details")
		}
	default:
		return nil, nil, errors.New("unsupported payment method type")
	}

	seqNr := p.SeqNr + 1

	event := PaymentIntentRequiresConfirmationEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: p.ID,
			SeqNr:           seqNr,
		},
		PaymentMethod: method,
	}

	aggregate := PaymentIntentRequiresConfirmation{
		paymentIntentMeta: paymentIntentMeta{
			ID:    p.ID,
			SeqNr: seqNr,
		},
		PaymentMethod: method,
	}

	return event, aggregate, nil
}
