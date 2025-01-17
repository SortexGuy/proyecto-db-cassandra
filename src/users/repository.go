package users

import (
	"log"
	"time"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
)

// Inserta un usuario en la tabla
func createUserRepository(user User) error {
	session := config.SESSION

	query := `INSERT INTO app.users (user_id, name, email, password) VALUES (?, ?, ?, ?)`
	err := session.Query(query, user.ID, user.Name, user.Email, user.Password).Exec()
	if err != nil {
		log.Println("Error inserting user:", err)
		return err
	}
	return nil
}

func addMovieToUserRepository(userID int64, movieID int64) error {
	session := config.SESSION

	query := `
        INSERT INTO app.movies_by_user (user_id, movie_id, watched)
        VALUES (?, ?, ?)
    `
	err := session.Query(query, userID, movieID, time.Now().Format(time.Layout)).Exec()
	if err != nil {
		log.Println("Error adding movie to user:", err)
		return err
	}
	return nil
}

// GetAllUsers obtiene todos los usuarios de la base de datos
func getAllUsersRepository() ([]User, error) {
	session := config.SESSION // Asegúrate de que config.SESSION esté correctamente inicializado
	var users []User
	query := "SELECT user_id, name, email FROM app.users"

	iter := session.Query(query).Iter()
	defer iter.Close()

	var user User
	for iter.Scan(&user.ID, &user.Name, &user.Email) {
		// Agregar el usuario al slice
		users = append(users, user)
	}

	if err := iter.Close(); err != nil {
		log.Println("Error closing iterator:", err)
		return nil, err
	}

	return users, nil
}

// UpdateUserRepository actualiza un usuario por ID
func updateUserRepository(user User) error {
	session := config.SESSION

	query := `UPDATE app.users SET name = ?, email = ?, password = ? WHERE user_id = ?`
	err := session.Query(query, user.Name, user.Email, user.Password, user.ID).Exec()
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}
	return nil
}

func deleteUserRepository(userID int64) error {
	session := config.SESSION

	query := `DELETE FROM app.users WHERE user_id = ?`
	err := session.Query(query, userID).Exec()
	if err != nil {
		log.Println("Error deleting user:", err)
		return err
	}
	return nil
}
