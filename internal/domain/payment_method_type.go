package domain

type (
	PaymentMethodType string

	PaymentMethodTypes []PaymentMethodType
)

const (
	PaymentMethodTypeCard PaymentMethodType = "card"
)
