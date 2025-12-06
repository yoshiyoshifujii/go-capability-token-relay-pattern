package domain

import (
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
)

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
		CaptureMethod PaymentCaptureMethod
	}

	PaymentIntentRequiresAction struct {
		paymentIntentMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
	}

	PaymentIntentRequiresCapture struct {
		paymentIntentMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
	}

	PaymentIntentProcessing struct {
		paymentIntentMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
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
		panic("payment method type not found in payment methods")
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

func (p PaymentIntentRequiresPaymentMethod) RequireConfirmation(method PaymentMethod, captureMethod PaymentCaptureMethod) (PaymentIntentEvent, PaymentIntent, error) {
	contract.AssertValidatable(method)
	contract.AssertValidatable(captureMethod)

	if method.PaymentMethodType != p.PaymentMethodType {
		panic("payment method type is not allowed")
	}

	seqNr := p.SeqNr + 1

	event := PaymentIntentRequiresConfirmationEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: p.ID,
			SeqNr:           seqNr,
		},
		PaymentMethod: method,
		CaptureMethod: captureMethod,
	}

	aggregate := PaymentIntentRequiresConfirmation{
		paymentIntentMeta: paymentIntentMeta{
			ID:    p.ID,
			SeqNr: seqNr,
		},
		PaymentMethod: method,
		CaptureMethod: captureMethod,
	}

	return event, aggregate, nil
}

func (p PaymentIntentRequiresConfirmation) RequireAction() (PaymentIntentEvent, PaymentIntent, error) {
	contract.AssertValidatable(p.PaymentMethod)

	seqNr := p.SeqNr + 1

	event := PaymentIntentRequiresActionEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: p.ID,
			SeqNr:           seqNr,
		},
		PaymentMethod: p.PaymentMethod,
	}

	aggregate := PaymentIntentRequiresAction{
		paymentIntentMeta: paymentIntentMeta{
			ID:    p.ID,
			SeqNr: seqNr,
		},
		PaymentMethod: p.PaymentMethod,
		CaptureMethod: p.CaptureMethod,
	}

	return event, aggregate, nil
}

func (p PaymentIntentRequiresConfirmation) RequireCapture() (PaymentIntentEvent, PaymentIntent, error) {
	contract.AssertValidatable(p.PaymentMethod)

	seqNr := p.SeqNr + 1

	event := PaymentIntentRequiresCaptureEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: p.ID,
			SeqNr:           seqNr,
		},
		PaymentMethod: p.PaymentMethod,
	}

	aggregate := PaymentIntentRequiresCapture{
		paymentIntentMeta: paymentIntentMeta{
			ID:    p.ID,
			SeqNr: seqNr,
		},
		PaymentMethod: p.PaymentMethod,
		CaptureMethod: p.CaptureMethod,
	}

	return event, aggregate, nil
}

func (p PaymentIntentRequiresConfirmation) StartProcessing() (PaymentIntentEvent, PaymentIntent, error) {
	contract.AssertValidatable(p.PaymentMethod)

	seqNr := p.SeqNr + 1

	event := PaymentIntentProcessingEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: p.ID,
			SeqNr:           seqNr,
		},
		PaymentMethod: p.PaymentMethod,
	}

	aggregate := PaymentIntentProcessing{
		paymentIntentMeta: paymentIntentMeta{
			ID:    p.ID,
			SeqNr: seqNr,
		},
		PaymentMethod: p.PaymentMethod,
		CaptureMethod: p.CaptureMethod,
	}

	return event, aggregate, nil
}
