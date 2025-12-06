package service

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type FakePaymentIntentIDGenerator struct {
	NextID domain.PaymentIntentID
	Err    error
}

func NewFakePaymentIntentIDGenerator(nextID domain.PaymentIntentID) service.PaymentIDGenerator {
	return &FakePaymentIntentIDGenerator{NextID: nextID}
}

func (f *FakePaymentIntentIDGenerator) GenerateID(ctx context.Context) (domain.PaymentIntentID, error) {
	if f.Err != nil {
		return "", f.Err
	}
	return f.NextID, nil
}
