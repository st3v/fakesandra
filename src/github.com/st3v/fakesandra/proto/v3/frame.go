package v3

import (
	"encoding/binary"
	"io"
)

const (
	versionRequest  = 0x03
	versionResponse = 0x83
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

func read(in io.Reader) (frame, error) {
	f := frame{}

	header := new(header)
	if err := binary.Read(in, binary.BigEndian, header); err != nil {
		return f, err
	}

	body := make([]byte, header.Length)
	if _, err := io.ReadFull(in, body); err != nil {
		return f, err
	}

	return frame{
		version: versionRequest,
		header:  *header,
		body:    body,
	}, nil
}

func write(out io.Writer, f frame) error {
	if _, err := out.Write([]byte{f.version}); err != nil {
		return err
	}

	if err := binary.Write(out, binary.BigEndian, f.header); err != nil {
		return err
	}

	_, err := out.Write(f.body)
	return err
}
