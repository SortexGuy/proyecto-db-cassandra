package users

import (
	"log"
	"time"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
	"github.com/SortexGuy/proyecto-db-cassandra/src/schema"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

// Inserta un usuario en la tabla
func createUserRepository(user User) error {
	session := config.SESSION

	query := `INSERT INTO app.users (id, name, email, password) VALUES (?, ?, ?, ?)`
	user.ID, _ = uuid.NewUUID()
	err := session.Query(query, user.ID, user.Name, user.Email, user.Password).Exec()
	if err != nil {
		log.Println("Error inserting user:", err)
		return err
	}
	return nil
}

func addMovieToUserRepository(userID, movieID uuid.UUID) error {
	rel := schema.DBMovieByUser{}
	// Primero se busca si ya existe la relacion entre usuario y pelicula
	query := `SELECT user_id, movie_id, watched, rewatched FROM app.movies_by_user
		WHERE user_id = ? AND movie_id = ?`
	err := config.SESSION.Query(query, userID, movieID).Scan(
		rel.User_ID, rel.Movie_ID, rel.Watched, rel.Rewatched)
	if err != nil && err != gocql.ErrNotFound {
		log.Println("Error adding movie to user:", err)
		return err
	}

	if err == gocql.ErrNotFound {
		// Si no existe se ingresa a la Base de datos
		query = `INSERT INTO app.movies_by_user (user_id, movie_id, watched, rewatched)
		VALUES (?, ?, ?, ?)`
		err := config.SESSION.Query(query, userID, movieID, time.Now(), time.Now()).Exec()
		if err != nil {
			log.Println("Error adding movie to user:", err)
			return err
		}
	} else {
		// Si existe se actualiza solo el tiempo en el que se volvio a ver
		query = `UPDATE app.movies_by_user SET rewatched = ?
			WHERE user_id = ? AND movie_id = ?`
		err := config.SESSION.Query(query, time.Now(), userID, movieID).Exec()
		if err != nil {
			log.Println("Error adding movie to user:", err)
			return err
		}
	}
	return nil
}

// GetAllUsers obtiene todos los usuarios de la base de datos
func getAllUsersRepository() ([]User, error) {
	session := config.SESSION // Asegúrate de que config.SESSION esté correctamente inicializado
	var users []User
	query := "SELECT user_id, name, email FROM users"

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

	query := `UPDATE app.users SET name = ?, email = ?, password = ? WHERE id = ?`
	err := session.Query(query, user.Name, user.Email, user.Password, user.ID).Exec()
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}
	return nil
}

func deleteUserRepository(userID uuid.UUID) error {
	session := config.SESSION

	query := `DELETE FROM app.users WHERE id = ?`
	err := session.Query(query, userID).Exec()
	if err != nil {
		log.Println("Error deleting user:", err)
		return err
	}
	return nil
}
