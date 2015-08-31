package query

import (
	"fmt"

	"github.com/st3v/fakesandra/cql/proto"
)

func Logger(log func(...interface{})) proto.QueryHandler {
	return proto.QueryHandlerFunc(
		func(qry string, req proto.Frame, rw proto.ResponseWriter) {
			log(fmt.Sprintf("Received query: %s", qry))
		},
	)
}
