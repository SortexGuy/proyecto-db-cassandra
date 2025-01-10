package movies

type MovieDTO struct {
	Movie_ID      int     `json:"movie_id"`
	Poster_Link   string  `json:"poster_link"`
	Series_Title  string  `json:"series_title"`
	Released_Year int     `json:"released_year"`
	Certificate   string  `json:"certificate"`
	Runtime       string  `json:"runtime"`
	Genre         string  `json:"genre"`
	IMDB_Rating   float64 `json:"imdb_rating"`
	Overview      string  `json:"overview"`
	Meta_score    int     `json:"meta_score"`
	Director      string  `json:"director"`
	Star1         string  `json:"star1"`
	Star2         string  `json:"star2"`
	Star3         string  `json:"star3"`
	Star4         string  `json:"star4"`
	No_of_Votes   int     `json:"no_of_votes"`
	Gross         string  `json:"gross"`
}
