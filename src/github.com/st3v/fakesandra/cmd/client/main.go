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

		session.SetPageSize(123)
		session.SetConsistency(gocql.LocalQuorum)

		qry := session.Query(`
			CREATE KEYSPACE foo 
			WITH REPLICATION {
				'class': 'SimpleStrategy', 
				'replication_factor': 3
			}`,
			"foo",
			"bar",
		)

		// qry.Consistency(gocql.Any).PageSize(987)
		qry.SerialConsistency(gocql.LocalSerial)
		if err := qry.Exec(); err != nil {
			log.Printf("Error executing query: %s\n", err.Error())
		}
	}
}
