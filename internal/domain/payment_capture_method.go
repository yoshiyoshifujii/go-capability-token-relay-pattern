package domain

import "errors"

type PaymentCaptureMethod string

const (
	PaymentCaptureMethodAutomatic PaymentCaptureMethod = "automatic"
	PaymentCaptureMethodManual    PaymentCaptureMethod = "manual"
)

func (p PaymentCaptureMethod) Validate() error {
	if len(p) == 0 {
		return errors.New("payment capture method is empty")
	}

	switch p {
	case PaymentCaptureMethodAutomatic, PaymentCaptureMethodManual:
		return nil
	default:
		return errors.New("unsupported payment capture method")
	}
}
