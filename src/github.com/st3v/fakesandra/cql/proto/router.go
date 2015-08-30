package proto

type router struct {
	handlers map[Opcode]Handler
}

func DefaultRouter() *router {
	return &router{
		handlers: map[Opcode]Handler{
			OpStartup: StartupHandler(),
			OpQuery:   QueryHandler(),
		},
	}
}

func (r *router) Route(f Frame) (Handler, error) {
	h, found := r.handlers[f.Opcode()]
	if !found {
		return nil, errMissingRoute
	}

	return h, nil
}
