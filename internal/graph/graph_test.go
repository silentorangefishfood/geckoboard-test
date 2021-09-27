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
		node := g.Nodes[index]
		if node == nil {
			t.Errorf("Expected graph to contain node: %s", index)
		}

		value := node.Value.(TestVal).Int
		if value != 1 {
			t.Errorf("Expected node to contain value: %d", value)
		}

		if node.TotalWeight != 0 {
			t.Errorf("Expected weight of node to equal 0: %d", node.TotalWeight)
		}
	})

	t.Run("AddSecondNode", func(t *testing.T) {
		index := "test2"
		g.AddNode(index, TestVal{ Int: 2 })
		node := g.Nodes[index]
		if node == nil {
			t.Errorf("Expected graph to contain node: %s", index)
		}

		value := node.Value.(TestVal).Int
		if value != 2 {
			t.Errorf("Expected node to contain value: %d", value)
		}

		if node.TotalWeight != 0 {
			t.Errorf("Expected weight of node to equal 0: %d", node.TotalWeight)
		}
	})

	t.Run("AddEdge", func(t *testing.T) {
		i1 := "test1"
		i2 := "test2"
		g.AddEdge(i1, i2)
		srcNode := g.Nodes[i1]
		if srcNode.TotalWeight != 1 {
			t.Errorf("Expected src node to have weight of 1: %d", srcNode.TotalWeight)
		}

		if len(srcNode.Edges) != 1 {
			t.Errorf("Expected src node to have 1 edge: %d", len(srcNode.Edges))
		}

		edgesTotal := 0
		for _, e := range srcNode.Edges {
			edgesTotal += e.Weight
		}

		if srcNode.TotalWeight != edgesTotal {
			t.Errorf("Expected src node total weight to equal sum of edge weights: %d", edgesTotal)
		}
	})

	t.Run("AddThirdNode", func(t *testing.T) {
		index := "test3"
		g.AddNode(index, TestVal{ Int: 3 })
		node := g.Nodes[index]
		if node == nil {
			t.Errorf("Expected graph to contain node: %s", index)
		}

		value := node.Value.(TestVal).Int
		if value != 3 {
			t.Errorf("Expected node to contain value: %d", value)
		}

		if node.TotalWeight != 0 {
			t.Errorf("Expected weight of node to equal 0: %d", node.TotalWeight)
		}
	})
}
