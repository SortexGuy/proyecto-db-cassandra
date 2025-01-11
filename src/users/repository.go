package users

import (
	"log"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
)

// GetAllUsers obtiene todos los usuarios de la base de datos
func getAllUsersRepository() ([]User, error) {
	session := config.SESSION
	var users []User
	query := "SELECT user_id FROM users"

	iter := session.Query(query).Iter()
	defer iter.Close()

	var userID int64
	for iter.Scan(&userID) {
		// Almacenar el ID del usuario
		users = append(users, User{ID: userID})
	}

	if err := iter.Close(); err != nil {
		log.Println("Error closing iterator:", err)
		return nil, err
	}

	return users, nil
}
