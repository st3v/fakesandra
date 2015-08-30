package v3

import (
	"encoding/binary"
	"io"
)

// Specification for CQL protocol v3 can be found under:
// https://github.com/apache/cassandra/blob/trunk/doc/native_protocol_v3.spec

const (
	version  uint8 = 3
	request  uint8 = 0x00
	response uint8 = 0x80
)

type opcode uint8

type header struct {
	Flags    uint8
	StreamID uint16
	Opcode   opcode
	Length   uint32
}

type frame struct {
	version uint8
	header  header
	body    []byte
}

type framer struct {
	direction uint8
	version   uint8
}

func RequestFramer() *framer {
	return &framer{
		version:   version,
		direction: request,
	}
}

func ResponseFramer() *framer {
	return &framer{
		version:   version,
		direction: response,
	}
}

func (f *framer) Frame(r io.Reader) (*frame, error) {
	frame := &frame{
		version: f.direction | f.version,
	}

	if err := readFrame(r, frame); err != nil {
		return nil, err
	}

	return frame, nil
}

func readFrame(in io.Reader, f *frame) error {
	if err := binary.Read(in, binary.BigEndian, &f.header); err != nil {
		return err
	}

	f.body = make([]byte, f.header.Length)
	if _, err := io.ReadFull(in, f.body); err != nil {
		return err
	}

	return nil
}
