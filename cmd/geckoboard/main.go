package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/silentorangefishfood/geckoboard-test/internal/graph"
)

type Corpus struct {
	Trigrams *graph.Graph
}

func NewCorpus() *Corpus {
	return &Corpus{
		Trigrams: graph.NewGraph(),
	}
}

func (c *Corpus) AddTrigram(w1, w2, w3 string) {
	c.Trigrams.AddNode(w1, w2)
	c.Trigrams.AddNode(w2, w3)
	c.Trigrams.AddEdge(w1+w2, w2+w3)
}

// Ingest takes an array of bytes, representing a new body of text to incorperate into the existing corpus
func (c *Corpus) Ingest(bs []byte) {
	arr := strings.Split(string(bs), " ")
	for i := 0; i < len(arr)-3; i++ {
		c.AddTrigram(arr[i], arr[i+1], arr[i+2])
	}
}

func (c *Corpus) Generate(maxLength int) ([]string, error) {
	start, err := c.GetRandomKey()
	if err != nil {
		return []string{}, err
	}
	return c.Trigrams.RandomWalk(start, 0, maxLength), nil
}

func (c *Corpus) GetRandomKey() (string, error) {
	if len(c.Trigrams.Nodes) == 0 {
		return "", errors.New("Empty corpus")
	}
	arr := []string{}
	for k, n := range c.Trigrams.Nodes {
		if len(n.Word1) > 0 && unicode.IsUpper(rune(n.Word1[0])) {
			arr = append(arr, k)
		}
	}

	return arr[rand.Intn(len(arr))], nil
}

func (s *Server) learnHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
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
	default:
		w.WriteHeader(404)
	}
}

func (s *Server) generateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sentence, err := s.Corpus.Generate(100)
		if err != nil {
			fmt.Println(err)
			// Must populate corpus before generating
			w.WriteHeader(400)
		}
		fmt.Println(sentence)
	default:
		w.WriteHeader(404)
	}
}

type Server struct {
	Corpus *Corpus
}

func NewServer() *Server {
	return &Server{
		Corpus: NewCorpus(),
	}
}

func main() {
	s := NewServer()
	http.Handle("/learn", http.HandlerFunc(s.learnHandler))
	http.Handle("/generate", http.HandlerFunc(s.generateHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
