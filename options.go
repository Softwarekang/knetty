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

package knetty

import (
	"github.com/Softwarekang/knetty/session"
)

/*
NewSessionCallBackFunc It is executed when a new session is established,
so some necessary parameters for drawing need to be set to ensure that the session starts properly.
*/
type NewSessionCallBackFunc func(s session.Session) error

// ServerOption option for server
type ServerOption func(*ServerOptions)

// ServerOptions options for server
type ServerOptions struct {
	network    string
	address    string
	newSession NewSessionCallBackFunc
}

// withServerNetwork set network
func withServerNetwork(network string) ServerOption {
	return func(opt *ServerOptions) {
		opt.network = network
	}
}

// withServerAddress set address
func withServerAddress(address string) ServerOption {
	return func(opt *ServerOptions) {
		opt.address = address
	}
}

// WithServiceNewSessionCallBackFunc set newSessionCallBackFunc
func WithServiceNewSessionCallBackFunc(f NewSessionCallBackFunc) ServerOption {
	return func(opt *ServerOptions) {
		opt.newSession = f
	}
}

func newDefaultServerOptions() []ServerOption {
	return []ServerOption{
		withServerAddress("127.0.0.1:8000"),
		withServerNetwork("tcp"),
	}
}

func mergeCustomServerOptions(customServerOptions ...ServerOption) []ServerOption {
	return append(newDefaultServerOptions(), customServerOptions...)
}

// ClientOption option for client
type ClientOption func(options *ClientOptions)

// ClientOptions options for client
type ClientOptions struct {
	network    string
	address    string
	newSession NewSessionCallBackFunc
}

// withClientNetwork set network
func withClientNetwork(network string) ClientOption {
	return func(opt *ClientOptions) {
		opt.network = network
	}
}

// withClientAddress set address
func withClientAddress(address string) ClientOption {
	return func(opt *ClientOptions) {
		opt.address = address
	}
}

// WithClientNewSessionCallBackFunc set newSessionCallBackFunc
func WithClientNewSessionCallBackFunc(f NewSessionCallBackFunc) ClientOption {
	return func(opt *ClientOptions) {
		opt.newSession = f
	}
}

func newDefaultClientOptions() []ClientOption {
	return []ClientOption{
		withClientAddress("127.0.0.1:8000"),
		withClientNetwork("tcp"),
	}
}

func mergeCustomClientOptions(customClientOptions ...ClientOption) []ClientOption {
	return append(newDefaultClientOptions(), customClientOptions...)
}
