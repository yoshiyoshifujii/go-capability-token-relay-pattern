package repository

import (
	"context"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type InMemoryBusinessRepository struct {
	Events   []domain.BusinessEvent
	Entities []domain.Business
}

func NewInMemoryBusinessRepository() *InMemoryBusinessRepository {
	return &InMemoryBusinessRepository{
		Events:   make([]domain.BusinessEvent, 0),
		Entities: make([]domain.Business, 0),
	}
}

func (i *InMemoryBusinessRepository) FindBy(ctx context.Context, aggregateID domain.BusinessID) (*domain.Business, error) {
	for _, entity := range i.Entities {
		if entity.ID == aggregateID {
			return &entity, nil
		}
	}
	return nil, nil
}

func (i *InMemoryBusinessRepository) Save(ctx context.Context, event domain.BusinessEvent, aggregate domain.Business) error {
	i.Events = append(i.Events, event)
	i.Entities = append(i.Entities, aggregate)
	return nil
}
