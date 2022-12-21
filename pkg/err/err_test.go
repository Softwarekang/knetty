package err

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKnettyError(t *testing.T) {
	var (
		netTimeoutErrp knettyErr
		connClosedErrp knettyErr
	)
	netTimeoutErrp = &netTimeoutErr{}
	connClosedErrp = &connCloseErr{}

	assert.Equal(t, "net io timeout", netTimeoutErrp.Error())
	assert.Equal(t, "net conn is closed", connClosedErrp.Error())

	assert.Equal(t, true, netTimeoutErrp.TimeoutError())
	assert.Equal(t, false, connClosedErrp.TimeoutError())
}
