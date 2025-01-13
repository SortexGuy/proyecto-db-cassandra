package recommendations

import (
	"log"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
)

func makeRecommendationRepository(recomendation Recommendation) error {
	session := config.SESSION

	// Verificar si ya existe un registro para el usuario
	userID := recomendation.UserID
	movies := recomendation.Movies

	query := `SELECT movie_id FROM recommendations WHERE user_id = ?`
	iter := session.Query(query, userID).Iter()
	movie := movies[0]
	if iter.Scan(&movie) {
		for i, v := range movies {
			if v == movie {
				// Remove item at i
				movies = append(movies[:i], movies[i+1:]...)
				// remove(movies, int64(i))
			}
		}
		updateQuery := `UPDATE recommendations SET movie_id = ? WHERE user_id = ?`
		err := session.Query(updateQuery, movies, userID).Exec()
		if err != nil {
			log.Println("Error updating recommendations:", err)
			return err
		}
	}

	if err := iter.Close(); err != nil || len(movies) > 0 {
		// Si no existe, creamos un nuevo registro
		for _, movieID := range movies {
			insertQuery := `INSERT INTO recommendations (user_id, movie_id) VALUES (?, ?)`
			err := session.Query(insertQuery, userID, movieID).Exec()
			if err != nil {
				log.Println("Error inserting recommendation:", err)
				return err
			}
		}
	}

	return nil
}
