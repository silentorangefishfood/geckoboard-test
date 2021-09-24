package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
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

func (c *Corpus) Ingest(s string) {
	arr := strings.Split(s, " ")
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

func main() {
	c := NewCorpus()
	bs, err := os.ReadFile("pg66366.txt")
	if err != nil {
		fmt.Println(err)
	}

	re := regexp.MustCompile("[^A-Za-z ]*")
	tidy := re.ReplaceAll(bs, []byte(""))

	c.Ingest(string(tidy))
	sentence := c.Generate(100)
	fmt.Println(sentence)
}
