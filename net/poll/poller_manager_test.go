/*
	Copyright 2022 ankangan

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
