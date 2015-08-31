package v3

import (
	"bytes"

	"github.com/st3v/fakesandra/cql/proto"
)

var StartupFrameHandler = proto.FrameHandlerFunc(startupFrameHandler)

var QueryFrameHandler = NewQueryFrameHandler(ResultVoidHandler)

var ResultVoidHandler = proto.QueryHandlerFunc(resultVoidHandler)

func resultVoidHandler(qry string, req proto.Frame, rw proto.ResponseWriter) {
	rw.WriteFrame(ResultVoidResponse(req))
}

func startupFrameHandler(req proto.Frame, rw proto.ResponseWriter) {
	// log.Println("Received STARTUP request")
	rw.WriteFrame(ReadyResponse(req))
}

func NewQueryFrameHandler(handler proto.QueryHandler) *queryFrameHandler {
	return &queryFrameHandler{
		queryHandler: handler,
	}
}

type queryFrameHandler struct {
	queryHandler proto.QueryHandler
}

func (qfm *queryFrameHandler) ServeCQL(req proto.Frame, rw proto.ResponseWriter) {
	var qry Query
	if err := readQuery(bytes.NewReader(req.Body()), &qry); err != nil {
		// TODO: write error to ResponseWriter
		return
	}

	qfm.queryHandler.ServeQuery(qry.TrimmedStatement(), req, rw)
}

func (qfm *queryFrameHandler) Prepend(handler proto.QueryHandler) {
	next := qfm.queryHandler
	qfm.queryHandler = proto.QueryHandlerFunc(
		func(qry string, req proto.Frame, rw proto.ResponseWriter) {
			handler.ServeQuery(qry, req, rw)
			next.ServeQuery(qry, req, rw)
		},
	)
}
