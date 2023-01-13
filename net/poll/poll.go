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

// Package poll impl io multiplexing on different systems
package poll

// Poll net poll interface
type Poll interface {
	// Register netFd in the poller. events is the type of event that the poller focus on
	Register(netFd *NetFileDesc, eventType EventType) error

	// Wait
	// poller will focus on all registered netFd, wait for netFd to satisfy the condition and
	// notify the registered listener, so it is blocked
	Wait() error

	// Close the poller
	Close() error
}

// NetFileDesc file-desc for net-fd
type NetFileDesc struct {
	// FD system fd
	FD int
	// listener for poller
	NetPollListener
}

// NetPollListener listener for net poller
type NetPollListener struct {
	// OnRead will run where fd is readable
	OnRead
	// OnWrite will run where fd is writeable
	OnWrite
	// OnInterrupt will run where fd is interrupted
	OnInterrupt
}

// OnRead the callback function when the net fd state is readable
type OnRead func() error

// OnWrite The callback function when the net fd state is writable
type OnWrite func() error

// OnInterrupt The callback function when the net fd state is interrupt
type OnInterrupt func() error

// EventType event type for poller
type EventType int

const (
	Read EventType = iota + 1
	DeleteRead
	ReadToRW
	RwToRead
	OnceWrite
)
