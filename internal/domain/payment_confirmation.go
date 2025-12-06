package domain

type PaymentConfirmationNext string

const (
	PaymentConfirmationNextProcessing      PaymentConfirmationNext = "processing"
	PaymentConfirmationNextRequiresAction  PaymentConfirmationNext = "requires_action"
	PaymentConfirmationNextRequiresCapture PaymentConfirmationNext = "requires_capture"
)
