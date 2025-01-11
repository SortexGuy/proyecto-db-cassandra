package movies

import (
	"log"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
)

// findMovieByIDRepo busca una película por su ID
func FindMovieByIDRepo(id int) (MovieDTOS, error) {
	movie := MovieDTOS{}

	// Obtener todas las películas
	movies, err := getAllMoviesRepository()
	if err != nil {
		log.Println("Error getting movies:", err)
		return movie, err // Retorna el MovieDTO vacío y el error
	}

	// Buscar la película por ID
	for _, m := range movies {
		if m.ID == int64(id) { // Asegúrate de que el tipo coincida
			movie.ID = m.ID   // Asigna el ID de la película encontrada
			return movie, nil // Retorna el MovieDTO con el ID encontrado
		}
	}

	return movie, nil // Retorna el MovieDTO vacío si no se encuentra
}

// GetAllMovies obtiene todas las películas de la base de datos
func getAllMoviesRepository() ([]Movie, error) {
	session := config.SESSION
	var movies []Movie
	query := "SELECT movie_id FROM movies" // Asegúrate de que este sea el nombre correcto de tu tabla

	iter := session.Query(query).Iter()
	defer iter.Close()

	var movie Movie
	for iter.Scan(&movie.ID) {
		movies = append(movies, movie)
	}

	if err := iter.Close(); err != nil {
		log.Println("Error closing iterator:", err)
		return nil, err
	}

	return movies, nil
}

// GetAllMoviesByUser  obtiene todas las películas de un usuario específico
func getAllMoviesByUser(userID int64) ([]MovieByUser, error) {
	session := config.SESSION
	var moviesByUser []MovieByUser
	query := "SELECT movie_id, user_id FROM movies_by_user WHERE user_id = ?"

	// Ejecuta la consulta con el userID
	iter := session.Query(query, userID).Iter()
	defer iter.Close()

	var movieByUser MovieByUser
	for iter.Scan(&movieByUser.MovieID, &movieByUser.UserID) {
		moviesByUser = append(moviesByUser, movieByUser)
	}

	if err := iter.Close(); err != nil {
		log.Println("Error closing iterator:", err)
		return nil, err
	}

	return moviesByUser, nil
}
