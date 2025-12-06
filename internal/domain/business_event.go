package domain

type (
	BusinessEvent interface {
		BusinessEvent()
	}

	businessEventMeta struct {
		BusinessID BusinessID
		SeqNr      uint64
	}

	BusinessInitializedEvent struct {
		businessEventMeta
		BusinessName string
	}
)

func (b *businessEventMeta) BusinessEvent() {
	panic("do not call this method")
}

func NewBusinessInitializedEvent(
	businessID BusinessID,
	seqNr uint64,
	businessName string,
) BusinessInitializedEvent {
	return BusinessInitializedEvent{
		businessEventMeta: businessEventMeta{
			BusinessID: businessID,
			SeqNr:      seqNr,
		},
		BusinessName: businessName,
	}
}

func (b BusinessInitializedEvent) BusinessEvent() {}
