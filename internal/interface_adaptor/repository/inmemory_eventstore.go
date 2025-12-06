package repository

// InMemoryEventStore keeps aggregates and their events in memory so they can be shared across repositories.
type InMemoryEventStore[Aggregate any, Event any] struct {
	Events   []Event
	Entities []Aggregate
}

func NewInMemoryEventStore[Aggregate any, Event any]() *InMemoryEventStore[Aggregate, Event] {
	return &InMemoryEventStore[Aggregate, Event]{
		Events:   make([]Event, 0),
		Entities: make([]Aggregate, 0),
	}
}
