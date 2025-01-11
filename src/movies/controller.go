package movies

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getMovieByIDController(c *gin.Context) {
	movieIDText := c.Query("movie_id")
	movieID, err := strconv.ParseInt(movieIDText, 10, 64)
	if movieIDText == "" || err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "movie_id not found",
		})
		return
	}

	movie, err := getMovieByIDService(movieID)
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

// GetMovies maneja la solicitud para obtener todas las películas
func getAllMoviesController(c *gin.Context) {
	movies, err := getAllMoviesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch movies"})
		return
	}
	c.JSON(http.StatusOK, movies)

}

func getMovieByUserController(c *gin.Context) {
	userIDText := c.Query("movie_id")
	userID, err := strconv.ParseInt(userIDText, 10, 64)
	if userIDText == "" || err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "movie_id not found",
		})
		return
	}
	movies, err := getMoviesByUserService(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting movie",
		})
		return
	}

	data := movies
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
