package domain

import "errors"

type (
	PaymentMethod struct {
		PaymentMethodType PaymentMethodType

		Card   *PaymentMethodCard
		PayPay *PaymentMethodPayPay
	}

	PaymentMethodCard struct {
		Number   string
		ExpYear  uint8
		ExpMonth uint8
	}

	PaymentMethodPayPay struct {
		AuthorizationURL string
	}
)

func NewPaymentMethod(
	paymentMethodType PaymentMethodType,
	card *PaymentMethodCard,
	payPay *PaymentMethodPayPay,
) PaymentMethod {
	method := PaymentMethod{
		PaymentMethodType: paymentMethodType,
		Card:              card,
		PayPay:            payPay,
	}
	if err := method.Validate(); err != nil {
		panic(err)
	}
	return method
}

func (p PaymentMethod) Validate() error {
	if err := p.PaymentMethodType.Validate(); err != nil {
		return err
	}

	switch p.PaymentMethodType {
	case PaymentMethodTypeCard:
		if p.Card == nil {
			return errors.New("card payment method requires card details")
		}
		return p.Card.Validate()
	case PaymentMethodTypePayPay:
		if p.PayPay == nil {
			return errors.New("paypay payment method requires paypay details")
		}
		return p.PayPay.Validate()
	default:
		return errors.New("unsupported payment method type")
	}
}

func (p PaymentMethodCard) Validate() error {
	if p.Number == "" {
		return errors.New("card number is empty")
	}
	if p.ExpYear == 0 {
		return errors.New("card exp year is empty")
	}
	if p.ExpMonth == 0 {
		return errors.New("card exp month is empty")
	}
	return nil
}

func (p PaymentMethodPayPay) Validate() error {
	if p.AuthorizationURL == "" {
		return errors.New("paypay authorization url is empty")
	}
	return nil
}
