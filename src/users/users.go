package users

// User representa la estructura de un usuario
type User struct {
	ID       int64  `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
CREATE TABLE users (
    id bigint PRIMARY KEY,
    name text,
    email text,
    password text
);
*/
