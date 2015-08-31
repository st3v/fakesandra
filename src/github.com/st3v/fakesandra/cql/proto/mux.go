package proto

func NewVersionMux() *versionMux {
	return &versionMux{
		handlers: make(map[Version]FrameHandler),
	}
}

type versionMux struct {
	handlers map[Version]FrameHandler
}

func (vmux *versionMux) ServeCQL(req Frame, rw ResponseWriter) {
	handler, found := vmux.handlers[req.Version()]
	if !found {
		// TODO: write error to ResponseWriter
	}

	handler.ServeCQL(req, rw)
}

func (vmux *versionMux) Handle(v Version, handler FrameHandler) {
	vmux.handlers[v] = handler
}
