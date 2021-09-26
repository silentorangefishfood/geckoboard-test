package graph

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

type Graph struct {
	Nodes map[string]*GraphNode
}

type GraphNode struct {
	Word1 string
	Word2 string
	TotalWeight int
	Edges []*EdgeNode
}

func (g *GraphNode) RandomEdge() *EdgeNode {
	stoppingPoint := rand.Intn(g.TotalWeight) + 1
	fmt.Printf("Stopping point: %d\n", stoppingPoint)
	count := 0
	for _, edge := range g.Edges {
		count += edge.Weight
		if count >= stoppingPoint {
			fmt.Printf("Reached stopping point, count: %d\n", count)
			return edge
		}
	}

	// We should never reach here
	return nil
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
	fmt.Println("Adding Node")
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
	fmt.Println("Adding Edge")
	n1 := g.Nodes[n1Index]
	if n1 == nil {
		fmt.Println("Source node must exist")
		return
	}

	// Keep track of the total number of edges adjacent to the GraphNode.  We use
	// this value to randomly pick an edge with the same frequency the edge
	// occours in the corpus.
	n1.TotalWeight += 1

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
	fmt.Println("RandomWalk")
	count += 1
	strs := []string{}
	startNode := g.Nodes[start]
	if len(startNode.Edges) == 0 ||
		(count >= maxLength &&
			len(startNode.Word2) > 0 &&
			startNode.Word2[len(startNode.Word2)-1:][0] == '.') {
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
