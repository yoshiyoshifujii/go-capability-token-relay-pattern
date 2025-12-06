package domain

import (
	"errors"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
)

type (
	CartID string

	CartItem struct {
		ItemID ItemID
		Price  ItemPrice
	}

	CartItems []CartItem

	Cart struct {
		BusinessID BusinessID
		CartID     CartID
		Items      CartItems
	}
)

func NewCartID(id string) CartID {
	if len(id) == 0 {
		panic("invalid cartID")
	}
	return CartID(id)
}

func (c CartID) Validate() error {
	if len(c) == 0 {
		return errors.New("invalid cart id")
	}
	return nil
}

func (c CartItem) Validate() error {
	if err := c.ItemID.Validate(); err != nil {
		return err
	}
	return Money(c.Price).Validate()
}

func NewCartItems(items ...CartItem) CartItems {
	if len(items) == 0 {
		panic("invalid cartItems")
	}
	return CartItems(items)
}

func (c CartItems) Validate() error {
	if len(c) == 0 {
		return errors.New("invalid cartItems")
	}
	for _, item := range c {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func NewCart(
	businessID BusinessID,
	cartID CartID,
	items CartItems,
) Cart {
	contract.AssertValidatable(businessID)
	contract.AssertValidatable(cartID)
	contract.AssertValidatable(items)

	return Cart{
		BusinessID: businessID,
		CartID:     cartID,
		Items:      items,
	}
}

func (c Cart) Validate() error {
	if err := c.BusinessID.Validate(); err != nil {
		return err
	}
	if err := c.CartID.Validate(); err != nil {
		return err
	}
	if err := c.Items.Validate(); err != nil {
		return err
	}
	return nil
}

func (c Cart) CalculateAmount() Money {
	var total Money
	for _, item := range c.Items {
		total += Money(item.Price)
	}
	return total
}
