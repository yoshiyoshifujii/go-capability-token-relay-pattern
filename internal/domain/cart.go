package domain

import (
	"errors"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
)

type (
	CartID string

	CartItems []ItemID

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

func NewCartItems(items ...ItemID) CartItems {
	if len(items) == 0 {
		panic("invalid cartItems")
	}
	return CartItems(items)
}

func (c CartItems) Validate() error {
	if len(c) == 0 {
		return errors.New("invalid cartItems")
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

func (c *Cart) Validate() error {
	if c == nil {
		return errors.New("invalid cart")
	}
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
