package main

import (
	"log"
	"os"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

func getClusterConfig() *gocql.ClusterConfig {
	cass_ip := os.Getenv("CASSANDRA_IPADDRESS")
	log.Print("Trying to connect to container at ", cass_ip)
	cluster := gocql.NewCluster(cass_ip)
	cluster.Consistency = gocql.Quorum
	return cluster
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("Please create a .env file in the root directory of the project")
	}

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
