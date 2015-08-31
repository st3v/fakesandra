package proto

type OpcodeMux interface {
	FrameHandler
	Handle(oc Opcode, handler FrameHandler)
	Handler(oc Opcode) (FrameHandler, bool)
}

func NewVersionMux() *VersionMux {
	return &VersionMux{
		handlers: make(map[Version]OpcodeMux),
	}
}

type VersionMux struct {
	handlers map[Version]OpcodeMux
}

func (vmux *VersionMux) ServeCQL(req Frame, rw ResponseWriter) {
	handler, found := vmux.Handler(req.Version())
	if !found {
		// TODO: write error to ResponseWriter
	}

	handler.ServeCQL(req, rw)
}

func (vmux *VersionMux) Handle(v Version, handler OpcodeMux) {
	vmux.handlers[v] = handler
}

func (vmux *VersionMux) Handler(v Version) (handler OpcodeMux, found bool) {
	handler, found = vmux.handlers[v]
	return
}
