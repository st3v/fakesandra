package cql

import (
	"encoding/binary"
	"io"
)

func readVersionDir(r io.Reader, v *VersionDir) error {
	return binary.Read(r, binary.BigEndian, v)
}
