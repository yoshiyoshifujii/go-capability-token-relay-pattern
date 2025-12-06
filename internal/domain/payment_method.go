package domain

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
