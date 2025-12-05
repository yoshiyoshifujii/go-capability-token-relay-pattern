package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_issuePointsTokenUseCase_Execute(t *testing.T) {
	// given
	ctx := t.Context()

	input := IssuePointsTokenUseCaseInput{
		OrderProcessingID: "op_123",
		UserID:            "user_123",
		PointsToUse:       100,
	}

	sut := NewIssuePointsTokenUseCase()

	// when
	actual, err := sut.Execute(ctx, input)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, actual)
}
