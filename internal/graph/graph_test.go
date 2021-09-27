package graph

import (
	"os"
	"testing"
)

// TestMain
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// TestGraph
func TestGraph(t *testing.T) {
	g := NewGraph()
	t.Run("AddFirstNode", func(t *testing.T) {
		w1 := "To"
		w2 := "be"
		g.AddNode(w1, w2)
		index := w1 + w2
		node := g.Nodes[index]
		if node == nil {
			t.Errorf("Expected graph to contain node: %s", index)
		}

		if node.Word1 != w1 {
			t.Errorf("Expected graph node to contain word1: %s", w1)
		}

		if node.Word2 != w2 {
			t.Errorf("Expected graph node to contain word2: %s", w1)
		}

		if node.TotalWeight != 0 {
			t.Errorf("Expected weight of node to equal 0: %d", node.TotalWeight)
		}
	})

	t.Run("AddSecondNode", func(t *testing.T) {
		w1 := "be"
		w2 := "or"
		g.AddNode(w1, w2)
		index := w1 + w2
		node := g.Nodes[index]
		if node == nil {
			t.Errorf("Expected graph to contain node: %s", index)
		}

		if node.Word1 != w1 {
			t.Errorf("Expected graph node to contain word1: %s", w1)
		}

		if node.Word2 != w2 {
			t.Errorf("Expected graph node to contain word2: %s", w1)
		}

		if node.TotalWeight != 0 {
			t.Errorf("Expected weight of node to equal 0: %d", node.TotalWeight)
		}
	})

	t.Run("AddEdge", func(t *testing.T) {
		i1 := "Tobe"
		i2 := "beor"
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
		w1 := "or"
		w2 := "not"
		g.AddNode(w1, w2)
		index := w1 + w2
		node := g.Nodes[index]
		if node == nil {
			t.Errorf("Expected graph to contain node: %s", index)
		}

		if node.Word1 != w1 {
			t.Errorf("Expected graph node to contain word1: %s", w1)
		}

		if node.Word2 != w2 {
			t.Errorf("Expected graph node to contain word2: %s", w1)
		}

		if node.TotalWeight != 0 {
			t.Errorf("Expected weight of node to equal 0: %d", node.TotalWeight)
		}
	})
}
