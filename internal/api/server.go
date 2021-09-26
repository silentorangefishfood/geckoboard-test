package api

import "github.com/silentorangefishfood/geckoboard-test/internal/db"

type Server struct {
	Corpus *db.Corpus
}

func NewServer() *Server {
	return &Server{
		Corpus: db.NewCorpus(),
	}
}
