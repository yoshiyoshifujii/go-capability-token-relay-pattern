package domain

import "errors"

type (
	PaymentMethodType string

	PaymentMethodTypes []PaymentMethodType
)

const (
	PaymentMethodTypeCard PaymentMethodType = "card"
)

func (p PaymentMethodTypes) Validate() error {
	if len(p) == 0 {
		return errors.New("payment_method_types must not be empty")
	}
	return nil
}
