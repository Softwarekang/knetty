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

// Package knetty .
package knetty

import (
	"github.com/Softwarekang/knetty/net/poll"
	"github.com/Softwarekang/knetty/pkg/log"
)

// SetPollerNums set reactor goroutine nums
func SetPollerNums(n int) error {
	return poll.PollerManager.SetPollerNums(n)
}

// SetLogger set custom log
func SetLogger(logger log.Logger) {
	log.DefaultLogger = logger
}
