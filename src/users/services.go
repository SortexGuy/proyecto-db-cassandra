package users

func GetUserByIDService(id int) (UserDTO, error) {
	user, err := findUserByIDRepo(id)

	return user, err
}
