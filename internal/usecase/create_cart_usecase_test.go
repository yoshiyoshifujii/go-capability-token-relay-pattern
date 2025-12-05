package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createCartUseCase_Execute(t *testing.T) {
	// given
	ctx := t.Context()

	input := CreateCartUseCaseInput{}

	sut := NewCreateCartUseCase()

	// when
	actual, err := sut.Execute(ctx, input)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, actual)
}
