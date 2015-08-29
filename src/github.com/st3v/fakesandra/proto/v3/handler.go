package v3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
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
	opStartup: startupHandler,
	opQuery:   queryHandler,
}

func handleFrame(f frame) (frame, error) {
	framer, ok := framers[f.header.Opcode]
	if !ok {
		return frame{}, fmt.Errorf("Unsupported opcode: ", f.header.Opcode)
	}

	return framer(f)
}

func startupHandler(f frame) (frame, error) {
	log.Println("Received STARTUP request")
	return readyResponse(f.header.StreamID), nil
}

func queryHandler(f frame) (frame, error) {
	var qry Query
	if err := readQuery(bytes.NewReader(f.body), &qry); err != nil {
		return frame{}, err
	}

	log.Printf("Received QUERY request: %s", qry)

	return voidResponse(f.header.StreamID), nil
}

func voidResponse(streamID uint16) frame {
	log.Println("Sending VOID response")
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(1))
	return response(streamID, opResult, buf.Bytes())
}

func readyResponse(streamID uint16) frame {
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
