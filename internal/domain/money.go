package domain

type (
	Money uint8
)

func NewMoney(money uint8) Money {
	if money == 0 {
		panic("money cannot be zero")
	}
	return Money(money)
}
