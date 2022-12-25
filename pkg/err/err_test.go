package err

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKnettyError(t *testing.T) {
	var (
		netTimeoutErrp   knettyErr
		connClosedErrp   knettyErr
		clientClosedErrp knettyErr
		serverClosedErrp knettyErr
	)
	netTimeoutErrp = &netTimeoutErr{}
	connClosedErrp = &connClosedErr{}
	clientClosedErrp = &clientClosedErr{}
	serverClosedErrp = &serverClosedErr{}
	assert.Equal(t, "net io timeout", netTimeoutErrp.Error())
	assert.Equal(t, "net conn is closed", connClosedErrp.Error())
	assert.Equal(t, "client has already been closed", clientClosedErrp.Error())
	assert.Equal(t, "server has already been closed", serverClosedErrp.Error())

	assert.Equal(t, true, netTimeoutErrp.TimeoutError())
	assert.Equal(t, false, connClosedErrp.TimeoutError())
	assert.Equal(t, false, clientClosedErrp.TimeoutError())
	assert.Equal(t, false, serverClosedErrp.TimeoutError())
}
