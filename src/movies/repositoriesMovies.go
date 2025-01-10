package movies

import (
	"github.com/gocql/gocql"
	"log"
)

// Movie representa la estructura de una película
type Movie struct {
	ID    int64  `json:"id"`
}

// MovieByUser  representa la relación entre una película y un usuario
type MovieByUser  struct {
	UserID  int64 `json:"user_id"`
    MovieID int64 `json:"movie_id"`
}

// MovieRepository es la estructura que maneja la sesión de Cassandra
type MovieRepository struct {
	session *gocql.Session
}

// NewMovieRepository crea una nueva instancia de MovieRepository
func NewMovieRepository(session *gocql.Session) *MovieRepository {
	return &MovieRepository{session: session}
} 

// GetAllMovies obtiene todas las películas de la base de datos
func (repo *MovieRepository) GetAllMovies() ([]Movie, error) {
	var movies []Movie
	query := "SELECT movie_id FROM movies" // Asegúrate de que este sea el nombre correcto de tu tabla

	iter := repo.session.Query(query).Iter()
	defer iter.Close()

	var movie Movie
	for iter.Scan(&movie.ID,) {
		movies = append(movies, movie)
	}

	if err := iter.Close(); err != nil {
		log.Println("Error closing iterator:", err)
		return nil, err
	}

	return movies, nil
}

// MovieByUser Repository es la estructura que maneja la sesión de Cassandra para movies_by_user
type MovieByUserRepository struct {
    session *gocql.Session
}

// NewMovieByUser Repository crea una nueva instancia de MovieByUser Repository
func NewMovieByUserRepository(session *gocql.Session) *MovieByUserRepository {
    return &MovieByUserRepository{session: session}
}

// GetAllMoviesByUser  obtiene todas las películas de un usuario específico
func (repo *MovieByUserRepository) GetAllMoviesByUser (userID int64) ([]MovieByUser , error) {
    var moviesByUser  []MovieByUser 
    query := "SELECT movie_id, user_id FROM movies_by_user WHERE user_id = ?"

    // Ejecuta la consulta con el userID
    iter := repo.session.Query(query, userID).Iter()
    defer iter.Close()

    var movieByUser  MovieByUser 
    for iter.Scan(&movieByUser .MovieID, &movieByUser .UserID) {
        moviesByUser  = append(moviesByUser , movieByUser )
    }

    if err := iter.Close(); err != nil {
        log.Println("Error closing iterator:", err)
        return nil, err
    }

    return moviesByUser , nil
}