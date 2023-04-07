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
