package v3

import (
	"bytes"

	"github.com/st3v/fakesandra/cql/proto"
)

type ResultCode int32

const (
	ResultVoid ResultCode = 1 + iota
	ResultRows
	ResultSetKeyspace
	ResultPrepared
	ResultSchemaChange
)

func ReadyResponse(request proto.Frame) proto.Frame {
	hdr := header{
		Opcode:   proto.OpReady,
		StreamID: request.StreamID(),
		Length:   0,
	}

	return &frame{
		versionDir: proto.VersionDir(Version) | response,
		header:     hdr,
		body:       make([]byte, 0),
	}
}

func ResultVoidResponse(request proto.Frame) proto.Frame {
	buf := new(bytes.Buffer)
	proto.WriteBinary(buf, ResultVoid)

	hdr := header{
		Opcode:   proto.OpResult,
		StreamID: request.StreamID(),
		Length:   uint32(buf.Len()),
	}

	return &frame{
		versionDir: proto.VersionDir(Version) | response,
		header:     hdr,
		body:       buf.Bytes(),
	}
}
