package domain

import (
	"errors"

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
		return nil, nil, errors.New("payment method type not found in payment methods")
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
		return nil, nil, errors.New("payment method type is not allowed")
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

	if p.CaptureMethod != PaymentCaptureMethodManual {
		return nil, nil, errors.New("capture method must be manual to require capture after confirmation")
	}

	seqNr := p.SeqNr + 1

	event := PaymentIntentRequiresCaptureEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: p.ID,
			SeqNr:           seqNr,
		},
		PaymentMethod: p.PaymentMethod,
		CaptureMethod: p.CaptureMethod,
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
		CaptureMethod: p.CaptureMethod,
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

func (p PaymentIntentRequiresConfirmation) ApplyConfirmationResult(next PaymentConfirmationNext) (PaymentIntentEvent, PaymentIntent, error) {
	switch next {
	case PaymentConfirmationNextProcessing:
		return p.StartProcessing()
	case PaymentConfirmationNextRequiresAction:
		return p.RequireAction()
	case PaymentConfirmationNextRequiresCapture:
		return p.RequireCapture()
	default:
		panic("invalid payment confirmation next")
	}
}

func (p PaymentIntentRequiresAction) RequireCapture() (PaymentIntentEvent, PaymentIntent, error) {
	contract.AssertValidatable(p.PaymentMethod)

	if p.CaptureMethod != PaymentCaptureMethodManual {
		return nil, nil, errors.New("capture method must be manual to require capture after action")
	}

	seqNr := p.SeqNr + 1

	event := PaymentIntentRequiresCaptureEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: p.ID,
			SeqNr:           seqNr,
		},
		PaymentMethod: p.PaymentMethod,
		CaptureMethod: p.CaptureMethod,
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

func (p PaymentIntentRequiresAction) StartProcessing() (PaymentIntentEvent, PaymentIntent, error) {
	contract.AssertValidatable(p.PaymentMethod)

	switch p.CaptureMethod {
	case PaymentCaptureMethodManual:
		return p.RequireCapture()
	case PaymentCaptureMethodAutomatic:
		seqNr := p.SeqNr + 1

		event := PaymentIntentProcessingEvent{
			paymentIntentEventMeta: paymentIntentEventMeta{
				PaymentIntentID: p.ID,
				SeqNr:           seqNr,
			},
			PaymentMethod: p.PaymentMethod,
			CaptureMethod: p.CaptureMethod,
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
	default:
		panic("invalid capture method")
	}
}
