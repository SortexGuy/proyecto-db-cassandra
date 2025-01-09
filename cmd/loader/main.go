package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

type CsvLine struct {
	Movie_ID      int
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
	No_of_Votes   int
	Gross         string
}

type UserCsv struct {
	User_ID  int
	Name     string
	Email    string
	Password string
}

type MovieByUserCsv struct {
	User_ID      int
	Movie_ID     int
	Username     string
	Movie_Title  string
	Director     string
	Release_Date int
}

var SESSION *gocql.Session

func getClusterConfig() *gocql.ClusterConfig {
	cass_ip := os.Getenv("CASSANDRA_IPADDRESS")
	log.Println("Trying to connect to container at ", cass_ip)
	cluster := gocql.NewCluster(cass_ip)
	cluster.Consistency = gocql.Quorum
	return cluster
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Please create a .env file in the root directory of the project")
	}

	cluster := getClusterConfig()

	log.Println("Trying to open Cassandra session")
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Unable to open up a session with the Cassandra database!", err)
	}
	SESSION = session
	result := SESSION.Query(`CREATE KEYSPACE IF NOT EXISTS app WITH REPLICATION = {
	'class' : 'SimpleStrategy',
	'replication_factor' : '1'
};`)
	err = result.Exec()
	if err != nil {
		log.Fatalln(err)
	}

	result = SESSION.Query(`CREATE TABLE IF NOT EXISTS app.movies(
	movie_id int,
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
	result = SESSION.Query(`CREATE TABLE IF NOT EXISTS app.users (
		user_id int,
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
	result = SESSION.Query(`CREATE TABLE IF NOT EXISTS app.movies_by_user (
				user_id int,
				movie_id int,
				username text,
				movie_title text,
				director text,
				release_date int,
				PRIMARY KEY( user_id,movie_id )
			);`)
	err = result.Exec()
	if err != nil {
		log.Fatalln(err)
	}

	// result = SESSION.Query(`CREATE INDEX movie_id ON app.movies(movie_id);`)
	// err = result.Exec()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	log.Println("Database Setup Finished")

	err = readFromCSVFile("./data/movies_copy.csv")
	if err != nil {
		log.Fatalln(err)
	}

	err = ReadUsersCSV("./data/usuarios.csv")
	if err != nil {
		log.Fatalln(err)
	}

	err = readUsersMovies("./data/peliculas_vistas.csv")
	if err != nil {

		log.Fatalln(err)
	}

	log.Println("All good")

	defer SESSION.Close()
}

// 1. Read a CSV file line-by-line (from local file)
func readFromCSVFile(filepath string) (err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	log.Println("File opened")

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

		processRecord(record)
		count += 1
	}
	log.Println("Records processed: ", count)

	return nil
}

// Función para leer el archivo CSV de usuarios
func ReadUsersCSV(filepath string) (err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	log.Println("Users file opened")

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

		ProcessUss(record) // Procesar el registro de usuarios
		count += 1
	}
	log.Println("Users records processed: ", count)

	return nil
}

// Función para leer el archivo CSV de usuarios
func readUsersMovies(filepath string) (err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	log.Println("Users movies file opened")

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

		ProcessUssMovies(record) // Procesar el registro de usuarios x pelicula 
		count += 1
	}
	log.Println("Users movies records processed: ", count)

	return nil
}

// 2. Process the records in each line
func processRecord(line []string) {
	if len(line) < 16 {
		log.Println("Invalid length, discarding line...")
		return
	}

	if line[0] == "Movie_ID" {
		// ignore first line
		return
	}

	// note: error checking omitted for brevity
	Movie_ID, _ := strconv.Atoi(line[0])
	Poster_Link := line[1]
	Series_Title := line[2]
	Released_Year, _ := strconv.Atoi(line[3])
	Certificate := line[4]
	Runtime := line[5]
	Genre := line[6]
	IMDB_Rating, _ := strconv.ParseFloat(line[7], 64)
	Overview := line[8]
	Meta_score, _ := strconv.Atoi(line[9])
	Director := line[10]
	Star1 := line[11]
	Star2 := line[12]
	Star3 := line[13]
	Star4 := line[14]
	No_of_Votes, _ := strconv.Atoi(line[15])
	Gross := line[16]

	buf := CsvLine{
		Movie_ID:      Movie_ID,
		Poster_Link:   Poster_Link,
		Series_Title:  Series_Title,
		Released_Year: Released_Year,
		Certificate:   Certificate,
		Runtime:       Runtime,
		Genre:         Genre,
		IMDB_Rating:   IMDB_Rating,
		Overview:      Overview,
		Meta_score:    Meta_score,
		Director:      Director,
		Star1:         Star1,
		Star2:         Star2,
		Star3:         Star3,
		Star4:         Star4,
		No_of_Votes:   No_of_Votes,
		Gross:         Gross,
	}

	insertIntoDb(buf)
}

// 2. Process the records in each line for users
func ProcessUss(line []string) {
	if len(line) < 4 {
		log.Println("Invalid length, discarding line...")
		return
	}

	if line[0] == "ID" {
		// ignore first line
		return
	}

	// note: error checking omitted for brevity
	User_ID, _ := strconv.Atoi(line[0])
	Name := line[1]
	Email := line[2]
	Password := line[3]

	buff := UserCsv{
		User_ID:  User_ID,
		Name:     Name,
		Email:    Email,
		Password: Password,
	}

	insertIntoDbUss(buff)
}

func ProcessUssMovies(line []string) {
	if len(line) < 6 {
		log.Println("Invalid length, discarding line...")
		return
	}

	if line[0] == "user_id" {
		// ignore first line
		return
	}

	// note: error checking omitted for brevity
	User_ID, _ := strconv.Atoi(line[0])
	Movie_ID, _ := strconv.Atoi(line[1])
	Username := line[2]
	Movie_Title := line[3]
	Director := line[4]
	Release_Date, _ := strconv.Atoi(line[5])

	buff := MovieByUserCsv{
		User_ID:      User_ID,
		Movie_ID:     Movie_ID,
		Username:     Username,
		Movie_Title:  Movie_Title,
		Director:     Director,
		Release_Date: Release_Date,
	}

	insertIntoDbUssMovies(buff)
}

// 3. Insert the values into the database
func insertIntoDb(record CsvLine) {
	query_obj := SESSION.Query(`INSERT INTO app.movies
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
		record.No_of_Votes,
		record.Gross,
	)
	err := query_obj.Exec()

	if err != nil {
		log.Printf("Insert failed: %s\n", err)
		log.Println("Failed record: ", record)
	}
}

func insertIntoDbUss(record UserCsv) {
	query_obj := SESSION.Query(`INSERT INTO app.users
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

func insertIntoDbUssMovies(record MovieByUserCsv) {
	query_obj := SESSION.Query(`INSERT INTO app.movies_by_user
	(user_id, movie_id, username, movie_title, director, release_date)
	VALUES (?, ?, ?, ?, ?, ?)`,
		record.User_ID,
		record.Movie_ID,
		record.Username,
		record.Movie_Title,
		record.Director,
		record.Release_Date,
	)
	err := query_obj.Exec()

	if err != nil {
		log.Printf("Insert failed: %s\n", err)
		log.Println("Failed record: ", record)
	}
}
