package main

import (
	"fmt"
	"log"

	"github.com/st3v/fakesandra"
	"github.com/st3v/fakesandra/middleware"
)

func main() {
	fmt.Println("Work in Progress!")

	// use FrameLogger middleware
	handler := middleware.FrameLogger(log.Print, fakesandra.DefaultHandler)

	if err := fakesandra.ListenAndServe(":9042", handler); err != nil {
		panic(err)
	}
}
