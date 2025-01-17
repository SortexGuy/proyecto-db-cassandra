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
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type MovieCSV struct {
	Movie_ID      int64
	Poster_Link   string
	Series_Title  string
	Released_Year int
	Certificate   string
	Runtime       string
	Genre         string
	IMDB_Rating   float64
	Overview      string
	Meta_score    int
	Director      string
	Star1         string
	Star2         string
	Star3         string
	Star4         string
	No_Votes   int
	Gross         string
}

type UserCsv struct {
	User_ID  int64
	Name     string
	Email    string
	Password string
}

type MovieByUserCsv struct {
	User_ID  int64
	Movie_ID int64
	Watched  time.Time
}

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
	defer config.SESSION.Close()

	dbQueries := []string{
		// Crear keyspace
		`CREATE KEYSPACE IF NOT EXISTS app WITH REPLICATION = {
			'class' : 'SimpleStrategy',
			'replication_factor' : '1'
		};`,
		// Crear tabla de peliculas
		`CREATE TABLE IF NOT EXISTS app.movies(
			movie_id bigint,
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
			no_Votes int,
			gross text,
			PRIMARY KEY( (movie_id), genre, series_title )
		);`,
		// Crear tabla de usuarios
		`CREATE TABLE IF NOT EXISTS app.users (
			user_id bigint,
			name text,
			email text,
			password text,
			PRIMARY KEY( (user_id), name, email )
		);`,
		// Crear tabla de peliculas por usuario
		// los id deben ser bigint
		`CREATE TABLE IF NOT EXISTS app.movies_by_user (
			user_id bigint,
			movie_id bigint,
			watched timestamp,
			PRIMARY KEY( (user_id), watched, movie_id )
		);`,
		// Crear tabla para contar los identificadores de varias tablas
		`CREATE TABLE IF NOT EXISTS app.counters (
			id_name text PRIMARY KEY,
			current_id bigint
		);`,
		// Crear tabla para guardar las recomendaciones
		`CREATE TABLE IF NOT EXISTS app.recommendations (
			user_id bigint,
			movie_id bigint,
			PRIMARY KEY( (user_id), movie_id )
		);`,
	}

	for _, query := range dbQueries {
		result := config.SESSION.Query(query)
		err = result.Exec()
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("Database Setup Finished")

	// err = readMoviesFromCSV("./data/movies_copy.csv")
	err = readFromCSV("./data/movies_copy.csv", processMovieRecord)
	if err != nil {
		log.Fatalln(err)
	}

	// err = readUsersFromCSV("./data/usuarios.csv")
	err = readFromCSV("./data/usuarios.csv", processUserRecord)
	if err != nil {
		log.Fatalln(err)
	}

	// err = readUsersMoviesFromCSV("./data/peliculas_vistas.csv")
	err = readFromCSV("./data/peliculas_vistas.csv", processUserMoviesRecord)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("All good")
}

// 1. Read a CSV file line-by-line (from local file)
func readFromCSV(filepath string, processFunc func([]string)) (err error) {
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

		processFunc(record) // Procesar el registro
		count += 1
	}
	log.Println("Records processed: ", count)

	return nil
}

// 2. Process the records in each line
func processMovieRecord(line []string) {
	if len(line) < 16 {
		log.Println("Invalid length, discarding line...")
		return
	}
	if line[0] == "Movie_ID" {
		// ignore first line
		return
	}

	// note: error checking omitted for brevity
	Movie_ID, _ := counters.IncrementCounter("movies")
	// Movie_ID, _ := strconv.ParseInt(line[0], 10, 64)
	Released_Year, _ := strconv.Atoi(line[3])
	IMDB_Rating, _ := strconv.ParseFloat(line[7], 64)
	Meta_score, _ := strconv.Atoi(line[9])
	No_of_Votes, _ := strconv.Atoi(line[15])

	buf := MovieCSV{
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
		No_Votes:   No_of_Votes,
		Gross:         line[16],
	}

	insertMovieIntoDb(buf)
	return
}

// 2. Process the records in each line for users
func processUserRecord(line []string) {
	if len(line) < 4 {
		log.Println("Invalid length, discarding line...")
		return
	}

	if line[0] == "ID" {
		// ignore first line
		return
	}

	// note: error checking omitted for brevity
	User_ID, _ := counters.IncrementCounter("users")
	// User_ID, _ := strconv.ParseInt(line[0], 10, 64)
	Name := line[1]
	Email := line[2]
	Password, _ := bcrypt.GenerateFromPassword([]byte(line[3]), bcrypt.DefaultCost)

	buff := UserCsv{
		User_ID:  User_ID,
		Name:     Name,
		Email:    Email,
		Password: string(Password),
	}

	insertUserIntoDb(buff)
}

func processUserMoviesRecord(line []string) {
	if len(line) < 6 {
		log.Println("Invalid length, discarding line...")
		return
	}

	if line[0] == "user_id" {
		// ignore first line
		return
	}

	// note: error checking omitted for brevity
	User_ID, _ := strconv.ParseInt(line[0], 10, 64)
	Movie_ID, _ := strconv.ParseInt(line[1], 10, 64)
	Watched := time.Now().Add(time.Duration(-24*7*rand.Intn(12)) * time.Hour)

	buff := MovieByUserCsv{
		User_ID:  User_ID,
		Movie_ID: Movie_ID,
		Watched:  Watched,
	}

	insertUserMoviesIntoDb(buff)
}

// 3. Insert the values into the database
func insertMovieIntoDb(record MovieCSV) {
	query_obj := config.SESSION.Query(`INSERT INTO app.movies
	(movie_id, poster_link, series_title, released_year, certificate, runtime, genre, imdb_rating, overview, meta_score, director, star1, star2, star3, star4, no_of_votes, gross)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		record.Movie_ID,
		record.Poster_Link,
		record.Series_Title,
		record.Released_Year,
		record.Certificate,
		record.Runtime,
		record.Genre,
		record.IMDB_Rating,
		record.Overview,
		record.Meta_score,
		record.Director,
		record.Star1,
		record.Star2,
		record.Star3,
		record.Star4,
		record.No_Votes,
		record.Gross,
	)
	err := query_obj.Exec()

	if err != nil {
		log.Printf("Insert failed: %s\n", err)
		log.Println("Failed record: ", record)
	}
}

func insertUserIntoDb(record UserCsv) {
	query_obj := config.SESSION.Query(`INSERT INTO app.users
	(user_id, name, Email, Password)
	VALUES (?, ?, ?, ?)`,
		record.User_ID,
		record.Name,
		record.Email,
		record.Password,
	)
	err := query_obj.Exec()

	if err != nil {
		log.Printf("Insert failed: %s\n", err)
		log.Println("Failed record: ", record)
	}
}

func insertUserMoviesIntoDb(record MovieByUserCsv) {
	query_obj := config.SESSION.Query(`INSERT INTO app.movies_by_user
	(user_id, movie_id, username, movie_title, director, release_date)
	VALUES (?, ?, ?, ?, ?, ?)`,
		record.User_ID,
		record.Movie_ID,
		record.Watched,
	)
	err := query_obj.Exec()

	if err != nil {
		log.Printf("Insert failed: %s\n", err)
		log.Println("Failed record: ", record)
	}
}
