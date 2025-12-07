package domain

type (
	PaymentIntentEvent interface {
		PaymentIntentEvent()
	}

	paymentIntentEventMeta struct {
		PaymentIntentID PaymentIntentID
		SeqNr           uint8
	}

	PaymentIntentRequiresPaymentMethodTypeEvent struct {
		paymentIntentEventMeta
		PaymentMethodTypes PaymentMethodTypes
		Amount             Money
	}

	PaymentIntentRequiresPaymentMethodEvent struct {
		paymentIntentEventMeta
		PaymentMethodType PaymentMethodType
		Amount            Money
	}

	PaymentIntentRequiresConfirmationEvent struct {
		paymentIntentEventMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
		Amount        Money
	}

	PaymentIntentRequiresActionEvent struct {
		paymentIntentEventMeta
		PaymentMethod PaymentMethod
		Amount        Money
	}

	PaymentIntentRequiresCaptureEvent struct {
		paymentIntentEventMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
		Amount        Money
	}

	PaymentIntentProcessingEvent struct {
		paymentIntentEventMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
		Amount        Money
	}

	PaymentIntentCompleteEvent struct {
		paymentIntentEventMeta
		PaymentMethod PaymentMethod
		Amount        Money
	}

	PaymentIntentFailedEvent struct {
		paymentIntentEventMeta
		PaymentMethodType PaymentMethodType
		PaymentMethod     PaymentMethod
		Amount            Money
		Reason            PaymentFailureReason
	}

	PaymentIntentCanceledEvent struct {
		paymentIntentEventMeta
		PaymentMethod PaymentMethod
		Amount        Money
		Reason        PaymentFailureReason
	}
)

func (e paymentIntentEventMeta) PaymentIntentEvent() {
	panic("do not call this method")
}
