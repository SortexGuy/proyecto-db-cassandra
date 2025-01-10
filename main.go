package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
	"log"
	//"net/http"
	"os"
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

	// Inicializa los repositorios
	movieRepo := movies.NewMovieRepository(session)
	movieByUserRepo := movies.NewMovieByUserRepository(session)

	// Inicializa el controlador
	movieController := movies.NewMovieController(movieRepo, movieByUserRepo)

	// Probar GetMoviesByUser
	userID := int64(6) // Cambia esto al ID de usuario que deseas probar
	moviesByUser, err := movieController.GetMoviesByUser(userID)
	if err != nil {
		log.Println("Error:", err)
	} else {
		log.Printf("Movies by user %d: %+v\n", userID, moviesByUser)
	}
}

// TODO: Execute code
//r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

//r.Run()
//}
