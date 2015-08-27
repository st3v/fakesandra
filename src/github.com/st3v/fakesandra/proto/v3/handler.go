package v3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

const (
	opError opcode = iota
	opStartup
	opReady
	opAuthenticate
	deprecated
	opOptions
	opSupported
	opQuery
	opResult
	opPrepare
	opExecute
	opRegister
	opEvent
	opBatch
	opAuthChallenge
	opAuthResponse
	opAuthSuccess
)

func Handle(rw io.ReadWriter) error {
	var out frame
	in, err := read(rw)
	if err != nil {
		// todo: write error frame
		return fmt.Errorf("Error reading v3 frame: %s", err)
	}

	out, err = handleFrame(in)
	if err != nil {
		// todo: write error frame
		return fmt.Errorf("Error handling v3 frame: %s", err)
	}

	if err := write(rw, out); err != nil {
		return fmt.Errorf("Error writing v3 frame: %s", err)
	}

	return nil
}

type framer func(frame) (frame, error)

var framers = map[opcode]framer{
	opStartup: startup,
	opQuery:   query,
}

func handleFrame(f frame) (frame, error) {
	framer, ok := framers[f.header.Opcode]
	if !ok {
		return frame{}, fmt.Errorf("Unsupported opcode: ", f.header.Opcode)
	}

	return framer(f)
}

func startup(f frame) (frame, error) {
	log.Println("Received STARTUP request")
	return ready(f.header.StreamID), nil
}

func query(f frame) (frame, error) {
	log.Printf("Received QUERY request: %s", strings.Trim(string(f.body), " \r\n"))
	return void(f.header.StreamID), nil
}

func void(streamID uint16) frame {
	log.Println("Sending VOID response")
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(1))
	return response(streamID, opResult, buf.Bytes())
}

func ready(streamID uint16) frame {
	log.Println("Sending READY response")
	return response(streamID, opReady, []byte{})
}

func response(streamID uint16, op opcode, body []byte) frame {
	h := header{
		Flags:    0,
		StreamID: streamID,
		Opcode:   op,
		Length:   uint32(len(body)),
	}

	return frame{
		version: versionResponse,
		header:  h,
		body:    body,
	}

}
