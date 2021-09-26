package graph

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestGraph(t *testing.T) {
	g := NewGraph()
	t.Run("AddNode", func(t *testing.T) {
		g.AddNode("To", "be")
	})
}
