package cql

import (
	"io"
	"log"
	"net"

	"github.com/st3v/fakesandra/cql/proto"
	"github.com/st3v/fakesandra/cql/proto/v3"
)

type server struct {
	addr      string
	versioner proto.Versioner
	router    proto.Router
}

func NewServer(addr string) *server {
	versioner := proto.NewVersioner()
	versioner.SetRequestFramer(3, v3.RequestFramer())
	versioner.SetResponseFramer(3, v3.ResponseFramer())

	return &server{
		addr:      addr,
		versioner: versioner,
		router:    proto.DefaultRouter(),
	}
}

func (s *server) ListenAndServe() error {
	addr := s.addr
	if addr == "" {
		addr = ":9042"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.Serve(ln)
}

func (s *server) Serve(l net.Listener) error {
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		go s.ServeConnection(c)
	}
}

func (s *server) ServeConnection(c net.Conn) {
	defer c.Close()
	for {
		framer, err := s.versioner.Version(c)
		if err == io.EOF {
			log.Println("Connection closed by client")
			return
		} else if err != nil {
			log.Printf("Error versioning request: %s", err)
			return
		}

		frame, err := framer.Frame(c)
		if err != nil {
			log.Printf("Error framing request: %s", err)
			return
		}

		// Forget about router, use handler chains instead
		// Each handler checks the Frame's opcode and either
		// handles the frame or passes it on to the next handler
		// down the chain
		//
		// This way the caller of ListenAndServe could pass us another
		// handler chain that we could prepend to the default handlers
		// for example a Proxy and Interceptor.
		//
		// Proxy -> Interceptor -> [ QueryHandler -> PrepareHandler -> ... ]
		//
		// Instead of writing the response directly to the ResponseWriter,
		// the HandlerFunc could just pass the response frame upstream
		// through the handler chain on the way out. The server would then
		// take care of writing the frame. That way handlers like the proxy
		// and the interceptor could hijack the response on its way up the
		// chain.

		handler, err := s.router.Route(frame)
		if err != nil {
			log.Printf("Error routing request: %s", err)
			return
		}

		if err := handler.Handle(frame, proto.FrameWriter(c)); err != nil {
			log.Printf("Error handling request: %s", err)
			return
		}
	}
}
