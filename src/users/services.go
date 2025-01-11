package users

import (
	"log"
	"github.com/SortexGuy/proyecto-db-cassandra/config"
)

// GetUser ByIDService obtiene un usuario por su ID
func GetUserByIDService(id int) (int64, error) {
	users, err := GetAllUsers() // Llama a GetAllUsers directamente
	if err != nil {
		log.Println("Error getting users:", err)
		return 0, err // Retorna un UserDTO vacío y el error
	}

	// Buscar el usuario por ID
	for _, u := range users {
		if u.ID == int64(id) {
			return u.ID, nil // Retorna el ID del usuario encontrado
		}
	}

	return 0, nil // Retorna 0 si no se encuentra el usuario
}

// GetAllUsers obtiene todos los usuarios de la base de datos
func GetAllUsers() ([]User , error) {
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