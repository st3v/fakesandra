package proto

import (
	"io"
	"log"
	"net"

	"github.com/st3v/fakesandra/proto/v3"
)

const (
	version1 uint8 = 1 + iota
	version2
	version3
)

var handlers = map[uint8]func(rw io.ReadWriter) error{
	version3: v3.Handle,
}

func Dispatch(c net.Conn) {
	defer c.Close()

	for {
		v, err := version(c)
		if err != nil {
			log.Printf("Error reading protocol version: %s", err)
			return
		}

		h, found := handlers[v]
		if !found {
			log.Printf("Unsupported protocol version '%d'", v)
			return
		}

		if err := h(c); err != nil {
			log.Println(err)
			return
		}
	}
}

func version(r io.Reader) (uint8, error) {
	version := make([]byte, 1)
	if _, err := io.ReadFull(r, version); err != nil {
		return 0, err
	}
	return version[0], nil
}
