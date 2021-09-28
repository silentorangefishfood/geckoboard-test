package graph

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
)

type Graph struct {
	Nodes sync.Map

	mu   sync.RWMutex
	size int
}

func (g *Graph) GetSize() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.size
}

func (g *Graph) IncrementSize(inc int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.size += inc
}

type GraphNode struct {
	mu     sync.RWMutex
	Value  interface{}
	weight int
	edges  []*EdgeNode
}

func (n *GraphNode) GetWeight() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.weight
}

func (n *GraphNode) IncrementWeight(inc int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.weight += inc
}

func (n *GraphNode) TotalEdges() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.edges)
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: sync.Map{},
	}
}

type EdgeNode struct {
	Index string

	mu     sync.RWMutex
	weight int
}

func newEdgeNode(index string) *EdgeNode {
	return &EdgeNode{
		Index:  index,
		weight: 1,
	}
}

func (e *EdgeNode) GetWeight() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.weight
}

func (e *EdgeNode) IncrementWeight(inc int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.weight += inc
}

// AddNode adds a new node in the graph
func (g *Graph) AddNode(index string, value interface{}) {
	_, ok := g.Nodes.Load(index)
	if !ok {
		node := &GraphNode{Value: value}
		g.Nodes.Store(index, node)
		g.IncrementSize(1)
		return
	}
}

// AddEdge adds an edge between two existing nodes in the graph
func (g *Graph) AddEdge(n1Index, n2Index string) {
	// Check if the source node exists
	n1, ok := g.Nodes.Load(n1Index)
	if !ok {
		fmt.Println("Source node must exist")
		return
	}

	// Keep track of the total number of edges adjacent to the GraphNode.  We use
	// this value to randomly pick an edge with the same frequency the edge
	// occours in the corpus.
	graphNode := n1.(*GraphNode)
	graphNode.IncrementWeight(1)

	// Check if the edge already exists
	graphNode.mu.RLock()
	defer graphNode.mu.RUnlock()
	for _, edge := range graphNode.edges {
		// if the edge already exists
		if edge.Index == n2Index {
			// Increment the weight (represents frequency)
			edge.IncrementWeight(1)
			return
		}
	}

	// Otherwise, it doesn't exist, create the new edge
	n2 := newEdgeNode(n2Index)

	// Insert the new edge
	graphNode.edges = append(graphNode.edges, n2)
}

func (g *GraphNode) RandomEdge() *EdgeNode {
	stoppingPoint := rand.Intn(g.GetWeight()) + 1
	count := 0

	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, edge := range g.edges {
		count += edge.GetWeight()
		if count >= stoppingPoint {
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
