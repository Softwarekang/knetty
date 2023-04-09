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

// Package utils generic tool method implementation.
package utils

// AdjustNToPowerOfTwo adjust n to the first value greater than or equal its 2^n.
func AdjustNToPowerOfTwo(n int) int {
	if IsPowerOfTwo(n) {
		return n
	}
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return n + 1
}

// IsPowerOfTwo check whether n is 2 to the nth power.
func IsPowerOfTwo(n int) bool {
	return n&(n-1) == 0
}
