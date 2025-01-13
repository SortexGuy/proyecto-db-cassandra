package movies

import "github.com/google/uuid"

type MovieDTO struct {
	ID            uuid.UUID `json:"id"`
	Poster_Link   string    `json:"poster_link"`
	Series_Title  string    `json:"series_title"`
	Released_Year int       `json:"released_year"`
	Certificate   string    `json:"certificate"`
	Runtime       string    `json:"runtime"`
	Genre         string    `json:"genre"`
	IMDB_Rating   float64   `json:"imdb_rating"`
	Overview      string    `json:"overview"`
	Meta_score    int       `json:"meta_score"`
	Director      string    `json:"director"`
	Star1         string    `json:"star1"`
	Star2         string    `json:"star2"`
	Star3         string    `json:"star3"`
	Star4         string    `json:"star4"`
	No_Votes      int       `json:"no_votes"`
	Gross         string    `json:"gross"`
}

// MovieByUser  representa la relación entre una película y un usuario
type MovieByUser struct {
	UserID  uuid.UUID `json:"user_id"`
	MovieID uuid.UUID `json:"movie_id"`
}

// Movie representa la estructura de una película
type MovieIDOnly struct {
	ID uuid.UUID `json:"id"`
}

/*
CREATE TABLE app.movies (
    movie_id uuid PRIMARY KEY,
    poster_link text,
    series_title text,
    released_year int,
    certificate text,
    runtime text,
    genre text,
    imdb_rating float,
    overview text,
    meta_score int,
    director text,
    star1 text,
    star2 text,
    star3 text,
    star4 text,
    no_of_votes int,
    gross text
);
*/
