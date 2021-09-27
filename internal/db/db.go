package db

import (
	"errors"
	"math/rand"
	"strings"
	"unicode"

	"github.com/silentorangefishfood/geckoboard-test/internal/graph"
)

type Corpus struct {
	Trigrams *graph.Graph
}

type Bigram struct {
	Word1 string
	Word2 string
}

func NewCorpus() *Corpus {
	return &Corpus{
		Trigrams: graph.NewGraph(),
	}
}

func (c *Corpus) AddTrigram(w1, w2, w3 string) {
	c.Trigrams.AddNode(w1+w2, Bigram{
		Word1: w1,
		Word2: w2,
	})
	c.Trigrams.AddNode(w2+w3, Bigram{
		Word1: w2,
		Word2: w3,
	})
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

	strs := []string{}
	c.Trigrams.RandomWalk(start, func(nodeCount int, node *graph.GraphNode) bool {
		bigram := node.Value.(Bigram)
		strs = append(strs, bigram.Word1)
		if len(node.Edges) == 0 ||
			(nodeCount >= maxLength &&
				len(bigram.Word2) > 0 &&
				bigram.Word2[len(bigram.Word2)-1:][0] == '.') {
			strs = append(strs, bigram.Word2)
			return true
		}

		return false
	})

	return strs, nil
}

func (c *Corpus) GetRandomKey() (string, error) {
	if len(c.Trigrams.Nodes) == 0 {
		return "", errors.New("Empty corpus")
	}
	arr := []string{}
	for k, n := range c.Trigrams.Nodes {
		bigram := n.Value.(Bigram)
		if len(bigram.Word1) > 0 && unicode.IsUpper(rune(bigram.Word1[0])) {
			arr = append(arr, k)
		}
	}

	return arr[rand.Intn(len(arr))], nil
}
