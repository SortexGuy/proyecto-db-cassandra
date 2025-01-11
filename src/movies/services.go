package movies

import(
	"log"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
)

// FindMovieByIDRepo busca una película por su ID y devuelve el ID de la película si se encuentra, o un error.
func FindMovieByIDRepoService(id int) (int64, error) {
    // Obtener todas las películas
    movies, err := getAllMoviesRepoService()
    if err != nil {
        log.Println("Error getting movies:", err)
        return 0, err // Retorna 0 y el error
    }

    // Buscar la película por ID
    for _, m := range movies {
        if m.ID == int64(id) { 
            return m.ID, nil // Retorna el ID de la película encontrada
        }
    }

    return 0, nil // Retorna 0 si no se encuentra la película
}


// GetAllMovies obtiene todas las películas de la base de datos
func getAllMoviesRepoService() ([]Movie, error) {
	session := config.SESSION
	var movies []Movie
	query := "SELECT movie_id FROM movies" 

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
func getAllMoviesByUserRepoService(userID int64) ([]int64, error) {
    var movieIDs []int64  // Slice para almacenar los IDs de las películas
    query := "SELECT movie_id FROM movies_by_user WHERE user_id = ?"

    // Ejecuta la consulta con el userID
    iter := config.SESSION.Query(query, userID).Iter()
    defer iter.Close()

    var movieID int64
    for iter.Scan(&movieID) {
        movieIDs = append(movieIDs, movieID)  // Agrega el ID de la película al slice
    }

    if err := iter.Close(); err != nil {
        log.Println("Error closing iterator:", err)
        return nil, err
    }

    return movieIDs, nil  // Retorna el slice de IDs de películas
}