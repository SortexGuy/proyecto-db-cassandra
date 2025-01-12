package counters

import (
	"log"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
)

// Incrementa el contador y devuelve el nuevo valor
func IncrementCounter(idName string) (int64, error) {
	session := config.SESSION
	query := `
        UPDATE app.counters
        SET current_id = current_id + 1
        WHERE id_name = ?
        IF EXISTS
    `
	applied, err := session.Query(query, idName).ScanCAS()
	if err != nil {
		log.Println("Error incrementing counter:", err)
		return 0, err
	}

	if !applied {
		// Si no existe, inicializamos el contador en 1
		err = session.Query("INSERT INTO app.counters (id_name, current_id) VALUES (?, 1)", idName).Exec()
		if err != nil {
			log.Println("Error initializing counter:", err)
			return 0, err
		}
		return 1, nil
	}

	// Recuperar el nuevo valor del contador
	var currentID int64
	err = session.Query("SELECT current_id FROM app.counters WHERE id_name = ?", idName).Scan(&currentID)
	if err != nil {
		log.Println("Error retrieving updated counter:", err)
		return 0, err
	}

	return currentID, nil
}

// Obtiene el valor actual del contador sin incrementarlo
func GetCounter(idName string) (int64, error) {
	session := config.SESSION
	var currentID int64
	err := session.Query("SELECT current_id FROM app.counters WHERE id_name = ?", idName).Scan(&currentID)
	if err != nil {
		log.Println("Error getting counter:", err)
		return 0, err
	}
	return currentID, nil
}
