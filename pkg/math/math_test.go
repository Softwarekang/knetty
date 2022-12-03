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
