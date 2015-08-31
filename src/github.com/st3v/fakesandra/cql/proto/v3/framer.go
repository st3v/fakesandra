package v3

import (
	"io"

	"github.com/st3v/fakesandra/cql/proto"
)

// Specification for CQL protocol v3 can be found under:
// https://github.com/apache/cassandra/blob/trunk/doc/native_protocol_v3.spec

const (
	Version   proto.Version    = 3
	request   proto.VersionDir = 0x00
	response  proto.VersionDir = 0x80
	headerLen                  = 9
)

type header struct {
	Flags    uint8
	StreamID uint16
	Opcode   proto.Opcode
	Length   uint32
}

type frame struct {
	versionDir proto.VersionDir
	header     header
	body       []byte
}

func (f *frame) Version() proto.Version {
	return proto.Version(Version)
}

func (f *frame) Response() bool {
	return f.versionDir&response == response
}

func (f *frame) Request() bool {
	return !f.Response()
}

func (f *frame) Opcode() proto.Opcode {
	return f.header.Opcode
}

func (f *frame) Body() []byte {
	return f.body
}

func (f *frame) WriteTo(w io.Writer) (int64, error) {
	return writeFrame(w, f)
}

type framer struct {
	direction proto.VersionDir
}

func RequestFramer() *framer {
	return &framer{
		direction: request,
	}
}

func ResponseFramer() *framer {
	return &framer{
		direction: response,
	}
}

func (f *framer) Frame(r io.Reader) (proto.Frame, error) {
	frame := &frame{
		versionDir: f.direction | proto.VersionDir(Version),
	}

	if err := readFrame(r, frame); err != nil {
		return nil, err
	}

	return frame, nil
}

func readFrame(in io.Reader, f *frame) error {
	if err := proto.ReadBinary(in, &f.header); err != nil {
		return err
	}

	f.body = make([]byte, f.header.Length)
	if _, err := io.ReadFull(in, f.body); err != nil {
		return err
	}

	return nil
}

func writeFrame(out io.Writer, f *frame) (int64, error) {
	var n int64 = 0
	if err := proto.WriteBinary(out, f.versionDir); err != nil {
		return n, err
	}
	n += 1

	if err := proto.WriteBinary(out, f.header); err != nil {
		return n, err
	}
	n += headerLen

	m, err := out.Write(f.body)
	return n + int64(m), err
}
