package repository

import (
	"context"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type (
	Repository[AggregateID, Aggregate, Event any] interface {
		FindBy(ctx context.Context, aggregateID AggregateID) (*Aggregate, error)
		Save(ctx context.Context, event Event, aggregate Aggregate) error
	}

	BusinessRepository interface {
		Repository[domain.BusinessID, domain.Business, domain.BusinessEvent]
	}
)
