package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	num := Max(1, 0)
	assert.Equal(t, 1, num)

	num2 := Max(2, 3)
	assert.Equal(t, 3, num2)
}

func TestMin(t *testing.T) {
	num := Min(1, 0)
	assert.Equal(t, 0, num)

	num2 := Min(2, 3)
	assert.Equal(t, 2, num2)
}
