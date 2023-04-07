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

// Package err wrapped err for knetty
package err

// knettyErr wrapped err for net
type knettyErr interface {
	error

	// TimeoutError if an error is caused by a timeout will return true
	TimeoutError() bool
}

var (
	// NetIOTimeoutErr net io err
	NetIOTimeoutErr = &netTimeoutErr{}
	// ConnClosedErr conn closed err
	ConnClosedErr = &connClosedErr{}
	// ClientClosedErr client closed err
	ClientClosedErr = &clientClosedErr{}
	// ServerClosedErr server closed err
	ServerClosedErr = &serverClosedErr{}
)

type netTimeoutErr struct{}

// Error implements error.
func (n netTimeoutErr) Error() string {
	return "net io timeout"
}

// TimeoutError implements knettyErr.
func (n netTimeoutErr) TimeoutError() bool {
	return true
}

type connClosedErr struct{}

// Error implements error.
func (c *connClosedErr) Error() string {
	return "net conn is closed"
}

// TimeoutError implements knettyErr.
func (c *connClosedErr) TimeoutError() bool {
	return false
}

type clientClosedErr struct{}

// Error implements error.
func (c *clientClosedErr) Error() string {
	return "client has already been closed"
}

// TimeoutError implements knettyErr.
func (c *clientClosedErr) TimeoutError() bool {
	return false
}

type serverClosedErr struct{}

// Error implements error.
func (s *serverClosedErr) Error() string {
	return "server has already been closed"
}

// TimeoutError implements knettyErr.
func (s *serverClosedErr) TimeoutError() bool {
	return false
}
