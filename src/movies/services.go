package movies

func GetMovieByIDService(id int) (MovieDTO, error) {
	movie, err := findMovieByIDRepo(id)

	return movie, err
}
