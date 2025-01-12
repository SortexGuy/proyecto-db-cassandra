package users

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getUserByIDController(c *gin.Context) {
	userIDText := c.Query("user_id")
	userID, err := strconv.ParseInt(userIDText, 10, 64)
	if userIDText == "" || err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "movie_id not found",
		})
		return
	}

	user, err := getUserByIDService(userID)
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
