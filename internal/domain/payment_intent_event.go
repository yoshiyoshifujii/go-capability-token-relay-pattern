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
	}

	PaymentIntentRequiresPaymentMethodEvent struct {
		paymentIntentEventMeta
		PaymentMethodType PaymentMethodType
	}

	PaymentIntentRequiresConfirmationEvent struct {
		paymentIntentEventMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
	}

	PaymentIntentRequiresActionEvent struct {
		paymentIntentEventMeta
		PaymentMethod PaymentMethod
	}

	PaymentIntentRequiresCaptureEvent struct {
		paymentIntentEventMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
	}

	PaymentIntentProcessingEvent struct {
		paymentIntentEventMeta
		PaymentMethod PaymentMethod
		CaptureMethod PaymentCaptureMethod
	}

	PaymentIntentCompleteEvent struct {
		paymentIntentEventMeta
		PaymentMethod PaymentMethod
	}
)

func (e paymentIntentEventMeta) PaymentIntentEvent() {
	panic("do not call this method")
}
