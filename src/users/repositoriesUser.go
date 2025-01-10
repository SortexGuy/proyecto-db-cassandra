package users

import (
	"github.com/gocql/gocql"
	"log"
)

// User representa la estructura de un usuario
type User struct {
	ID int64 `json:"id"`
}

// UserRepository es la estructura que maneja la sesión de Cassandra para usuarios
type UserRepository struct {
	session *gocql.Session
}

// NewUser Repository crea una nueva instancia de UserRepository
func NewUserRepository(session *gocql.Session) *UserRepository {
	return &UserRepository{session: session}
}

// GetAllUsers obtiene todos los usuarios de la base de datos
func (repo *UserRepository) GetAllUsers() ([]User , error) {
	var users []User 
	query := "SELECT user_id FROM users" // Asegúrate de que este sea el nombre correcto de tu tabla

	iter := repo.session.Query(query).Iter()
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