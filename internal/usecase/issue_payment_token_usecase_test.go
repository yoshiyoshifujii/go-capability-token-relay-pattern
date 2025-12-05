package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_issuePaymentTokenUseCase_Execute(t *testing.T) {
	// given
	ctx := t.Context()

	input := IssuePaymentTokenUseCaseInput{
		OrderProcessingID: "op_123",
		UserID:            "user_123",
		PaymentMethod:     "credit-card",
	}

	sut := NewIssuePaymentTokenUseCase()

	// when
	actual, err := sut.Execute(ctx, input)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, actual)
}
