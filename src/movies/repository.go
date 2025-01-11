package movies

import(
    "github.com/gocql/gocql" 
)


// NewMovieRepository crea una nueva instancia de MovieRepository
func NewMovieRepository(session *gocql.Session) *MovieRepository {
	return &MovieRepository{session: session}
} 


// NewMovieByUser Repository crea una nueva instancia de MovieByUser Repository
func NewMovieByUserRepository(session *gocql.Session) *MovieByUserRepository {
    return &MovieByUserRepository{session: session}
}


