package domain

import (
	"errors"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
)

type (
	PaymentIntentRequiresPaymentMethodType struct {
		paymentIntentMeta
		PaymentMethodTypes PaymentMethodTypes
		Amount             Money
	}

	PaymentIntentRequiresPaymentMethod struct {
		paymentIntentMeta
		PaymentMethodType PaymentMethodType
		Amount            Money
	}

	PaymentIntentRequiresConfirmation struct {
		paymentIntentMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
		Amount        Money
	}

	PaymentIntentRequiresAction struct {
		paymentIntentMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
		Amount        Money
	}

	PaymentIntentRequiresCapture struct {
		paymentIntentMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
		Amount        Money
	}

	PaymentIntentProcessing struct {
		paymentIntentMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
		Amount        Money
	}

	PaymentIntentSucceeded struct {
		paymentIntentMeta
		PaymentMethod PaymentMethod
		Amount        Money
	}
)

func GeneratePaymentIntent(id PaymentIntentID, types PaymentMethodTypes, amount Money) (PaymentIntentEvent, PaymentIntent, error) {
	seqNr := uint8(1)

	event := PaymentIntentRequiresPaymentMethodTypeEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: id,
			SeqNr:           seqNr,
		},
		PaymentMethodTypes: types,
		Amount:             amount,
	}

	aggregate := PaymentIntentRequiresPaymentMethodType{
		paymentIntentMeta: paymentIntentMeta{
			ID:     id,
			SeqNr:  seqNr,
			Amount: amount,
		},
		PaymentMethodTypes: types,
		Amount:             amount,
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
		Amount:            p.Amount,
	}

	aggregate := PaymentIntentRequiresPaymentMethod{
		paymentIntentMeta: paymentIntentMeta{
			ID:     p.ID,
			SeqNr:  seqNr,
			Amount: p.Amount,
		},
		PaymentMethodType: methodType,
		Amount:            p.Amount,
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
		Amount:        p.Amount,
	}

	aggregate := PaymentIntentRequiresConfirmation{
		paymentIntentMeta: paymentIntentMeta{
			ID:     p.ID,
			SeqNr:  seqNr,
			Amount: p.Amount,
		},
		PaymentMethod: method,
		CaptureMethod: captureMethod,
		Amount:        p.Amount,
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
		Amount:        p.Amount,
	}

	aggregate := PaymentIntentRequiresAction{
		paymentIntentMeta: paymentIntentMeta{
			ID:     p.ID,
			SeqNr:  seqNr,
			Amount: p.Amount,
		},
		PaymentMethod: p.PaymentMethod,
		CaptureMethod: p.CaptureMethod,
		Amount:        p.Amount,
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
		Amount:        p.Amount,
	}

	aggregate := PaymentIntentRequiresCapture{
		paymentIntentMeta: paymentIntentMeta{
			ID:     p.ID,
			SeqNr:  seqNr,
			Amount: p.Amount,
		},
		PaymentMethod: p.PaymentMethod,
		CaptureMethod: p.CaptureMethod,
		Amount:        p.Amount,
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
		Amount:        p.Amount,
	}

	aggregate := PaymentIntentProcessing{
		paymentIntentMeta: paymentIntentMeta{
			ID:     p.ID,
			SeqNr:  seqNr,
			Amount: p.Amount,
		},
		PaymentMethod: p.PaymentMethod,
		CaptureMethod: p.CaptureMethod,
		Amount:        p.Amount,
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
		Amount:        p.Amount,
	}

	aggregate := PaymentIntentRequiresCapture{
		paymentIntentMeta: paymentIntentMeta{
			ID:     p.ID,
			SeqNr:  seqNr,
			Amount: p.Amount,
		},
		PaymentMethod: p.PaymentMethod,
		CaptureMethod: p.CaptureMethod,
		Amount:        p.Amount,
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
			Amount:        p.Amount,
		}

		aggregate := PaymentIntentProcessing{
			paymentIntentMeta: paymentIntentMeta{
				ID:     p.ID,
				SeqNr:  seqNr,
				Amount: p.Amount,
			},
			PaymentMethod: p.PaymentMethod,
			CaptureMethod: p.CaptureMethod,
			Amount:        p.Amount,
		}

		return event, aggregate, nil
	default:
		panic("invalid capture method")
	}
}

func (p PaymentIntentRequiresCapture) StartProcessing() (PaymentIntentEvent, PaymentIntent, error) {
	contract.AssertValidatable(p.PaymentMethod)

	if p.CaptureMethod != PaymentCaptureMethodManual {
		return nil, nil, errors.New("capture method must be manual to start processing")
	}

	seqNr := p.SeqNr + 1

	event := PaymentIntentProcessingEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: p.ID,
			SeqNr:           seqNr,
		},
		PaymentMethod: p.PaymentMethod,
		CaptureMethod: p.CaptureMethod,
		Amount:        p.Amount,
	}

	aggregate := PaymentIntentProcessing{
		paymentIntentMeta: paymentIntentMeta{
			ID:     p.ID,
			SeqNr:  seqNr,
			Amount: p.Amount,
		},
		PaymentMethod: p.PaymentMethod,
		CaptureMethod: p.CaptureMethod,
		Amount:        p.Amount,
	}

	return event, aggregate, nil
}

func (p PaymentIntentProcessing) Complete() (PaymentIntentEvent, PaymentIntent, error) {
	contract.AssertValidatable(p.PaymentMethod)

	seqNr := p.SeqNr + 1

	event := PaymentIntentCompleteEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: p.ID,
			SeqNr:           seqNr,
		},
		PaymentMethod: p.PaymentMethod,
		Amount:        p.Amount,
	}

	aggregate := PaymentIntentSucceeded{
		paymentIntentMeta: paymentIntentMeta{
			ID:     p.ID,
			SeqNr:  seqNr,
			Amount: p.Amount,
		},
		PaymentMethod: p.PaymentMethod,
		Amount:        p.Amount,
	}

	return event, aggregate, nil
}
