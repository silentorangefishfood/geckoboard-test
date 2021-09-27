package graph

import (
	"os"
	"testing"
)

// TestMain
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type TestVal struct {
	Int int
}

// TestGraph
func TestGraph(t *testing.T) {
	g := NewGraph()
	t.Run("AddFirstNode", func(t *testing.T) {
		index := "test1"
		g.AddNode(index, TestVal{ Int: 1 })
		node, ok := g.Nodes.Load(index)
		if !ok {
			t.Errorf("Failed to load node")
		}

		if node == nil {
			t.Errorf("Expected graph to contain node: %s", index)
		}

		graphNode := node.(*GraphNode)
		value := graphNode.Value.(TestVal).Int
		if value != 1 {
			t.Errorf("Expected graphNode to contain value: %d", value)
		}

		if graphNode.TotalWeight != 0 {
			t.Errorf("Expected weight of graphNode to equal 0: %d", graphNode.TotalWeight)
		}
	})

	t.Run("AddSecondNode", func(t *testing.T) {
		index := "test2"
		g.AddNode(index, TestVal{ Int: 2 })
		node, ok := g.Nodes.Load(index)
		if !ok {
			t.Errorf("Failed to load node")
		}

		if node == nil {
			t.Errorf("Expected graph to contain node: %s", index)
		}

		graphNode := node.(*GraphNode)
		value := graphNode.Value.(TestVal).Int
		if value != 2 {
			t.Errorf("Expected graphNode to contain value: %d", value)
		}

		if graphNode.TotalWeight != 0 {
			t.Errorf("Expected weight of graphNode to equal 0: %d", graphNode.TotalWeight)
		}
	})

	t.Run("AddEdge", func(t *testing.T) {
		i1 := "test1"
		i2 := "test2"
		g.AddEdge(i1, i2)
		srcNode, ok := g.Nodes.Load(i1)
		if !ok {
			t.Errorf("Failed to load node")
		}

		graphNode := srcNode.(*GraphNode)
		if graphNode.TotalWeight != 1 {
			t.Errorf("Expected src node to have weight of 1: %d", graphNode.TotalWeight)
		}

		if len(graphNode.Edges) != 1 {
			t.Errorf("Expected src node to have 1 edge: %d", len(graphNode.Edges))
		}

		edgesTotal := 0
		for _, e := range graphNode.Edges {
			edgesTotal += e.Weight
		}

		if graphNode.TotalWeight != edgesTotal {
			t.Errorf("Expected src node total weight to equal sum of edge weights: %d", edgesTotal)
		}
	})

	t.Run("AddThirdNode", func(t *testing.T) {
		index := "test3"
		g.AddNode(index, TestVal{ Int: 3 })
		node, ok := g.Nodes.Load(index)
		if !ok {
			t.Errorf("Failed to load node")
		}

		if node == nil {
			t.Errorf("Expected graph to contain node: %s", index)
		}

		graphNode := node.(*GraphNode)
		value := graphNode.Value.(TestVal).Int
		if value != 3 {
			t.Errorf("Expected graphNode to contain value: %d", value)
		}

		if graphNode.TotalWeight != 0 {
			t.Errorf("Expected weight of graphNode to equal 0: %d", graphNode.TotalWeight)
		}
	})
}
