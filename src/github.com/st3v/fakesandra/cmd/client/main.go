package main

import (
	"log"

	"github.com/gocql/gocql"
)

func main() {
	cluster := gocql.NewCluster("127.0.0.1")

	for v := 3; v < 4; v++ {
		cluster.ProtoVersion = v
		session, err := cluster.CreateSession()
		if err != nil {
			log.Printf("Error openning session: %s\n", err.Error())
			continue
		}

		err = session.Query(`
			CREATE KEYSPACE foo 
			WITH REPLICATION {
				'class': 'SimpleStrategy', 
				'replication_factor': 3,
			}`,
		).Exec()

		if err != nil {
			log.Printf("Error executing query: %s\n", err.Error())
		}
	}
}
