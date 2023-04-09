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

package ring_buffer

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name    string
		args    args
		wantCap int
	}{
		{name: "test case 1", args: args{size: 1}, wantCap: 1 << 6},
		{name: "test case 2", args: args{size: 1 << minAllocBit}, wantCap: 1 << minAllocBit},
		{name: "test case 3", args: args{size: 1<<minAllocBit + 1}, wantCap: 1 << 7},
		{name: "test case 4", args: args{size: 1<<maxAllocBit + 1}, wantCap: 1 << maxAllocBit},
		{name: "test case 5", args: args{size: 1 << maxAllocBit}, wantCap: 1 << maxAllocBit},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.size); !reflect.DeepEqual(cap(got), tt.wantCap) {
				t.Errorf("Get() = %v, want %v", cap(got), tt.wantCap)
			}
		})
	}
}

func TestPut(t *testing.T) {
	bytes1 := Get(1)
	bytes1[1] = '1'
	Put(bytes1)
	// get released buf
	oldBytes := Get(1)
	// check index 1 value
	assert.Equal(t, uint8('1'), oldBytes[1])
}
