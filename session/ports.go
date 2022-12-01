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

// EventListener listen for session
type EventListener interface {
	// OnMessage runs when the session gets a pkg
	OnMessage(s Session, pkg interface{})
	// OnConnect runs when the connection initialized
	OnConnect(s Session)
	// OnClose runs before the session closed
	OnClose(s Session)
	// OnError runs when the session err
	OnError(s Session, e error)
}
