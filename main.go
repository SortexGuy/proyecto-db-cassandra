package main

import (
	"log"

	"github.com/gocql/gocql"
)

func getClusterConfig() *gocql.ClusterConfig {
	cluster := gocql.NewCluster("172.18.0.2")
	cluster.Consistency = gocql.Quorum
	return cluster
}

func main() {
	cluster := getClusterConfig()

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Unable to open up a session with the Cassandra database!", err)
	}
	result := session.Query("CREATE KEYSPACE IF NOT EXISTS app WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };")
	err = result.Exec()
	if err != nil {
		log.Println(err)
	}

	log.Println("All good")

	defer session.Close()
}
