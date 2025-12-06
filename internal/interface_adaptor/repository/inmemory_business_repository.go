package repository

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type InMemoryBusinessRepository struct {
	store *InMemoryEventStore[domain.Business, domain.BusinessEvent]
}

func NewInMemoryBusinessRepository() *InMemoryBusinessRepository {
	return &InMemoryBusinessRepository{
		store: NewInMemoryEventStore[domain.Business, domain.BusinessEvent](),
	}
}

func (i *InMemoryBusinessRepository) FindBy(ctx context.Context, aggregateID domain.BusinessID) (*domain.Business, error) {
	for _, entity := range i.store.Entities {
		if entity.ID == aggregateID {
			e := entity
			return &e, nil
		}
	}
	return nil, nil
}

func (i *InMemoryBusinessRepository) Save(ctx context.Context, event domain.BusinessEvent, aggregate domain.Business) error {
	i.store.Events = append(i.store.Events, event)
	i.store.Entities = append(i.store.Entities, aggregate)
	return nil
}

func (i *InMemoryBusinessRepository) Events() []domain.BusinessEvent {
	return i.store.Events
}
