package api

import "github.com/silentorangefishfood/geckoboard-test/internal/db"

// Server holds any shared data the handlers should need.
// The Corpus field is emulating an in memory database.
type Server struct {
	Corpus *db.Corpus
}

// NewServer returns a new server object to define handlers on.
func NewServer() *Server {
	return &Server{
		Corpus: db.NewCorpus(),
	}
}
