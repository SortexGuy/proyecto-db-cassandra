package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getAuthController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "auth",
	})
}
