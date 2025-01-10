package movies

import (
	"log"
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

	data := []MovieDTOS{movie}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// // MovieController es la estructura que maneja las operaciones relacionadas con las películas
type MovieController struct {
	movieRepo       *MovieRepository
	movieByUserRepo *MovieByUserRepository // Repositorio para movies_by_user
}


// NewMovieController crea una nueva instancia de MovieController
func NewMovieController(movieRepo *MovieRepository, movieByUserRepo *MovieByUserRepository) *MovieController {
	return &MovieController{
		movieRepo:       movieRepo,
		movieByUserRepo: movieByUserRepo,
	}
}

// GetMovies maneja la solicitud para obtener todas las películas
func (ctrl *MovieController) GetMovies() ([]Movie, error) {
	movies, err := ctrl.movieRepo.GetAllMovies()
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
func (ctrl *MovieController) GetMoviesByUser (userID int64) ([]MovieByUser , error) {
    moviesByUser , err := ctrl.movieByUserRepo.GetAllMoviesByUser (userID)
    if err != nil {
        log.Println("Error fetching movies by user:", err)
        return nil, err
    }

    // Contabilizar cuántas películas se extrajeron
    count := len(moviesByUser )
    log.Printf("Total movies by user %d extracted: %d\n", userID, count)

    return moviesByUser , nil
}