package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

type Corpus struct {
	Trigrams *Graph
}

func NewCorpus() *Corpus {
	return &Corpus{
		Trigrams: NewGraph(),
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

func (c *Corpus) Generate(maxLength int) []string {
	start := c.GetRandomKey()
	return c.Trigrams.RandomWalk(start, 0, maxLength)
}

func (c *Corpus) GetRandomKey() string {
	arr := []string{}
	for k := range c.Trigrams.Nodes {
		arr = append(arr, k)
	}

	return arr[rand.Intn(len(arr))]
}

type Graph struct {
	Nodes map[string]*GraphNode
}

type GraphNode struct {
	Word1 string
	Word2 string
	Edges []*EdgeNode
}

func (g *GraphNode) RandomEdge() *EdgeNode {
	return g.Edges[rand.Intn(len(g.Edges))]
}

type EdgeNode struct {
	Index  string
	Weight int
}

func (n *GraphNode) GetIndex() string {
	return n.Word1 + n.Word2
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*GraphNode),
	}
}

// AddNode adds a new node in the graph
func (g *Graph) AddNode(word1, word2 string) {
	index := word1 + word2
	node := g.Nodes[index]
	if node == nil {
		node := &GraphNode{
			Word1: word1,
			Word2: word2,
		}
		g.Nodes[index] = node
		return
	}
}

// AddEdge adds a node to the
func (g *Graph) AddEdge(n1Index, n2Index string) {
	n1 := g.Nodes[n1Index]
	if n1 == nil {
		fmt.Println("Source node must exist")
		return
	}

	for _, edge := range n1.Edges {
		// if the edge already exists
		if edge.Index == n2Index {
			// Increment the weight (represents frequency)
			edge.Weight += 1
			return
		}
	}

	// Otherwise, it doesn't exist, create the new edge
	n2 := &EdgeNode{
		Index:  n2Index,
		Weight: 1,
	}

	// Insert the new edge
	g.Nodes[n1Index].Edges = append(g.Nodes[n1Index].Edges, n2)
}

// RandomWalk given a starting index will walk through a graph returning a list of strings it comes across
func (g *Graph) RandomWalk(start string, count, maxLength int) []string {
	count += 1
	strs := []string{}
	startNode := g.Nodes[start]
	if len(startNode.Edges) == 0 || count == maxLength {
		strs = append(strs, startNode.Word2)
		return strs
	}

	strs = append(strs, startNode.Word1)
	randomNext := startNode.RandomEdge()
	randomNextIndex := randomNext.Index
	strs = append(strs, g.RandomWalk(randomNextIndex, count, maxLength)...)
	return strs
}

func (g *Graph) Print() {
	jsonBs, err := json.MarshalIndent(g.Nodes, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonBs))
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

		re := regexp.MustCompile("[^A-Za-z ]*")
		tidy := re.ReplaceAll(body, []byte(""))
		s.Corpus.Ingest(tidy)
	default:
		w.WriteHeader(404)
	}
}

func (s *Server) generateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sentence := s.Corpus.Generate(100)
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
