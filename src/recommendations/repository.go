package recommendations

import (
	"log"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
)

func makeRecommendationRepository(recomendation Recommendation) error {
	session := config.SESSION

	// Verificar si ya existe un registro para el usuario
	var user_id int64
	var num_recommendations int
	var movies []int64
	query := `SELECT movies, num_recommendations FROM app.recommendations WHERE user_id = ?`
	if err := session.Query(query, user_id).Scan(&movies, &num_recommendations); err == nil {
		updateQuery := `UPDATE recommendations SET movies = ?, num_recommendations = ? WHERE user_id = ?`
		err := session.Query(updateQuery, movies, num_recommendations, user_id).Exec()
		if err != nil {
			log.Println("Error updating recommendations:", err)
			return err
		}
	} else {
		// Si no existe, creamos un nuevo registro
		insertQuery := `INSERT INTO app.recommendations (user_id, movies, num_recommendations) VALUES (?, ?, ?)`
		err := session.Query(insertQuery, user_id, []int64{user_id}, 1).Exec()
		if err != nil {
			log.Println("Error inserting recommendation:", err)
			return err
		}
	}

	return nil
}
