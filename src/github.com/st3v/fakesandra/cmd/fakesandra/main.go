package main

import (
	"fmt"
	"log"

	"github.com/st3v/fakesandra"
	"github.com/st3v/fakesandra/middleware/frame"
	"github.com/st3v/fakesandra/middleware/query"
)

func main() {
	fmt.Println("Work in Progress!")

	// use middleware to log frames
	frameHandler := frame.Logger(log.Print, fakesandra.DefaultHandler)

	// use middleware to log queries
	fakesandra.HandleQuery(query.Logger(log.Print))

	if err := fakesandra.ListenAndServe(":9042", frameHandler); err != nil {
		panic(err)
	}
}
