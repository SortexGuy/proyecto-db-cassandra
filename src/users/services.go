package users

import (
	"errors"
	"log"

	. "github.com/SortexGuy/proyecto-db-cassandra/src/counters"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserService(user User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	id, err := IncrementCounter("users")
	if err != nil {
		return err
	}
	user.ID = id
	err = CreateUserRepository(user)
	if err != nil {
		return err
	}

	return nil

}

func AddMovieToUserService(userID int64, movieID int64) error {
	// Validación básica
	if userID == 0 {
		return errors.New("user ID is required")
	}
	if movieID == 0 {
		return errors.New("movie ID is required")
	}

	// Llamar al repositorio para agregar la película
	err := AddMovieToUserRepository(userID, movieID)
	if err != nil {
		log.Println("Error in AddMovieToUserService:", err)
		return err
	}

	return nil
}

// GetUser ByIDService obtiene un usuario por su ID
func getUserByIDService(id int64) (User, error) {
	user := User{}
	users, err := getAllUsersRepository()
	if err != nil {
		log.Println("Error getting users:", err)
		return user, err // Retorna un UserDTO vacío y el error
	}
	// Buscar el usuario por ID
	for _, u := range users {
		if u.ID == id {
			user.ID = u.ID
			user.Name = u.Name
			user.Email = u.Email
			return user, nil // Retorna el ID del usuario encontrado
		}
	}

	return user, errors.New("User Not Found") // Retorna 0 si no se encuentra el usuario
}

func getAllUserIDsService() ([]int64, error) {
	ids := []int64{}
	users, err := getAllUsersRepository()
	if err != nil {
		log.Println("Error getting users:", err)
		return ids, err // Retorna un UserDTO vacío y el error
	}
	for _, u := range users {
		ids = append(ids, u.ID)
	}

	return ids, nil
}

func verifyEmailService(emailText string) (bool, error) {
	exist := false
	users, err := getAllUsersRepository()
	if err != nil {
		log.Println("Error getting users:", err)
		return exist, err
	}
	// Buscar el usuario por ID
	for _, u := range users {
		if u.Email == emailText {
			exist = true
			return exist, nil
		}
	}

	return exist, errors.New("email not found")
}

// Actualiza los datos de un usuario
func UpdateUserService(user User) error {
	// Validar datos del usuario
	if user.ID == 0 {
		return errors.New("user ID is required")
	}
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return errors.New("name, email, and password are required")
	}

	// Llamar al repositorio para actualizar el usuario
	err := UpdateUserRepository(user)
	if err != nil {
		log.Println("Error in UpdateUserService:", err)
		return err
	}

	return nil
}

// Elimina un usuario por ID
func DeleteUserService(userID int64) error {
	// Validar el ID del usuario
	if userID == 0 {
		return errors.New("user ID is required")
	}

	// Llamar al repositorio para eliminar el usuario
	err := DeleteUserRepository(userID)
	if err != nil {
		log.Println("Error in DeleteUserService:", err)
		return err
	}

	return nil
}
