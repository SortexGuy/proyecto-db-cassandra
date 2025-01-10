package movies

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getMovieController(c *gin.Context) {
	movieID, exists := c.Get("movie_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "movie_id not found",
		})
		return
	}

	movie, err := GetMovieByIDService(movieID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting movie",
		})
		return
	}

	data := []MovieDTO{movie}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
