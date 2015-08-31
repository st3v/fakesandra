package v3

import (
	"bytes"
	"log"

	"github.com/st3v/fakesandra/cql/proto"
)

var StartupHandler = proto.StartupHandler

var QueryHandler = proto.FrameHandlerFunc(queryHandler)

func queryHandler(req proto.Frame, rw proto.ResponseWriter) {
	var qry Query
	if err := readQuery(bytes.NewReader(req.Body()), &qry); err != nil {
		// TODO: write error to ResponseWriter
		return
	}

	log.Printf("Received QUERY request: %s", qry)

	// return voidResponse(f.header.StreamID), nil
}
