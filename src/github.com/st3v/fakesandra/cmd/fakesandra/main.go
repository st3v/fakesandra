package main

import (
	"fmt"
	"log"
	"net"

	"github.com/st3v/fakesandra/proto"
)

func main() {
	fmt.Println("Work in Progress!")

	l, err := net.Listen("tcp", ":9042")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go proto.Dispatch(conn)
	}
}
