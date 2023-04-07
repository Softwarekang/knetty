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

package session

// Codec for session
type Codec interface {
	// Encode will convert object to binary network data
	Encode(pkg interface{}) ([]byte, error)

	// Decode will convert binary network data into upper-layer protocol objects.
	// The following three conditions are used to distinguish abnormal, half - wrapped, normal and sticky packets.
	// Exceptions: nil,0,err
	// Half-pack: nil,0,nil
	// Normal & Sticky package: pkg,pkgLen,nil
	Decode([]byte) (interface{}, int, error)
}

// EventListener listener for session event
type EventListener interface {
	// OnConnect runs when the connection initialized
	OnConnect(s Session)
	// OnMessage runs when the session gets a pkg
	OnMessage(s Session, pkg interface{})
	// OnError runs when the session err
	OnError(s Session, e error)
	// OnClose runs before the session closed
	OnClose(s Session)
}
