package proto

import "io"

// Version represents the version of a CQL frame.
type Version uint8

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
	io.WriterTo
	Version() Version
	Request() bool
	Response() bool
	Opcode() Opcode
	Body() []byte
	QueryHandler() HandlerFunc
}

type ResponseWriter interface {
	Write(Frame) error
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
	Handle(f Frame, rw ResponseWriter) error
}

type HandlerFunc func(f Frame, rw ResponseWriter) error
