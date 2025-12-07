package domain

import "errors"

type (
	PaymentIntentID string

	PaymentIntent interface {
		RequirePaymentMethod(PaymentMethodType) (PaymentIntentEvent, PaymentIntent, error)
		RequireConfirmation(PaymentMethod, PaymentCaptureMethod) (PaymentIntentEvent, PaymentIntent, error)
		RequireAction() (PaymentIntentEvent, PaymentIntent, error)
		RequireCapture() (PaymentIntentEvent, PaymentIntent, error)
		StartProcessing() (PaymentIntentEvent, PaymentIntent, error)
		Complete() (PaymentIntentEvent, PaymentIntent, error)
		Fail(PaymentFailureReason, bool) (PaymentIntentEvent, PaymentIntent, error)
	}

	paymentIntentMeta struct {
		ID      PaymentIntentID
		SeqNr   uint8
		Version uint8
		Amount  Money
	}
)

func (p PaymentIntentID) Validate() error {
	if len(p) == 0 {
		return errors.New("invalid payment intent id")
	}
	return nil
}

func (p paymentIntentMeta) RequirePaymentMethod(methodType PaymentMethodType) (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
}

func (p paymentIntentMeta) RequireConfirmation(method PaymentMethod, captureMethod PaymentCaptureMethod) (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
}

func (p paymentIntentMeta) RequireAction() (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
}

func (p paymentIntentMeta) RequireCapture() (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
}

func (p paymentIntentMeta) StartProcessing() (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
}

func (p paymentIntentMeta) Complete() (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
}

func (p paymentIntentMeta) Fail(reason PaymentFailureReason, retryable bool) (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
}
