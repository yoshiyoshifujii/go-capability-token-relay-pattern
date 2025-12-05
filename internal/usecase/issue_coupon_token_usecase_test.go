package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_issueCouponTokenUseCase_Execute(t *testing.T) {
	// given
	ctx := t.Context()

	input := IssueCouponTokenUseCaseInput{
		OrderProcessingID: "op_123",
		UserID:            "user_123",
		CouponRef:         "coupon_abc",
	}

	sut := NewIssueCouponTokenUseCase()

	// when
	actual, err := sut.Execute(ctx, input)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, actual)
}
