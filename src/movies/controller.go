package movies

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getMovieController(c *gin.Context) {
	movieIDText := c.Query("movie_id")
	movieID, err := strconv.Atoi(movieIDText)
	if movieIDText == "" || err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "movie_id not found",
		})
		return
	}

	movie, err := GetMovieByIDService(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting movie",
		})
		return
	}

	data := []MovieDTOS{movie}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// GetMovies maneja la solicitud para obtener todas las películas
func getAllMoviesController() ([]Movie, error) {
	movies, err := getAllMoviesService()
	if err != nil {
		log.Println("Error fetching movies:", err)
		return nil, err
	}

	// Contabilizar cuántas películas se extrajeron
	count := len(movies)
	log.Printf("Total movies extracted: %d\n", count)
	return movies, nil
}

// GetMoviesByUser  obtiene todas las películas de un usuario específico
func GetMoviesByUser(userID int64) ([]MovieByUser, error) {
	moviesByUser, err := getAllMoviesByUser(userID)
	if err != nil {
		log.Println("Error fetching movies by user:", err)
		return nil, err
	}

	// Contabilizar cuántas películas se extrajeron
	count := len(moviesByUser)
	log.Printf("Total movies by user %d extracted: %d\n", userID, count)

	return moviesByUser, nil
}
