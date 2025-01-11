package users

import(


	"github.com/gocql/gocql"

)

// NewUser Repository crea una nueva instancia de UserRepository
func NewUserRepository(session *gocql.Session) *UserRepository {
	return &UserRepository{session: session}
}







