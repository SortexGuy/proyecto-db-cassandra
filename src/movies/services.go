package movies

import (
	"errors"
	"log"

	. "github.com/SortexGuy/proyecto-db-cassandra/src/counters"
)

func createMovieService(movie MovieDTO) error {
	id, err := IncrementCounter("movies")
	if err != nil {
		return err
	}
	movie.ID = id

	return createMovieRepository(movie)
}

func GetAllMoviesService() ([]MovieDTO, error) {
	movies, err := GetAllMoviesRepository()
	return movies, err
}

func GetAllMoviesIDsService() ([]int64, error) {
	movies, err := getAllMoviesIDRepository()

	// Contabilizar cuántas películas se extrajeron
	count := len(movies)
	log.Printf("Total movies extracted: %d\n", count)
	return movies, err
}

func GetMovieByIDService(movieID int64) (MovieDTO, error) {
	if movieID == 0 {
		return MovieDTO{}, errors.New("movie ID is required")
	}

	return getMovieByIDRepository(movieID)
}

// GetMoviesByUser  obtiene todas las películas de un usuario específico
func GetMoviesWatchedByUserService(userID int64) ([]MovieByUser, error) {
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
	if movie.ID == 0 {
		return errors.New("movie ID is required")
	}

	return UpdateMovieRepository(movie)
}

func DeleteMovieService(movieID int64) error {
	if movieID == 0 {
		return errors.New("movie ID is required")
	}

	return DeleteMovieRepository(movieID)
}
