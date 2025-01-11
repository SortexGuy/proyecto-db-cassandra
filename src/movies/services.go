package movies

func GetMovieByIDService(id int) (MovieDTOS, error) {
	movie, err := FindMovieByIDRepo(id)

	return movie, err
}

func getAllMoviesService() ([]Movie, error) {
	movies, err := getAllMoviesRepository()
	return movies, err
}
