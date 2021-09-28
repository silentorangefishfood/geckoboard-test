package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func (s *Server) Learn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	re := regexp.MustCompile("\n")
	tidy := re.ReplaceAll(body, []byte(" "))
	re = regexp.MustCompile("[^A-Za-z .]*")
	tidy = re.ReplaceAll(tidy, []byte(""))
	s.Corpus.Ingest(tidy)
}

func (s *Server) Generate(w http.ResponseWriter, r *http.Request) {
	sentence, err := s.Corpus.Generate(100)
	if err != nil {
		fmt.Println(err)
		// Must populate corpus before generating
		w.WriteHeader(400)
		return
	}
	fmt.Println(sentence)
}
