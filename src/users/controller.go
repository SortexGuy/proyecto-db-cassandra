package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUsersController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "users",
	})
}
