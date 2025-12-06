package repository

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type InMemoryPaymentIntentRepository struct {
	store *InMemoryEventStore[domain.PaymentIntent, domain.PaymentIntentEvent]
}

func NewInMemoryPaymentIntentRepository() *InMemoryPaymentIntentRepository {
	return &InMemoryPaymentIntentRepository{
		store: NewInMemoryEventStore[domain.PaymentIntent, domain.PaymentIntentEvent](),
	}
}

func (i *InMemoryPaymentIntentRepository) FindBy(ctx context.Context, aggregateID domain.PaymentIntentID) (*domain.PaymentIntent, error) {
	for _, entity := range i.store.Entities {
		if paymentIntentID(entity) == aggregateID {
			e := entity
			return &e, nil
		}
	}
	return nil, nil
}

func (i *InMemoryPaymentIntentRepository) Save(ctx context.Context, event domain.PaymentIntentEvent, aggregate domain.PaymentIntent) error {
	i.store.Events = append(i.store.Events, event)
	i.store.Entities = append(i.store.Entities, aggregate)
	return nil
}

func (i *InMemoryPaymentIntentRepository) Events() []domain.PaymentIntentEvent {
	return i.store.Events
}

func paymentIntentID(intent domain.PaymentIntent) domain.PaymentIntentID {
	switch v := intent.(type) {
	case domain.PaymentIntentInitialized:
		return v.ID
	case domain.PaymentIntentRequiresPaymentMethodType:
		return v.ID
	default:
		panic("unsupported payment intent aggregate for in-memory repository")
	}
}
