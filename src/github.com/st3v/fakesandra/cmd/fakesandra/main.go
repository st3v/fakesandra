package main

import (
	"fmt"

	"github.com/st3v/fakesandra"
)

func main() {
	fmt.Println("Work in Progress!")
	if err := fakesandra.ListenAndServe(":9042"); err != nil {
		panic(err)
	}
}
