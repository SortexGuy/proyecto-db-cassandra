package movies

import (
	"errors"
	"log"

	. "github.com/SortexGuy/proyecto-db-cassandra/src/counters"
	"github.com/google/uuid"
)

func createMovieService(movie MovieDTO) error {
	_, err := IncrementCounter("movies")
	if err != nil {
		return err
	}
	// movie.ID = id

	return createMovieRepository(movie)
}

func GetAllMoviesService() ([]MovieDTO, error) {
	movies, err := getAllMoviesRepository()
	return movies, err
}

func GetAllMoviesIDsService() (uuid.UUIDs, error) {
	movies, err := getAllMoviesIDRepository()

	// Contabilizar cuántas películas se extrajeron
	count := len(movies)
	log.Printf("Total movies extracted: %d\n", count)
	return movies, err
}

func getMovieByIDService(movieID uuid.UUID) (MovieDTO, error) {
	// if movieID == 0 {
	// 	return MovieDTO{}, errors.New("movie ID is required")
	// }

	return getMovieByIDRepository(movieID)
}

// GetMoviesByUser  obtiene todas las películas de un usuario específico
func GetMoviesWatchedByUserService(userID uuid.UUID) ([]MovieByUser, error) {
	moviesByUser, err := getMoviesByUserRepository(userID)
	if err != nil {
		log.Println("Error fetching movies by user:", err)
		return nil, err
	}

	// Contabilizar cuántas películas se extrajeron
	count := len(moviesByUser)
	log.Printf("Total movies by user %d extracted: %d\n", userID, count)

	return moviesByUser, nil
}

func UpdateMovieService(movie MovieDTO) error {
	// if movie.ID == 0 {
	// 	return errors.New("movie ID is required")
	// }

	return UpdateMovieRepository(movie)
}

func DeleteMovieService(movieID uuid.UUID) error {
	// if movieID == 0 {
	// 	return errors.New("movie ID is required")
	// }

	return DeleteMovieRepository(movieID)
}
