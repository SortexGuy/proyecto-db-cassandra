package movies

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createMovieController(c *gin.Context) {
	var movie MovieDTO
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := createMovieService(movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Movie created successfully"})
}

func getMovieByIDController(c *gin.Context) {
	movieIDText := c.Param("id")
	movieID, err := strconv.ParseInt(movieIDText, 0, 64)
	if movieIDText == "" || err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "movie_id not found",
		})
		return
	}

	movie, err := getMovieByIDService(movieID)
	if err != nil {
		log.Println(err)
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

func updateMovieController(c *gin.Context) {
	var movie MovieDTO
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := UpdateMovieService(movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully"})
}

// GetMovies maneja la solicitud para obtener todas las películas
func getAllMoviesController(c *gin.Context) {
	movies, err := GetAllMoviesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch movies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": movies,
	})
}

func getMovieWatchedByUserController(c *gin.Context) {
	userIDText := c.Query("user_id")
	userID, err := strconv.ParseInt(userIDText, 10, 64)
	if userIDText == "" || err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user_id not found",
		})
		return
	}
	movies, err := GetMoviesWatchedByUserService(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting movie",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": movies,
	})
}

func deleteMovieController(c *gin.Context) {
	movieID, err := strconv.ParseInt(c.Param("movie_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	err = DeleteMovieService(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}
