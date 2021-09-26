package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func (s *Server) LearnHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		w.WriteHeader(400)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}

	re := regexp.MustCompile("\n")
	tidy := re.ReplaceAll(body, []byte(" "))
	re = regexp.MustCompile("[^A-Za-z .]*")
	tidy = re.ReplaceAll(tidy, []byte(""))
	fmt.Println(string(tidy))
	s.Corpus.Ingest(tidy)
	s.Corpus.Trigrams.Print()
}

func (s *Server) GenerateHandler(w http.ResponseWriter, r *http.Request) {
	sentence, err := s.Corpus.Generate(100)
	if err != nil {
		fmt.Println(err)
		// Must populate corpus before generating
		w.WriteHeader(400)
	}
	fmt.Println(sentence)
}
