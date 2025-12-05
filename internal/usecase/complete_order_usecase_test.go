package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_completeOrderUseCase_Execute(t *testing.T) {
	// given
	ctx := t.Context()

	input := CompleteOrderUseCaseInput{
		OrderProcessingID: "op_123",
		CapabilityTokens: []string{
			"coupon.token",
			"points.token",
			"payment.token",
		},
	}

	sut := NewCompleteOrderUseCase()

	// when
	actual, err := sut.Execute(ctx, input)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, actual)
}
