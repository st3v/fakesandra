package v3

import (
	"bytes"
	"log"

	"github.com/st3v/fakesandra/cql/proto"
)

func queryHandler(f proto.Frame, w proto.ResponseWriter) error {
	var qry Query
	if err := readQuery(bytes.NewReader(f.Body()), &qry); err != nil {
		return err
	}

	log.Printf("Received QUERY request: %s", qry)

	// return voidResponse(f.header.StreamID), nil
	return nil
}
