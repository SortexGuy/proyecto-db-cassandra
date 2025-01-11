package movies

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMovieController(c *gin.Context) {
	movieIDText := c.Query("movie_id")
	movieID, err := strconv.Atoi(movieIDText)
	if movieIDText == "" || err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "movie_id not found",
		})
		return
	}
}