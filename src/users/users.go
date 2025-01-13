package users

import "github.com/google/uuid"

// User representa la estructura de un usuario
type User struct {
	ID       uuid.UUID `json:"user_id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

/*
CREATE TABLE users (
    id uuid PRIMARY KEY,
    name text,
    email text,
    password text
);
*/
