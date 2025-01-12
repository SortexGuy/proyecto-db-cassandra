package movies

import "log"

func getMovieByIDService(id int64) (MovieDTO, error) {
	movie, err := findMovieByIDRepository(id)
	return movie, err
}

func getAllMoviesService() ([]Movie, error) {
	movies, err := getAllMoviesRepository()

	// Contabilizar cuántas películas se extrajeron
	count := len(movies)
	log.Printf("Total movies extracted: %d\n", count)
	return movies, err
}

// GetMoviesByUser  obtiene todas las películas de un usuario específico
func getMoviesByUserService(userID int64) ([]MovieByUser, error) {
	moviesByUser, err := getAllMoviesByUserRepository(userID)
	if err != nil {
		log.Println("Error fetching movies by user:", err)
		return nil, err
	}

	// Contabilizar cuántas películas se extrajeron
	count := len(moviesByUser)
	log.Printf("Total movies by user %d extracted: %d\n", userID, count)

	return moviesByUser, nil
}
