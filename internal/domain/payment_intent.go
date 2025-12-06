package domain

import "errors"

type (
	PaymentIntentID string

	PaymentIntent interface {
		RequirePaymentMethod(PaymentMethodType) (PaymentIntentEvent, PaymentIntent, error)
		RequireConfirmation(PaymentMethod) (PaymentIntentEvent, PaymentIntent, error)
		RequireAction(PaymentMethod) (PaymentIntentEvent, PaymentIntent, error)
		RequireCapture(PaymentMethod) (PaymentIntentEvent, PaymentIntent, error)
		StartProcessing(PaymentMethod) (PaymentIntentEvent, PaymentIntent, error)
		Complete(PaymentMethod) (PaymentIntentEvent, PaymentIntent, error)
	}

	paymentIntentMeta struct {
		ID      PaymentIntentID
		SeqNr   uint8
		Version uint8
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

func (p paymentIntentMeta) RequireConfirmation(method PaymentMethod) (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
}

func (p paymentIntentMeta) RequireAction(method PaymentMethod) (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
}

func (p paymentIntentMeta) RequireCapture(method PaymentMethod) (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
}

func (p paymentIntentMeta) StartProcessing(method PaymentMethod) (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
}

func (p paymentIntentMeta) Complete(method PaymentMethod) (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
}
