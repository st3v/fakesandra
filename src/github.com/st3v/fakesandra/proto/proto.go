package proto

import (
	"encoding/binary"
	"io"
)

func Read(in io.Reader, f *Frame) error {
	if err := binary.Read(in, binary.BigEndian, &f.header); err != nil {
		return err
	}

	f.body = make([]byte, f.header.Length)

	_, err := io.ReadFull(in, f.body)
	return err
}

type FrameHeader struct {
	Version uint8
	Flags   uint8
	Stream  uint16
	OpCode  uint8
	Length  uint32
}

type FrameBody []byte

type Frame struct {
	header FrameHeader
	body   FrameBody
}

func NewFrame(h FrameHeader, b FrameBody) Frame {
	return Frame{h, b}
}

func Write(out io.Writer, f Frame) error {
	if err := binary.Write(out, binary.BigEndian, f.header); err != nil {
		return err
	}

	_, err := out.Write(f.body)
	return err
}
