package v3

import (
	"bytes"

	"github.com/st3v/fakesandra/cql/proto"
)

var StartupHandler = proto.FrameHandlerFunc(startupHandler)

func startupHandler(req proto.Frame, rw proto.ResponseWriter) {
	// log.Println("Received STARTUP request")
	rw.WriteFrame(ReadyResponse(req))
}

var QueryHandler = proto.FrameHandlerFunc(queryHandler)

func queryHandler(req proto.Frame, rw proto.ResponseWriter) {
	var qry Query
	if err := readQuery(bytes.NewReader(req.Body()), &qry); err != nil {
		// TODO: write error to ResponseWriter
		return
	}

	// log.Printf("Received QUERY request: %s", qry)
	rw.WriteFrame(ResultVoidResponse(req))
}
