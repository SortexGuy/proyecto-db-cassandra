package users

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func getUsersController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user_id not found",
		})
		return
	}

}