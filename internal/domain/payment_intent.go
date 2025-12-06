package domain

type (
	PaymentIntentID string

	PaymentIntent interface {
		RequirePaymentMethodTypes(PaymentMethodTypes) (PaymentIntentEvent, PaymentIntent, error)
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

func (p paymentIntentMeta) RequirePaymentMethodTypes(types PaymentMethodTypes) (PaymentIntentEvent, PaymentIntent, error) {
	panic("intentionally unimplemented; concrete states must override as needed")
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
