package poll

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test  for poller manager
func TestPollerManager(t *testing.T) {
	var err error
	err = PollerManager.SetPollerNums(0)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "SetPollerNums(n int):@n < 0")

	err = PollerManager.SetPollerNums(10)
	assert.Nil(t, err)

	err = PollerManager.SetPollerNums(2)
	assert.Nil(t, err)

	poller := PollerManager.Pick()
	assert.NotNil(t, poller)

	err = PollerManager.Close()
	assert.Nil(t, err)
}
