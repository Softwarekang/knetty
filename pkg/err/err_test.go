package err

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestknettyError(t *testing.T) {
	var (
		netTimeoutErrp knettyErr
		connClosedErrp knettyErr
	)
	netTimeoutErrp = &netTimeoutErr{}
	connClosedErrp = &connCloseErr{}

	assert.Equal(t, "net io timeout", netTimeoutErrp.Error())
	assert.Equal(t, "net conn interrupted", connClosedErrp.Error())

	assert.Equal(t, true, netTimeoutErrp.TimeoutError())
	assert.Equal(t, false, connClosedErrp.TimeoutError())
}
