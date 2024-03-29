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

// Package math common mathematical methods
package math

// Max return the larger of the two numbers, a and b.
func Max(a, b int) int {
	if a < b {
		return b
	}

	return a
}

// Min return the smaller of the two numbers, a and b.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
