package domain

type (
	PaymentMethod struct {
		PaymentMethodType PaymentMethodType

		Card *PaymentMethodCard
	}

	PaymentMethodCard struct {
		Number   string
		ExpYear  uint8
		ExpMonth uint8
	}
)
