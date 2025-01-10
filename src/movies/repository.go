package movies

import (
	"log"
	"github.com/gocql/gocql"
)

// Movie representa la estructura de una película
type Movies struct {
	ID    int64  `json:"id"`
}

// MovieDTO representa la estructura de datos que deseas devolver
type MovieDTOS struct {
	ID int64 `json:"id"` // Asegúrate de que este campo esté presente
}

// MovieRepository es la estructura que maneja la sesión de Cassandra// Variable global para el repositorio
var MovieRepo *MovieRepository

// NewMovieRepositorys crea una nueva instancia de MovieRepository
func NewMovieRepositorys(session *gocql.Session) *MovieRepository {
	return &MovieRepository{session: session}
}


// findMovieByIDRepo busca una película por su ID
func FindMovieByIDRepo(id int) (MovieDTOS, error) {
	movie := MovieDTOS{}

	// Obtener todas las películas
	movies, err := MovieRepo.GetAllMovie() // Usa la variable global
	if err != nil {
		log.Println("Error getting movies:", err)
		return movie, err // Retorna el MovieDTO vacío y el error
	}

	// Buscar la película por ID
	for _, m := range movies {
		if m.ID == int64(id) { // Asegúrate de que el tipo coincida
			movie.ID = m.ID // Asigna el ID de la película encontrada
			return movie, nil // Retorna el MovieDTO con el ID encontrado
		}
	}

	return movie, nil // Retorna el MovieDTO vacío si no se encuentra
}
// GetAllMovies obtiene todas las películas de la base de datos
func (repo *MovieRepository) GetAllMovie() ([]Movie, error) {
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
