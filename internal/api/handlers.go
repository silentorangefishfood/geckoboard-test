package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// Learn is a HTTP handler that accepts a stream of plain text.  It cleans and
// normalises this text by replacing all whitespace with single spaces, and
// removing any characters that are not alphanumeric, space, commas, or full
// stops.
// 
// Learn then ingests the data into it's corpus of learned text.
func (s *Server) Learn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	// First even out any whitespace by replacing one or more whitespace with a single space
	re := regexp.MustCompile(`\s+`)
	tidy := re.ReplaceAll(body, []byte(" "))

	// Then remove anything that's not a-z, A-Z, space, commas, or full stops
	re = regexp.MustCompile("[^A-Za-z .,]*")
	tidy = re.ReplaceAll(tidy, []byte(""))

	// Load the cleaned text into the corpus
	s.Corpus.Ingest(tidy)
}

// GenerateResp represents a JSON structure, containing a single string
// sentence generated from the corpus.
type GenerateResp struct {
	Sentence string `json:"sentence"`
}

// Generate is a HTTP handler that generates a new sentence from the corpus,
// and returns it to the caller. It returns 400 in response to an empty corpus,
// indicating the caller should first POST to /learn
func (s *Server) Generate(w http.ResponseWriter, r *http.Request) {
	sentence, err := s.Corpus.Generate(100)
	if err != nil {
		log.Println(err)

		// Must populate corpus before generating
		w.WriteHeader(400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&GenerateResp{
		Sentence: strings.Join(sentence, " "),
	})
}
