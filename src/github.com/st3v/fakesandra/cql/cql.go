package cql

import (
	"io"
	"net"
)

// Version represents the version of a CQL frame.
type Version uint8

// VersionDir represents the version AND direction of a CQL frame.
type VersionDir uint8

// Represents a CQL frame.
type Frame interface {
	Version() Version
	Request() bool
	Response() bool
	Header() []byte
	Body() []byte
	Bytes() []byte
}

// Server serves a CQL connection.
type Server interface {
	Serve(con net.Conn) error
}

// Versioner identifies the right Framer that should be used to frame
// incomming bytes.
type Versioner interface {
	Version(in io.Reader) (Framer, error)
}

// Framer reads raw bytes off a reader and frames them according to a
// particular version of the CQL protocol.
type Framer interface {
	Frame(in io.Reader) (Frame, error)
}

// Router routes frames to handler chains.
type Router interface {
	Route(request Frame) (Handler, error)
}

// Handler handles a CQL frame. Might write replies to the response writer.
type Handler interface {
	Handle(request Frame, responseWriter io.Writer) error
}

type HandlerFunc func(request Frame, responseWriter io.Writer) error
