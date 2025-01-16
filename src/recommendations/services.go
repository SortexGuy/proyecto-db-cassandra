package recommendations

import (
	"fmt"
	"log"

	"github.com/SortexGuy/proyecto-db-cassandra/src/movies"
	"github.com/SortexGuy/proyecto-db-cassandra/src/users"
)

func makeRecommendationService(userID int64, lambda float64) (Recommendation, error) {
	// Aquí puedes añadir lógica adicional, como validaciones.
	var recomendation Recommendation

	log.Println("Obteniendo todos los usuarios")
	users, err := users.GetAllUserIDsService()
	if err != nil {
		log.Println("Error getting users:", err)
		return recomendation, err
	}
	movieArr, err := movies.GetAllMoviesIDsService()
	if err != nil {
		log.Println("Error getting movies:", err)
		return recomendation, err
	}

	relations, err := movies.GetAllMoviesByUserRepository()
	if err != nil {
		log.Println("Error getting relations:", err)
		return recomendation, err
	}
	grafo := CreateGraph(users, movieArr, relations)

	recomendation = HybridRecommendation(grafo, userID, 10, lambda, 5)

	err = makeRecommendationRepository(recomendation)
	if err != nil {
		log.Println("Error saving recommendation:", err)
		return recomendation, err
	}
	return recomendation, nil
}

func getRecommendationService(userID int64) ([]movies.MovieDTO, error) {
	var m []movies.MovieDTO
	recommendation, err := getRecommendationRepository(userID)
	if err != nil {
		return m, fmt.Errorf("error getting recommendations: %w", err)
	}

	for _, movieID := range recommendation.Movies {
		movie, err := movies.GetMovieByIDService(movieID) // Llama a la función getMovieByIDService
		if err != nil {
			// Manejo de errores al obtener la película.  Puedes registrar el error,
			// devolver un error, o simplemente continuar con la siguiente película.
			// Aquí se opta por continuar con la siguiente película.
			fmt.Printf("Error getting movie %d: %v\n", movieID, err)
			continue
		}
		m = append(m, movie)
	}

	return m, nil
}
