/*
	Copyright 2022 Phoenix

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

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
