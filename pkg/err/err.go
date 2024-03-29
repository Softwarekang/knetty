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
}

var (
	// ConnClosedErr conn closed err
	ConnClosedErr = &connClosedErr{}
	// ClientClosedErr client closed err
	ClientClosedErr = &clientClosedErr{}
	// ServerClosedErr server closed err
	ServerClosedErr = &serverClosedErr{}
	// BufferFullErr buffer is full err
	BufferFullErr = &bufferFullErr{}
	// BufferEmptyErr is empty err
	BufferEmptyErr = &bufferEmptyErr{}
)

type connClosedErr struct{}

// Error implements error.
func (c *connClosedErr) Error() string {
	return "net connection is closed"
}

type clientClosedErr struct{}

// Error implements error.
func (c *clientClosedErr) Error() string {
	return "client has already been closed"
}

type serverClosedErr struct{}

// Error implements error.
func (s *serverClosedErr) Error() string {
	return "server has already been closed"
}

type bufferFullErr struct {
}

// Error implements error.
func (o *bufferFullErr) Error() string {
	return "buffer is full"
}

type bufferEmptyErr struct {
}

func (o *bufferEmptyErr) Error() string {
	return "buffer is empty"
}

type UnKnowNetworkErr string

func (e UnKnowNetworkErr) Error() string { return "unKnowErr network " + string(e) }

type IllegalListenerErr string

func (e IllegalListenerErr) Error() string { return "illegal listener " + string(e) }
