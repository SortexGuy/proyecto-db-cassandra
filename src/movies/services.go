package movies

func GetMovieByIDService(id int) (MovieDTOS, error) {
	movie, err := FindMovieByIDRepo(id)

	return movie, err
}
