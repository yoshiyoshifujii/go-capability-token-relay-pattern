package domain

import "errors"

type PaymentFailureReason string

const (
	PaymentFailureReasonConfirmationFailed PaymentFailureReason = "confirmation_failed"
	PaymentFailureReasonCaptureFailed      PaymentFailureReason = "capture_failed"
)

func (p PaymentFailureReason) Validate() error {
	if p == "" {
		return errors.New("payment failure reason is empty")
	}
	return nil
}
