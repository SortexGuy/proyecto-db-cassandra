package movies

import (
	"log"
	"time"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
)

func createMovieRepository(movie MovieDTO) error {
	session := config.SESSION
	query := `
		INSERT INTO app.movies
		(movie_id, poster_link, series_title, released_year, certificate,
		runtime, genre, imdb_rating, overview, 
		meta_score, director, star1, star2, star3, star4, no_Votes, gross)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	err := session.Query(query,
		movie.ID,
		movie.Poster_Link,
		movie.Series_Title,
		movie.Released_Year,
		movie.Certificate,
		movie.Runtime,
		movie.Genre,
		movie.IMDB_Rating,
		movie.Overview,
		movie.Meta_score,
		movie.Director,
		movie.Star1,
		movie.Star2,
		movie.Star3,
		movie.Star4,
		movie.No_Votes,
		movie.Gross,
	).Exec()

	if err != nil {
		log.Println("Error creating movie:", err)
		return err
	}
	return nil
}

func GetAllMoviesRepository() ([]MovieDTO, error) {
	session := config.SESSION
	var movies []MovieDTO

	// Consulta completa con todos los campos de MovieDTO
	query := `SELECT 
		movie_id, poster_link, series_title, released_year, certificate, 
		runtime, genre, imdb_rating, overview, meta_score, director, 
		star1, star2, star3, star4, no_Votes, gross 
		FROM app.movies`

	iter := session.Query(query).Iter()
	defer iter.Close()

	// Variables temporales para escanear los datos
	var (
		id           int64
		posterLink   string
		seriesTitle  string
		releasedYear int
		certificate  string
		runtime      string
		genre        string
		imdbRating   float64
		overview     string
		metaScore    int
		director     string
		star1        string
		star2        string
		star3        string
		star4        string
		noVotes      int
		gross        string
	)

	// Iterar sobre los resultados y asignar a MovieDTO
	for iter.Scan(
		&id, &posterLink, &seriesTitle, &releasedYear, &certificate,
		&runtime, &genre, &imdbRating, &overview, &metaScore, &director,
		&star1, &star2, &star3, &star4, &noVotes, &gross,
	) {
		movie := MovieDTO{
			ID:            id,
			Poster_Link:   posterLink,
			Series_Title:  seriesTitle,
			Released_Year: releasedYear,
			Certificate:   certificate,
			Runtime:       runtime,
			Genre:         genre,
			IMDB_Rating:   imdbRating,
			Overview:      overview,
			Meta_score:    metaScore,
			Director:      director,
			Star1:         star1,
			Star2:         star2,
			Star3:         star3,
			Star4:         star4,
			No_Votes:      noVotes,
			Gross:         gross,
		}
		movies = append(movies, movie)
	}

	// Verificar errores al cerrar el iterador
	if err := iter.Close(); err != nil {
		log.Println("Error closing iterator:", err)
		return nil, err
	}

	return movies, nil
}

// GetAllMovies obtiene todas las películas de la base de datos
func getAllMoviesIDRepository() ([]int64, error) {
	session := config.SESSION
	var movies []int64
	query := "SELECT movie_id FROM app.movies" // Asegúrate de que este sea el nombre correcto de tu tabla

	iter := session.Query(query).Iter()
	defer iter.Close()

	var movie int64
	for iter.Scan(&movie) {
		movies = append(movies, movie)
	}

	if err := iter.Close(); err != nil {
		log.Println("Error closing iterator:", err)
		return nil, err
	}

	return movies, nil
}

func getMovieByIDRepository(movieID int64) (MovieDTO, error) {
	session := config.SESSION
	query := `SELECT 
		movie_id, poster_link, series_title, released_year, certificate, 
		runtime, genre, imdb_rating, overview, meta_score, director, 
		star1, star2, star3, star4, no_Votes, gross 
	 FROM movies WHERE movie_id = ?`
	var movie MovieDTO

	err := session.Query(query, movieID).Scan(
		&movie.ID,
		&movie.Poster_Link,
		&movie.Series_Title,
		&movie.Released_Year,
		&movie.Certificate,
		&movie.Runtime,
		&movie.Genre,
		&movie.IMDB_Rating,
		&movie.Overview,
		&movie.Meta_score,
		&movie.Director,
		&movie.Star1,
		&movie.Star2,
		&movie.Star3,
		&movie.Star4,
		&movie.No_Votes,
		&movie.Gross,
	)

	if err != nil {
		log.Println("Error fetching movie by ID:", err)
		return movie, err
	}
	return movie, nil
}

func GetAllMoviesByUserRepository() ([]MovieByUser, error) {
	session := config.SESSION
	var moviesByUser []MovieByUser

	query := "SELECT user_id, movie_id, watched FROM movies_by_user"

	iter := session.Query(query).Iter()
	defer iter.Close()

	var watched string
	var movieByUser MovieByUser
	for iter.Scan(&movieByUser.UserID, &movieByUser.MovieID, &watched) {
		movieByUser.Watched, _ = time.Parse(time.Layout, watched)
		moviesByUser = append(moviesByUser, movieByUser)
	}

	if err := iter.Close(); err != nil {
		log.Println("Error closing iterator:", err)
		return nil, err
	}

	return moviesByUser, nil
}

// GetAllMoviesByUser  obtiene todas las películas de un usuario específico
func getMoviesByUserRepository(userID int64) ([]MovieByUser, error) {
	session := config.SESSION
	var moviesByUser []MovieByUser
	query := "SELECT user_id, movie_id, watched FROM movies_by_user WHERE user_id = ?"

	// Ejecuta la consulta con el userID
	iter := session.Query(query, userID).Iter()
	defer iter.Close()

	var watched string
	var movieByUser MovieByUser
	for iter.Scan(&movieByUser.UserID, &movieByUser.MovieID, &watched) {
		movieByUser.Watched, _ = time.Parse(time.Layout, watched)
		moviesByUser = append(moviesByUser, movieByUser)
	}

	if err := iter.Close(); err != nil {
		log.Println("Error closing iterator:", err)
		return nil, err
	}

	return moviesByUser, nil
}

func UpdateMovieRepository(movie MovieDTO) error {
	session := config.SESSION
	query := `
		UPDATE app.movies SET 
		poster_link = ?, series_title = ?, released_year = ?, certificate = ?, runtime = ?, genre = ?, 
		imdb_rating = ?, overview = ?, meta_score = ?, 
		director = ?, star1 = ?, star2 = ?, star3 = ?, star4 = ?, no_of_votes = ?, gross = ?
		WHERE movie_id = ?
	`
	err := session.Query(query,
		movie.Poster_Link, movie.Series_Title, movie.Released_Year, movie.Certificate,
		movie.Runtime, movie.Genre, movie.IMDB_Rating, movie.Overview,
		movie.Meta_score, movie.Director, movie.Star1, movie.Star2,
		movie.Star3, movie.Star4, movie.No_Votes, movie.Gross,
		movie.ID,
	).Exec()

	if err != nil {
		log.Println("Error updating movie:", err)
		return err
	}
	return nil
}

func DeleteMovieRepository(movieID int64) error {
	session := config.SESSION
	query := `DELETE FROM app.movies WHERE movie_id = ?`
	err := session.Query(query, movieID).Exec()

	if err != nil {
		log.Println("Error deleting movie:", err)
		return err
	}
	return nil
}
