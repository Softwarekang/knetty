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

package buffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteBuffer_Write(t *testing.T) {
	var (
		err error
	)
	byteBuffer := NewByteBuffer()

	data := []byte{'1', '2', '3'}
	err = byteBuffer.Write(data)
	assert.Nil(t, err)
	assert.Equal(t, 3, byteBuffer.Len())

	err = byteBuffer.WriteString("test data")
	assert.Nil(t, err)
	assert.Equal(t, 12, byteBuffer.Len())

	assert.Equal(t, "123test data", string(byteBuffer.Bytes()))

	outData1 := make([]byte, 10)
	n, err := byteBuffer.Read(outData1)
	assert.Nil(t, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, "123test da", string(outData1[:n]))

	outData2 := make([]byte, 10)
	n, err = byteBuffer.Read(outData2)
	assert.Nil(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, "ta", string(outData2[:n]))

	err = byteBuffer.WriteString("12345")
	assert.Nil(t, err)
	assert.Equal(t, 5, byteBuffer.Len())

	byteBuffer.Release(2)
	assert.Equal(t, "345", string(byteBuffer.Bytes()))

	assert.Equal(t, false, byteBuffer.IsEmpty())

	byteBuffer.Clear()

	assert.Equal(t, true, byteBuffer.IsEmpty())
}
