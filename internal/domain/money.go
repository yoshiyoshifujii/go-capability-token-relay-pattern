package domain

import "errors"

type (
	Money uint8
)

func NewMoney(money uint8) Money {
	if money == 0 {
		panic("money cannot be zero")
	}
	return Money(money)
}

func (m Money) Validate() error {
	if m == 0 {
		return errors.New("money is zero")
	}
	return nil
}
