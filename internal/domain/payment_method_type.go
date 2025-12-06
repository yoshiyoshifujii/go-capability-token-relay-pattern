package domain

import "errors"

type (
	PaymentMethodType string

	PaymentMethodTypes []PaymentMethodType
)

const (
	PaymentMethodTypeCard   PaymentMethodType = "card"
	PaymentMethodTypePayPay PaymentMethodType = "paypay"
)

func (p PaymentMethodType) Validate() error {
	if p == "" {
		return errors.New("payment method type is empty")
	}
	return nil
}

func (p PaymentMethodTypes) Validate() error {
	if len(p) == 0 {
		return errors.New("payment_method_types must not be empty")
	}
	return nil
}

func (p PaymentMethodTypes) Contains(methodType PaymentMethodType) bool {
	for _, t := range p {
		if t == methodType {
			return true
		}
	}
	return false
}
