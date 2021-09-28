package db

import (
	"errors"
	"math/rand"
	"strings"
	"unicode"

	"github.com/silentorangefishfood/geckoboard-test/internal/graph"
)

// Corpus has a single field, Trigrams, which
// is a graph structure representing a body of text.
type Corpus struct {
	Trigrams *graph.Graph
}

// Bigram holds two words
type Bigram struct {
	Word1 string
	Word2 string
}

// NewCorpus creates a new corpus.
func NewCorpus() *Corpus {
	return &Corpus{
		Trigrams: graph.NewGraph(),
	}
}

// addTrigram takes three strings representing a trigram, and adds it to the
// graph datastructure. Internally this adds two nodes and an 'edge' between
// them. For example, the trigram ["you", "with", "the"] would be represented
// by the following structure
// {
//   "youwith": {
//     "Value": {
//       "Word1": "you",
//       "Word2": "with"
//     },
//     "Edges": [
//       {
//         "Index": "withthe"
//       }
//     ]
//   },
//   "withthe": {
//     "Value": {
//       "Word1": "with",
//       "Word2": "the"
//     },
//   },
// }
func (c *Corpus) addTrigram(w1, w2, w3 string) error {
	c.Trigrams.AddNode(w1+w2, Bigram{
		Word1: w1,
		Word2: w2,
	})
	c.Trigrams.AddNode(w2+w3, Bigram{
		Word1: w2,
		Word2: w3,
	})
	if err := c.Trigrams.AddEdge(w1+w2, w2+w3); err != nil {
		return err
	}

	return nil
}

// Ingest takes an array of bytes, representing a new body of text to
// incorperate into the existing corpus. Responsibility is on the caller to
// clean and normalise any text before calling Ingest.
func (c *Corpus) Ingest(bs []byte) {
	arr := strings.Split(string(bs), " ")
	for i := 0; i < len(arr)-3; i++ {
		c.addTrigram(arr[i], arr[i+1], arr[i+2])
	}
}

// Generate returns a slice representing a generated sentence from the corpus
// of learned text. The function takes a single argument, approximateLength.
// After producing a sentence greater than approximate length, Generate will
// return upon finding a word ending with a full stop.
func (c *Corpus) Generate(approximateLength int) ([]string, error) {
	start, err := c.GetRandomKey()
	if err != nil {
		return []string{}, err
	}

	strs := []string{}
	c.Trigrams.RandomWalk(start, func(nodeCount int, node *graph.Node) bool {
		bigram := node.Value.(Bigram)
		strs = append(strs, bigram.Word1)
		if node.TotalEdges() == 0 ||
			(nodeCount >= approximateLength &&
				len(bigram.Word2) > 0 &&
				bigram.Word2[len(bigram.Word2)-1:][0] == '.') {
			strs = append(strs, bigram.Word2)
			return true
		}

		return false
	})

	return strs, nil
}

// GetRandomKey returns an entrypoint to the corpus by returning a key to the
// map containing graph nodes. Valid keys are those with the first word
// beginning with a capital letter. An empty corpus returns an error.
func (c *Corpus) GetRandomKey() (string, error) {
	if c.Trigrams.GetSize() == 0 {
		return "", errors.New("Empty corpus")
	}
	arr := []string{}
	c.Trigrams.Nodes.Range(func(key, value interface{}) bool {
		graphNode := value.(*graph.Node)
		index := key.(string)
		bigram := graphNode.Value.(Bigram)
		if len(bigram.Word1) > 0 && unicode.IsUpper(rune(bigram.Word1[0])) {
			arr = append(arr, index)
		}

		return true
	})

	return arr[rand.Intn(len(arr))], nil
}
