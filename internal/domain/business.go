package domain

import (
	"errors"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/lib/contract"
)

type (
	BusinessID string

	Business struct {
		ID   BusinessID
		Name string
	}
)

func NewBusinessID(id string) BusinessID {
	if len(id) == 0 {
		panic("invalid business id")
	}
	return BusinessID(id)
}

func (b BusinessID) Validate() error {
	if len(b) == 0 {
		return errors.New("invalid business id")
	}
	return nil
}

func NewBusiness(id BusinessID, name string) Business {
	contract.AssertValidatable(id)
	if len(name) == 0 {
		panic("invalid business name")
	}
	return Business{
		ID:   id,
		Name: name,
	}
}
