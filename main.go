package main

import (
	"log"
	"net/http"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
	"github.com/SortexGuy/proyecto-db-cassandra/src/auth"
	"github.com/SortexGuy/proyecto-db-cassandra/src/movies"
	"github.com/SortexGuy/proyecto-db-cassandra/src/recommendations"
	"github.com/SortexGuy/proyecto-db-cassandra/src/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Please create a .env file in the root directory of the project")
	}

	cluster := config.GetClusterConfig(false)
	log.Println("Trying to open Cassandra session")
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Unable to open up a session with the Cassandra database!", err)
	}
	config.SESSION = session
	defer config.SESSION.Close()

	// TODO: Execute code
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	auth.RegisterRoutes(r)
	movies.RegisterRoutes(r)
	users.RegisterRoutes(r)
	recommendations.RegisterRoutes(r)

	r.Run()
}
