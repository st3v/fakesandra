package proto

import (
	"encoding/binary"
	"io"
)

func WriteBinary(w io.Writer, data interface{}) error {
	return binary.Write(w, binary.BigEndian, data)
}

func WriteByte(w io.Writer, n uint8) error {
	_, err := w.Write([]byte{n})
	return err
}

func WriteShort(w io.Writer, n uint16) error {
	return WriteBinary(w, n)
}

func WriteInt(w io.Writer, n int32) error {
	return WriteBinary(w, n)
}

func WriteLong(w io.Writer, n int64) error {
	return WriteBinary(w, n)
}

func WriteShortBytes(w io.Writer, b []byte) error {
	if len(b) > 1<<16-1 {
		return errMaxLenExceeded
	}

	if err := WriteShort(w, uint16(len(b))); err != nil {
		return err
	}

	_, err := w.Write(b)
	return err
}

func WriteBytes(w io.Writer, b []byte) error {
	if len(b) > 1<<32-1 {
		return errMaxLenExceeded
	}

	if err := WriteInt(w, int32(len(b))); err != nil {
		return err
	}

	_, err := w.Write(b)
	return err
}

func WriteString(w io.Writer, str string) error {
	return WriteShortBytes(w, []byte(str))
}

func WriteLongString(w io.Writer, str string) error {
	return WriteBytes(w, []byte(str))
}

func ReadBinary(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.BigEndian, data)
}

func ReadByte(r io.Reader, n *uint8) error {
	return ReadBinary(r, n)
}

func ReadShort(r io.Reader, n *uint16) error {
	return ReadBinary(r, n)
}

func ReadInt(r io.Reader, n *int32) error {
	return ReadBinary(r, n)
}

func ReadLong(r io.Reader, n *int64) error {
	return ReadBinary(r, n)
}

func ReadBytes(r io.Reader) ([]byte, error) {
	var n int32
	if err := ReadInt(r, &n); err != nil {
		return []byte{}, err
	}

	b := make([]byte, n)
	if _, err := io.ReadFull(r, b); err != nil {
		return []byte{}, err
	}

	return b, nil
}

func ReadShortBytes(r io.Reader) ([]byte, error) {
	var n uint16
	if err := ReadShort(r, &n); err != nil {
		return []byte{}, err
	}

	b := make([]byte, n)
	if _, err := io.ReadFull(r, b); err != nil {
		return []byte{}, err
	}

	return b, nil
}

func ReadString(r io.Reader) (string, error) {
	str, err := ReadShortBytes(r)
	return string(str), err
}

func ReadLongString(r io.Reader) (string, error) {
	str, err := ReadBytes(r)
	return string(str), err
}

func ReadConsistency(r io.Reader, c *Consistency) error {
	return ReadBinary(r, c)
}

func readVersionDir(r io.Reader, v *VersionDir) error {
	return binary.Read(r, binary.BigEndian, v)
}

type frameWriter struct {
	out io.Writer
}

func FrameWriter(w io.Writer) *frameWriter {
	return &frameWriter{w}
}

func (fw *frameWriter) WriteFrame(f Frame) error {
	_, err := f.WriteTo(fw.out)
	return err
}
