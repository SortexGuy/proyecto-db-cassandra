package main

import (
	"log"
	"os"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

var SESSION *gocql.Session

func getClusterConfig() *gocql.ClusterConfig {
	cass_ip := os.Getenv("CASSANDRA_IPADDRESS")
	log.Println("Trying to connect to container at ", cass_ip)
	cluster := gocql.NewCluster(cass_ip)
	cluster.Keyspace = "app"
	cluster.Consistency = gocql.Quorum
	return cluster
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Please create a .env file in the root directory of the project")
	}

	cluster := getClusterConfig()

	log.Println("Trying to open Cassandra session")
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Unable to open up a session with the Cassandra database!", err)
	}
	SESSION = session

	// TODO: Execute code

	log.Println("All good")

	defer SESSION.Close()
}
