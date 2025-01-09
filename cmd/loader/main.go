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
