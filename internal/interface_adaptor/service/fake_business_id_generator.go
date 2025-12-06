package service

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type FakeBusinessIDGenerator struct {
	NextID domain.BusinessID
	Err    error
}

func NewFakeBusinessIDGenerator(nextID domain.BusinessID) *FakeBusinessIDGenerator {
	return &FakeBusinessIDGenerator{NextID: nextID}
}

func (f *FakeBusinessIDGenerator) GenerateID(ctx context.Context) (domain.BusinessID, error) {
	if f.Err != nil {
		return "", f.Err
	}
	return f.NextID, nil
}
