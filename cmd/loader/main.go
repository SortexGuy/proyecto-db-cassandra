package main

import (
	"encoding/csv"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
	"github.com/SortexGuy/proyecto-db-cassandra/src/counters"
	"github.com/SortexGuy/proyecto-db-cassandra/src/schema"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Please create a .env file in the root directory of the project")
	}

	cluster := config.GetClusterConfig(true)

	log.Println("Trying to open Cassandra session")
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Unable to open up a session with the Cassandra database!", err)
	}
	config.SESSION = session

	// Crear keyspace
	result := config.SESSION.Query(`CREATE KEYSPACE IF NOT EXISTS app WITH REPLICATION = {
		'class' : 'SimpleStrategy',
		'replication_factor' : '1'
	};`)
	err = result.Exec()
	if err != nil {
		log.Fatalln(err)
	}

	// Crear tabla de peliculas
	result = config.SESSION.Query(`CREATE TABLE IF NOT EXISTS app.movies(
		movie_id uuid,
		poster_link text,
		series_title text,
		released_year int,
		certificate text,
		runtime text,
		genre text,
		imdb_rating double,
		overview text,
		meta_score int,
		director text,
		star1 text,
		star2 text,
		star3 text,
		star4 text,
		no_of_votes int,
		gross text,
		PRIMARY KEY( movie_id )
	);`)
	err = result.Exec()
	if err != nil {
		log.Fatalln(err)
	}

	// Crear tabla de usuarios
	result = config.SESSION.Query(`CREATE TABLE IF NOT EXISTS app.users (
		user_id uuid,
		name text,
		email text,
		password text,
		PRIMARY KEY( user_id )
	);`)
	err = result.Exec()
	if err != nil {
		log.Fatalln(err)
	}

	// Crear tabla de peliculas por usuario
	//los id deben ser bigint
	result = config.SESSION.Query(`CREATE TABLE IF NOT EXISTS app.movies_by_user (
		user_id uuid,
		movie_id uuid,
		watched timestamp,
		rewatched timestamp,
		PRIMARY KEY( user_id, movie_id )
	);`)
	err = result.Exec()
	if err != nil {
		log.Fatalln(err)
	}

	// Crear tabla para contar los identificadores de varias tablas
	result = config.SESSION.Query(`CREATE TABLE app.counters (
		id_name text PRIMARY KEY,
		current_id bigint
	);`)
	err = result.Exec()
	if err != nil {
		log.Fatalln(err)
	}

	// Crear tabla para guardar las recomendaciones
	result = config.SESSION.Query(`CREATE TABLE app.recommendations (
		user_id uuid,
		movie_id uuid,
		PRIMARY KEY( (user_id), movie_id )
	);`)
	err = result.Exec()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Database Setup Finished")

	err = readMoviesFromCSV("./data/movies_copy.csv")
	if err != nil {
		log.Fatalln(err)
	}

	err = readUsersFromCSV("./data/usuarios.csv")
	if err != nil {
		log.Fatalln(err)
	}

	// err = readUsersMoviesFromCSV("./data/peliculas_vistas.csv")
	err = generateMoviesByUsers()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("All good")

	defer config.SESSION.Close()
}

// 1. Read a CSV file line-by-line (from local file)
func readMoviesFromCSV(filepath string) (err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	count := -1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
			continue
		}

		err = processMovieRecord(record)
		if err != nil {
			log.Println("Processing movies failed")
			return err
		}
		count += 1
	}
	log.Println("Records processed: ", count)

	return nil
}

// Función para leer el archivo CSV de usuarios
func readUsersFromCSV(filepath string) (err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	count := -1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
			continue
		}

		processUserRecord(record) // Procesar el registro de usuarios
		count += 1
	}
	log.Println("Users records processed: ", count)

	return nil
}

// Función para leer el archivo CSV de usuarios
// func readUsersMoviesFromCSV(filepath string) (err error) {
// 	file, err := os.Open(filepath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
//
// 	reader := csv.NewReader(file)
// 	count := -1
// 	for {
// 		record, err := reader.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatalln(err)
// 			continue
// 		}
//
// 		processUserMoviesRecord(record) // Procesar el registro de usuarios x pelicula
// 		count += 1
// 	}
// 	log.Println("Users movies records processed: ", count)
//
// 	return nil
// }

func generateMoviesByUsers() error {
	var err error
	var movieID uuid.UUID
	var movieIDs uuid.UUIDs
	iter := config.SESSION.Query(`SELECT movie_id FROM app.movies`).Iter()
	for iter.Scan(&movieID) {
		movieIDs = append(movieIDs, movieID)
	}
	if err = iter.Close(); err != nil {
		log.Println("Error on scanning movies")
		return err
	}

	var userID uuid.UUID
	var userIDs uuid.UUIDs
	iter = config.SESSION.Query(`SELECT user_id FROM app.users`).Iter()
	for iter.Scan(&userID) {
		userIDs = append(userIDs, userID)
	}
	if err = iter.Close(); err != nil {
		log.Println("Error on scanning users")
		return err
	}

	err = processMoviesByUser(movieIDs, userIDs)
	return err
}

func processMoviesByUser(movieIDs, userIDs uuid.UUIDs) error {
	for _, userID := range userIDs {
		p := rand.Perm(len(movieIDs))
		for _, movieID := range p[:5] {
			timeframeSec := 24 * 7 * rand.Intn(12) * -int(time.Hour)
			watched := time.Now().Add(time.Duration(timeframeSec))
			max := time.Now().Unix()
			min := watched.Unix()
			rewatched := time.Unix(rand.Int63n(max-min)+min, 0)

			record := schema.DBMovieByUser{
				User_ID: userID, Movie_ID: movieIDs[movieID],
				Watched: watched, Rewatched: rewatched,
			}

			err := insertMoviesByUserIntoDb(record)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 2. Process the records in each line
func processMovieRecord(line []string) error {
	if len(line) < 16 {
		log.Println("Invalid length, discarding line...")
		return nil
	}
	if line[0] == "Movie_ID" {
		// ignore first line
		return nil
	}

	// note: error checking omitted for brevity
	// Movie_ID, _ := strconv.Atoi(line[0])
	Movie_ID, _ := uuid.NewUUID()
	Released_Year, _ := strconv.Atoi(line[3])
	IMDB_Rating, _ := strconv.ParseFloat(line[7], 64)
	Meta_score, _ := strconv.Atoi(line[9])
	No_of_Votes, _ := strconv.Atoi(line[15])

	buf := schema.DBMovie{
		Movie_ID:      Movie_ID,
		Poster_Link:   line[1],
		Series_Title:  line[2],
		Released_Year: Released_Year,
		Certificate:   line[4],
		Runtime:       line[5],
		Genre:         line[6],
		IMDB_Rating:   IMDB_Rating,
		Overview:      line[8],
		Meta_score:    Meta_score,
		Director:      line[10],
		Star1:         line[11],
		Star2:         line[12],
		Star3:         line[13],
		Star4:         line[14],
		No_of_Votes:   No_of_Votes,
		Gross:         line[16],
	}

	err := insertMovieIntoDb(buf)
	if err != nil {
		log.Println("Error inserting movie to db")
	}
	return err
}

// 2. Process the records in each line for users
func processUserRecord(line []string) error {
	if len(line) < 4 {
		log.Println("Invalid length, discarding line...")
		return nil
	}
	if line[0] == "ID" {
		// ignore first line
		return nil
	}

	// note: error checking omitted for brevity
	User_ID, _ := uuid.NewUUID()
	buff := schema.DBUser{
		User_ID:  User_ID,
		Name:     line[1],
		Email:    line[2],
		Password: line[3],
	}

	err := insertUserIntoDb(buff)
	if err != nil {
		log.Println("Error inserting user to db")
	}
	return err
}

// 3. Insert the values into the database
func insertMovieIntoDb(record schema.DBMovie) error {
	counters.IncrementCounter("movies")
	queryText := `INSERT INTO app.movies
		(movie_id, poster_link, series_title, released_year,
			certificate, runtime, genre, imdb_rating, overview,
			meta_score, director, star1, star2, star3, star4, no_of_votes, gross)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	query_obj := config.SESSION.Query(queryText,
		record.Movie_ID, record.Poster_Link, record.Series_Title, record.Released_Year,
		record.Certificate, record.Runtime, record.Genre, record.IMDB_Rating,
		record.Overview, record.Meta_score, record.Director, record.Star1,
		record.Star2, record.Star3, record.Star4, record.No_of_Votes, record.Gross,
	)
	err := query_obj.Exec()

	if err != nil {
		log.Println("Failed record: ", record)
	}
	return err
}

func insertUserIntoDb(record schema.DBUser) error {
	counters.IncrementCounter("users")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(record.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed hashing pass")
		return err
	}

	query_obj := config.SESSION.Query(`INSERT INTO app.users
	(user_id, name, Email, Password)
	VALUES (?, ?, ?, ?)`,
		record.User_ID, record.Name, record.Email, hashedPassword,
	)
	err = query_obj.Exec()

	if err != nil {
		log.Println("Failed record: ", record)
	}
	return err
}

func insertMoviesByUserIntoDb(record schema.DBMovieByUser) error {
	query_obj := config.SESSION.Query(`INSERT INTO app.movies_by_user
		(user_id, movie_id, watched, rewatched) VALUES (?, ?, ?, ?)`,
		record.User_ID, record.Movie_ID,
		record.Watched, record.Rewatched,
	)
	err := query_obj.Exec()

	if err != nil {
		log.Println("Failed record: ", record)
	}
	return err
}
