package service

import (
	"context"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
)

type (
	IDGenerator[ID any] interface {
		GenerateID(context.Context) (ID, error)
	}

	BusinessIDGenerator interface {
		IDGenerator[domain.BusinessID]
	}

	CartIDGenerator interface {
		IDGenerator[domain.CartID]
	}
)
