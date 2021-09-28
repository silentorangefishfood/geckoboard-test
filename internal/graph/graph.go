package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

// Graph represents a low level graph datastructure, implemented as an
// adjacency list. It is thread-safe, using sync.Map to hold the graph nodes,
// and a Read Write mutex to synchronise access to the 'size' which holds the
// total number of nodes in the graph.
type Graph struct {
	Nodes sync.Map

	mu   sync.RWMutex
	size int
}

// GetSize returns the total number of nodes in the graph.
func (g *Graph) GetSize() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.size
}

func (g *Graph) incrementSize(inc int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.size += inc
}

// Node represents a node in the Graph.  It's Value is kept intentionally
// generic, represented by the empty interface. Each GraphNode also maintains a
// totalWeight field, which is the sum of each of the graphs edges.
// Maintaining this value allows for randomly selecting an edge node,
// proportional to the weight of the edge nodes.
type Node struct {
	mu     sync.RWMutex
	Value  interface{}
	totalWeight int
	edges  []*Edge
}

func (n *Node) getTotalWeight() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.totalWeight
}

func (n *Node) incrementWeight(inc int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.totalWeight += inc
}

// TotalEdges returns the number of edges a node has in the graph
func (n *Node) TotalEdges() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.edges)
}

// NewGraph returns a new Graph datastructure
func NewGraph() *Graph {
	return &Graph{
		Nodes: sync.Map{},
	}
}

// Edge represents an edge between two given nodes
// The Index field holds the destination nodes index
type Edge struct {
	Index string

	mu     sync.RWMutex
	weight int
}

func newEdge(index string) *Edge {
	return &Edge{
		Index:  index,
		weight: 1,
	}
}

func (e *Edge) getWeight() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.weight
}

func (e *Edge) incrementWeight(inc int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.weight += inc
}

// AddNode adds a new node in the graph
func (g *Graph) AddNode(index string, value interface{}) {
	_, ok := g.Nodes.Load(index)
	if !ok {
		node := &Node{Value: value}
		g.Nodes.Store(index, node)
		g.incrementSize(1)
		return
	}
}

// AddEdge takes two indexes to the graph, and add an edge between them.
// Both nodes must exist in the graph prior to creating an edge.
func (g *Graph) AddEdge(n1Index, n2Index string) error {
	// Check if the source node exists
	n1, ok := g.Nodes.Load(n1Index)
	if !ok {
		return errors.New("source node must exist")
	}

	// Check if the destination node exists
	_, ok = g.Nodes.Load(n2Index)
	if !ok {
		return errors.New("destination node must exist")
	}

	// Keep track of the total number of edges adjacent to the GraphNode. We use
	// this value to randomly pick an edge with the same frequency the edge
	// occurs in the corpus.
	graphNode := n1.(*Node)
	graphNode.incrementWeight(1)

	// Check if the edge already exists
	graphNode.mu.RLock()
	defer graphNode.mu.RUnlock()
	for _, edge := range graphNode.edges {
		// if the edge already exists
		if edge.Index == n2Index {
			// Increment the weight (represents frequency)
			edge.incrementWeight(1)
			return nil
		}
	}

	// Otherwise, it doesn't exist, create the new edge
	edge := newEdge(n2Index)

	// Insert the new edge
	graphNode.edges = append(graphNode.edges, edge)
	return nil
}

// RandomEdge selects a random edge adjacent to the node it is called on.
// This has a worst case of O(n).
func (n *Node) RandomEdge() *Edge {
	stoppingPoint := rand.Intn(n.getTotalWeight()) + 1
	count := 0

	n.mu.RLock()
	defer n.mu.RUnlock()
	for _, edge := range n.edges {
		count += edge.getWeight()
		if count >= stoppingPoint {
			return edge
		}
	}

	// We should never reach here
	return nil
}

// StopCase is a function for testing the point at which RandomWalk should stop
// traversing the graph.  It takes a count representing the number of nodes
// processed, and the current node being processed.
type StopCase func(int, *Node) bool

// RandomWalk takes a starting key, and walks through the graph nodes, testing
// each node it visits using the StopCase function. When the StopCase returns
// false, the graph will stop traversing the graph.
func (g *Graph) RandomWalk(startIndex string, fn StopCase) {
	randomWalk(g, startIndex, 0, fn)
}

func randomWalk(g *Graph, start string, count int, fn StopCase) {
	count++
	node, ok := g.Nodes.Load(start)
	if ok {
		graphNode := node.(*Node)
		if ok && fn(count, graphNode) {
			return
		}

		randomWalk(g, graphNode.RandomEdge().Index, count, fn)
	}
}

// Print prints a JSON representation of the graph datastructure to the
// console. Mostly used for debugging purposes.
func (g *Graph) Print() {
	m := map[string]interface{}{}
	g.Nodes.Range(func(key, value interface{}) bool {
		val, ok := value.(*Node)
		if ok {
			// Create an anonymous struct with exported fields for the purpose of printing edges.
			m[fmt.Sprint(key)] = struct {
				Value interface{}
				Edges []*Edge
			}{
				Value: val.Value,
				Edges: val.edges,
			}
		}

		return true
	})

	jsonBs, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonBs))
}
