package schema

import (
	"time"

	"github.com/google/uuid"
)

type DBMovie struct {
	Movie_ID      uuid.UUID
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

type DBUser struct {
	User_ID  uuid.UUID
	Name     string
	Email    string
	Password string
}

type DBMovieByUser struct {
	User_ID   uuid.UUID
	Movie_ID  uuid.UUID
	Watched   time.Time
	Rewatched time.Time
}
