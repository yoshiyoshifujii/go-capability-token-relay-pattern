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
		BusinessName       string
		PaymentMethodTypes PaymentMethodTypes
	}
)

func (b *businessEventMeta) BusinessEvent() {
	panic("do not call this method")
}

func NewBusinessInitializedEvent(
	businessID BusinessID,
	seqNr uint64,
	businessName string,
	paymentMethodTypes PaymentMethodTypes,
) BusinessInitializedEvent {
	return BusinessInitializedEvent{
		businessEventMeta: businessEventMeta{
			BusinessID: businessID,
			SeqNr:      seqNr,
		},
		BusinessName:       businessName,
		PaymentMethodTypes: paymentMethodTypes,
	}
}

func (b BusinessInitializedEvent) BusinessEvent() {}
