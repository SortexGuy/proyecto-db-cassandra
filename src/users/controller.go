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

	user, err := GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting user",
		})
		return
	}

	data := []UserDTO{user}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
