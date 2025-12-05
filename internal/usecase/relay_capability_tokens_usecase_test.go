package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_relayCapabilityTokensUseCase_Execute(t *testing.T) {
	// given
	ctx := t.Context()

	input := RelayCapabilityTokensUseCaseInput{
		OrderProcessingID: "op_123",
		CouponToken:       "coupon.token",
		PointsToken:       "points.token",
		PaymentToken:      "payment.token",
	}

	sut := NewRelayCapabilityTokensUseCase()

	// when
	actual, err := sut.Execute(ctx, input)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.NotNil(t, actual.VerifiedTokens)
}
