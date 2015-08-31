package proto

import (
	"fmt"
	"io"
)

// Version represents the version of a CQL frame.
type Version uint8

const (
	Version1 Version = 1 + iota
	Version2
	Version3
)

var Versions = []Version{
	Version1,
	Version2,
	Version3,
}

type Consistency uint16

const (
	Any Consistency = iota
	One
	Two
	Three
	Quorum
	All
	LocalQuorum
	EachQuorum
	Serial
	LocalSerial
	LocalOne
)

func (c Consistency) String() string {
	switch c {
	case Any:
		return "ANY"
	case One:
		return "ONE"
	case Two:
		return "TWO"
	case Three:
		return "THREE"
	case Quorum:
		return "QUORUM"
	case All:
		return "ALL"
	case LocalQuorum:
		return "LOCAL_QUORUM"
	case EachQuorum:
		return "EACH_QUORUM"
	case Serial:
		return "SERIAL"
	case LocalSerial:
		return "LOCAL_SERIAL"
	case LocalOne:
		return "LOCAL_SERIAL"
	default:
		return "UNKNOWN"
	}
}

// Represents a CQL frame.
type Frame interface {
	fmt.Stringer
	io.WriterTo
	Version() Version
	Request() bool
	Response() bool
	Opcode() Opcode
	StreamID() uint16

	// TODO: make Body() return an io.Reader, i.e. the framer should not read
	// the body but wrap the underlying readr into a io.LimitReader and attach it
	// to the frame
	Body() []byte
}

// TODO: Pull the pipline related stuff below out of the proto package and
// put it into top-level packages for versioning, framing, frame handlers,
// and query handlers. That way it will be easier to write version agnostic
// handlers and middleware.
//
// In the end proto should only contain low-level code related to reading,
// parsing and writing frames, queries, etc.

type ResponseWriter interface {
	WriteFrame(response Frame) error
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

// FrameHandler handles a CQL request frame. Might write replies to the
// response writer.
// TODO: Rename ServeCQL to ServeFrame, for clarity reasons
type FrameHandler interface {
	ServeCQL(request Frame, rw ResponseWriter)
}

type FrameHandlerFunc func(request Frame, rw ResponseWriter)

func (fn FrameHandlerFunc) ServeCQL(request Frame, rw ResponseWriter) {
	fn(request, rw)
}

type QueryFrameHandler interface {
	FrameHandler
	Prepend(handler QueryHandler)
}

// TODO: instead of passing a string, create interface for query
type QueryHandler interface {
	ServeQuery(query string, request Frame, rw ResponseWriter)
}

type QueryHandlerFunc func(query string, request Frame, rw ResponseWriter)

func (fn QueryHandlerFunc) ServeQuery(q string, r Frame, rw ResponseWriter) {
	fn(q, r, rw)
}
