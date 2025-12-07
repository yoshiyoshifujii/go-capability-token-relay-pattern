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
		FailureReason     PaymentFailureReason
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

	PaymentIntentCanceled struct {
		paymentIntentMeta
		PaymentMethod PaymentMethod
		Amount        Money
		FailureReason PaymentFailureReason
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

	if p.CaptureMethod != PaymentCaptureMethodAutomatic {
		return nil, nil, errors.New("capture method must be automatic to start processing after action")
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

func (p PaymentIntentRequiresConfirmation) Fail(reason PaymentFailureReason, retryable bool) (PaymentIntentEvent, PaymentIntent, error) {
	return failPaymentIntent(p.paymentIntentMeta, p.PaymentMethod, p.Amount, reason, retryable)
}

func (p PaymentIntentRequiresAction) Fail(reason PaymentFailureReason, retryable bool) (PaymentIntentEvent, PaymentIntent, error) {
	return failPaymentIntent(p.paymentIntentMeta, p.PaymentMethod, p.Amount, reason, retryable)
}

func (p PaymentIntentRequiresCapture) Fail(reason PaymentFailureReason, retryable bool) (PaymentIntentEvent, PaymentIntent, error) {
	return failPaymentIntent(p.paymentIntentMeta, p.PaymentMethod, p.Amount, reason, retryable)
}

func (p PaymentIntentProcessing) Fail(reason PaymentFailureReason, retryable bool) (PaymentIntentEvent, PaymentIntent, error) {
	return failPaymentIntent(p.paymentIntentMeta, p.PaymentMethod, p.Amount, reason, retryable)
}

func failPaymentIntent(
	meta paymentIntentMeta,
	paymentMethod PaymentMethod,
	amount Money,
	reason PaymentFailureReason,
	retryable bool,
) (PaymentIntentEvent, PaymentIntent, error) {
	contract.AssertValidatable(paymentMethod)
	contract.AssertValidatable(reason)

	seqNr := meta.SeqNr + 1

	if retryable {
		event := PaymentIntentFailedEvent{
			paymentIntentEventMeta: paymentIntentEventMeta{
				PaymentIntentID: meta.ID,
				SeqNr:           seqNr,
			},
			PaymentMethodType: paymentMethod.PaymentMethodType,
			PaymentMethod:     paymentMethod,
			Amount:            amount,
			Reason:            reason,
		}

		aggregate := PaymentIntentRequiresPaymentMethod{
			paymentIntentMeta: paymentIntentMeta{
				ID:     meta.ID,
				SeqNr:  seqNr,
				Amount: amount,
			},
			PaymentMethodType: paymentMethod.PaymentMethodType,
			Amount:            amount,
			FailureReason:     reason,
		}

		return event, aggregate, nil
	}

	event := PaymentIntentCanceledEvent{
		paymentIntentEventMeta: paymentIntentEventMeta{
			PaymentIntentID: meta.ID,
			SeqNr:           seqNr,
		},
		PaymentMethod: paymentMethod,
		Amount:        amount,
		Reason:        reason,
	}

	aggregate := PaymentIntentCanceled{
		paymentIntentMeta: paymentIntentMeta{
			ID:     meta.ID,
			SeqNr:  seqNr,
			Amount: amount,
		},
		PaymentMethod: paymentMethod,
		Amount:        amount,
		FailureReason: reason,
	}

	return event, aggregate, nil
}
