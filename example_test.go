package cspf_test

import (
	"fmt"

	"github.com/bigmikes/cspf"
)

func ExampleGraph() {
	// A cspf.Graph needs no initialization.
	var graph cspf.Graph

	// Create two vertices.
	a := cspf.Vertex{ID: "A"}
	b := cspf.Vertex{ID: "B"}

	// Add an directed edge between
	// the two vertices: A -> B.
	graph.AddEdge(a, b, 1)

	// List all the paths from vertex A
	// to vertex B.
	paths := graph.Paths(a, b)

	fmt.Println(paths)
	// Output: [[{{A} {B} 1 map[]}]]
}

func ExampleGraph_SPF() {
	// Create a graph with four vertices.
	// A -> B -> D
	// A -> C -> D
	var graph cspf.Graph
	a := cspf.Vertex{ID: "A"}
	b := cspf.Vertex{ID: "B"}
	c := cspf.Vertex{ID: "C"}
	d := cspf.Vertex{ID: "D"}

	// Add the edges
	graph.AddEdge(a, b, 1)
	graph.AddEdge(a, c, 2)
	graph.AddEdge(b, d, 1)
	graph.AddEdge(c, d, 1)

	// Run the Dijkstra algorithm to
	// derive the graph containing only
	// the shortest paths (equal cost)
	// from A to D.
	spfGraph, _ := graph.SPF(a, d)

	// List the path from vertex A
	// to vertex D.
	paths := spfGraph.Paths(a, d)

	fmt.Println(paths)
	// Output: [[{{A} {B} 1 map[]} {{B} {D} 1 map[]}]]
}

func ExampleGraph_CSPF() {
	// Create a graph with four vertices.
	// A -> B -> D
	// A -> C -> D
	var graph cspf.Graph
	a := cspf.Vertex{ID: "A"}
	b := cspf.Vertex{ID: "B"}
	c := cspf.Vertex{ID: "C"}
	d := cspf.Vertex{ID: "D"}

	// Create a Tag
	tagBlue := cspf.Tag{
		Key:   "link",
		Value: "blue",
	}

	// Add the edges with a label
	// to exclude one specific path
	graph.AddEdge(a, b, 1, tagBlue)
	graph.AddEdge(a, c, 2)
	graph.AddEdge(b, d, 1)
	graph.AddEdge(c, d, 1)

	// Run the CSPF algorithm to
	// derive the graph containing only
	// the shortest paths (equal cost)
	// from A to D.
	cspfGraph, _ := graph.CSPF(a, d, `link != "blue"`)

	// List the path from vertex A
	// to vertex D.
	paths := cspfGraph.Paths(a, d)

	fmt.Println(paths)
	// Output: [[{{A} {C} 2 map[]} {{C} {D} 1 map[]}]]
}
