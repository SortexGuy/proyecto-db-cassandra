package auth

import (
	"net/http"

	"github.com/SortexGuy/proyecto-db-cassandra/src/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/auth")

	group.POST("/login", loginUserController)
	group.POST("/register", registerUserController)

	// Protected Route. Example of how to protect a route (Delete later)
	group.GET("/test", middlewares.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Acceso concedido"})
	})
}
