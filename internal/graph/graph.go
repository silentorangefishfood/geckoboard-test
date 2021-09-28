package graph

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
)

type Graph struct {
	Nodes sync.Map
	Size int
}

type GraphNode struct {
	Value       interface{}
	TotalWeight int
	Edges       []*EdgeNode
}

type EdgeNode struct {
	Index  string
	Weight int
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: sync.Map{},
	}
}

// AddNode adds a new node in the graph
func (g *Graph) AddNode(index string, value interface{}) {
	_, ok := g.Nodes.Load(index)
	if !ok {
		node := &GraphNode{ Value: value }
		g.Nodes.Store(index, node)
		g.Size++
		return
	}
}

// AddEdge adds a node to the
func (g *Graph) AddEdge(n1Index, n2Index string) {
	n1, ok := g.Nodes.Load(n1Index)
	if ok && n1 == nil {
		fmt.Println("Source node must exist")
		return
	}

	// Keep track of the total number of edges adjacent to the GraphNode.  We use
	// this value to randomly pick an edge with the same frequency the edge
	// occours in the corpus.
	graphNode := n1.(*GraphNode)
	graphNode.TotalWeight += 1

	for _, edge := range graphNode.Edges {
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
	graphNode.Edges = append(graphNode.Edges, n2)
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


type StopCase func(int, *GraphNode) bool

// RandomWalk given a starting index will walk through a graph returning a list of strings it comes across
func (g *Graph) RandomWalk(startIndex string, fn StopCase) {
	randomWalk(g, startIndex, 0, fn)
}

func randomWalk(g *Graph, start string, count int, fn StopCase) {
	count += 1
	node, ok := g.Nodes.Load(start)
	if ok {
		graphNode := node.(*GraphNode)
		if ok && fn(count, graphNode) {
			return
		}

		randomWalk(g, graphNode.RandomEdge().Index, count, fn)
	}
}

func (g *Graph) Print() {
	jsonBs, err := json.MarshalIndent(g.Nodes, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonBs))
}
