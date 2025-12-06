package service

import (
	"context"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type FakeCartIDGenerator struct {
	NextID domain.CartID
	Err    error
}

func NewFakeCartIDGenerator(nextID domain.CartID) service.CartIDGenerator {
	return &FakeCartIDGenerator{NextID: nextID}
}

func (f *FakeCartIDGenerator) GenerateID(ctx context.Context) (domain.CartID, error) {
	if f.Err != nil {
		return "", f.Err
	}
	return f.NextID, nil
}
