package users

import(
	"github.com/gocql/gocql"
)

type UserDTO struct {
	User_ID  int    `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User representa la estructura de un usuario
type User struct {
	ID int64 `json:"id"`
}

// UserRepository es la estructura que maneja la sesión de Cassandra para usuarios
type UserRepository struct {
	session *gocql.Session
}
