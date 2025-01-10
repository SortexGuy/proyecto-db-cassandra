package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SortexGuy/proyecto-db-cassandra/src/movies"
	"github.com/SortexGuy/proyecto-db-cassandra/src/users"
	"github.com/gin-gonic/gin"
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
	defer SESSION.Close()

	// TODO: Execute code
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	movies.RegisterRoutes(r)
	users.RegisterRoutes(r)

	r.Run()
}
