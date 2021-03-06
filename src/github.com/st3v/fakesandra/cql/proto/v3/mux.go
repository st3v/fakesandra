package v3

import "github.com/st3v/fakesandra/cql/proto"

var DefaultMux = NewOpcodeMux()

func NewOpcodeMux() *opcodeMux {
	return &opcodeMux{
		handlers: map[proto.Opcode]proto.FrameHandler{
			proto.OpQuery:   QueryFrameHandler,
			proto.OpStartup: StartupFrameHandler,
		},
	}
}

type opcodeMux struct {
	handlers map[proto.Opcode]proto.FrameHandler
}

func (opmux *opcodeMux) ServeCQL(req proto.Frame, rw proto.ResponseWriter) {
	handler, found := opmux.Handler(req.Opcode())
	if !found {
		// TODO: write error to ResponseWriter
	}

	handler.ServeCQL(req, rw)
}

func (opmux *opcodeMux) Handle(oc proto.Opcode, handler proto.FrameHandler) {
	opmux.handlers[oc] = handler
}

func (opmux *opcodeMux) Handler(oc proto.Opcode) (proto.FrameHandler, bool) {
	handler, found := opmux.handlers[oc]
	return handler, found
}
