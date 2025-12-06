package domain

type (
	ItemID string

	ItemName string

	ItemPrice Money

	Item struct {
		BusinessID BusinessID
		ID         ItemID
		Name       ItemName
		Price      ItemPrice
	}
)

func NewItem(
	businessID BusinessID,
	id ItemID,
	name ItemName,
	price ItemPrice,
) Item {
	return Item{
		ID:         id,
		BusinessID: businessID,
		Name:       name,
		Price:      price,
	}
}
