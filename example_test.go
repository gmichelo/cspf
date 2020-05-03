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
	// Create a graph with five vertices.
	// A -> B -> C -> E
	// A -> D -> E
	var graph cspf.Graph
	a := cspf.Vertex{ID: "A"}
	b := cspf.Vertex{ID: "B"}
	c := cspf.Vertex{ID: "C"}
	d := cspf.Vertex{ID: "D"}
	e := cspf.Vertex{ID: "E"}

	// Create red and blue tag
	tagBlue := cspf.Tag{
		Key:   "link",
		Value: "blue",
	}
	tagRed := cspf.Tag{
		Key:   "link",
		Value: "red",
	}

	// Add the edges with labels
	graph.AddEdge(a, b, 2, tagRed)
	graph.AddEdge(b, c, 2, tagRed)
	graph.AddEdge(c, e, 2, tagRed)
	graph.AddEdge(a, d, 1, tagBlue)
	graph.AddEdge(d, e, 1, tagBlue)

	// Run the CSPF algorithm to
	// derive the graph containing
	// the shortest path from A to E that
	// includes only red edges.
	cspfGraph, _ := graph.CSPF(a, d, `link == "red"`)

	// List the path from vertex A
	// to vertex E.
	paths := cspfGraph.Paths(a, e)

	fmt.Println(paths)
	// Output: [[{{A} {B} 2 map[link:red]} {{B} {C} 2 map[link:red]} {{C} {E} 2 map[link:red]}]]
}
