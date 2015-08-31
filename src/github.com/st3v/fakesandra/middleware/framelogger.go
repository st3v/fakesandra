package middleware

import (
	"fmt"

	"github.com/st3v/fakesandra/cql/proto"
)

func FrameLogger(log func(...interface{}), next proto.FrameHandler) proto.FrameHandler {
	return proto.FrameHandlerFunc(func(req proto.Frame, rw proto.ResponseWriter) {
		log(fmt.Sprintf("Read: %s", req))
		next.ServeCQL(req, &responseLogger{rw, log})
	})
}

type responseLogger struct {
	out proto.ResponseWriter
	log func(...interface{})
}

func (rl *responseLogger) WriteFrame(f proto.Frame) error {
	rl.log(fmt.Sprintf("Sent: %s", f))
	return rl.out.WriteFrame(f)
}
