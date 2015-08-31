package fakesandra

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/st3v/fakesandra/cql/proto"
	"github.com/st3v/fakesandra/cql/proto/v3"
)

const DefaultPort = 9042

func ListenAndServe(addr string, handler proto.FrameHandler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServe()
}

type server struct {
	addr      string
	versioner proto.Versioner
	handler   proto.FrameHandler
}

var DefaultVersioner = func() proto.Versioner {
	versioner := proto.NewVersioner()
	versioner.SetRequestFramer(proto.Version3, v3.RequestFramer())
	// TODO: Do we really need response framers? We are a server after all.
	versioner.SetResponseFramer(proto.Version3, v3.ResponseFramer())
	return versioner
}()

var DefaultHandler = func() proto.FrameHandler {
	vmux := proto.NewVersionMux()
	vmux.Handle(3, v3.DefaultMux)
	return vmux
}()

func NewServer(addr string, handler proto.FrameHandler) *server {
	if handler == nil {
		handler = DefaultHandler
	}

	return &server{
		addr:      addr,
		versioner: DefaultVersioner,
		handler:   handler,
	}
}

func (s *server) ListenAndServe() error {
	addr := s.addr
	if addr == "" {
		addr = fmt.Sprintf(":%d", DefaultPort)
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

	log.Println("Serving new connection ...")

	for {
		// TODO: Handle timeouts
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

		s.handler.ServeCQL(frame, proto.FrameWriter(c))
	}
}
