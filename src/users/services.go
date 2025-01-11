package users

import (
	"errors"
	"log"
)

// GetUser ByIDService obtiene un usuario por su ID
func getUserByIDService(id int64) (UserDTO, error) {
	user := UserDTO{}
	users, err := getAllUsersRepository() // Llama a GetAllUsers directamente
	if err != nil {
		log.Println("Error getting users:", err)
		return user, err // Retorna un UserDTO vacío y el error
	}

	// Buscar el usuario por ID
	for _, u := range users {
		if u.ID == id {
			user.ID = u.ID

			return user, nil // Retorna el ID del usuario encontrado
		}
	}

	return user, errors.New("User Not Found") // Retorna 0 si no se encuentra el usuario
}
