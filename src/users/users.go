package users

type UserDTO struct {
	ID       int64  `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User representa la estructura de un usuario
type User struct {
	ID int64 `json:"id"`
}
