package recommendations

import (
	"log"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
)

func makeRecommendationRepository(recomendation Recommendation) error {
	session := config.SESSION

	// Verificar si ya existe un registro para el usuario
	user_id := recomendation.UserID
	num_recommendations := recomendation.NumRecommendations
	movies := recomendation.Movies
	query := `SELECT movie_id, num_recommendations FROM recommendations WHERE user_id = ?`
	// TODO: esto tiene que ser un iterador porque la base de datos no puede guardar arreglos
	if err := session.Query(query, user_id).Scan(&movies, &num_recommendations); err == nil {
		updateQuery := `UPDATE recommendations SET movies = ?, num_recommendations = ? WHERE user_id = ?`
		err := session.Query(updateQuery, movies, num_recommendations, user_id).Exec()
		if err != nil {
			log.Println("Error updating recommendations:", err)
			return err
		}
	} else {
		for _, movieID := range movies {
			// Si no existe, creamos un nuevo registro
			insertQuery := `INSERT INTO recommendations (user_id, movies, num_recommendations) VALUES (?, ?, ?)`
			err := session.Query(insertQuery, user_id, movieID, 1).Exec()
			if err != nil {
				log.Println("Error inserting recommendation:", err)
				return err
			}
		}
	}

	return nil
}
