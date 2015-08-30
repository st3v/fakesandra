package main

import (
	"fmt"

	"github.com/st3v/fakesandra/cql"
)

func main() {
	fmt.Println("Work in Progress!")
	if err := cql.ListenAndServe(":9042"); err != nil {
		panic(err)
	}
}
