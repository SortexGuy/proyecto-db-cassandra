package recommendations

import (
	"log"
	"time"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
)

func makeRecommendationRepository(recomendation Recommendation) error {
	session := config.SESSION

	// Verificar si ya existe un registro para el usuario
	userID := recomendation.UserID
	movies := recomendation.Movies

	query := `SELECT movie_id FROM recommendations WHERE user_id = ?`
	iter := session.Query(query, userID).Iter()
	movieID := movies[0]
	if iter.Scan(&movieID) {
		for i, v := range movies {
			if v == movieID {
				// Remove item at i
				movies = append(movies[:i], movies[i+1:]...)
				// remove(movies, int64(i))
			}
		}
		updateQuery := `UPDATE recommendations SET rewatched = ? WHERE user_id = ? AND movie_id = ?`
		err := session.Query(updateQuery, time.Now(), userID, movieID).Exec()
		if err != nil {
			log.Println("Error updating recommendations:", err)
			return err
		}
	}

	if err := iter.Close(); err != nil || len(movies) > 0 {
		// Si no existe, creamos un nuevo registro
		for _, movieID := range movies {
			insertQuery := `INSERT INTO recommendations (user_id, movie_id, watched, rewatched) VALUES (?, ?, ?, ?)`
			err := session.Query(insertQuery, userID, movieID, time.Now(), time.Now()).Exec()
			if err != nil {
				log.Println("Error inserting recommendation:", err)
				return err
			}
		}
	}

	return nil
}
