// Package cspf implements the Costrained Shortest Path First
// algorithm to find the shortest paths in a graph that satisfy
// a set of generic conditions.
// Every edge that connects two vertices of the graph can be labeled
// with generic key/value pairs.
// The conditions are parsed and evaluated using
// github.com/PaesslerAG/gval package. Thus, any expression and
// value type currently supported by this package can be used
// to state the constraints.
package cspf

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/PaesslerAG/gval"
)

var (
	// ErrDuplicateTagKey is returned by AddEdge method
	// when one Tag's key was specified more than once.
	ErrDuplicateTagKey = errors.New("DuplicateTagKey")
	// ErrNilGraph is returned whenever one method
	// was called on a nil cspf.Graph object
	ErrNilGraph = errors.New("NilGraph")
)

const infinity = uint64(math.MaxUint64)

// Tag contains a generic key/value pair that can
// be attached to cspf.Edge.
// Key must be unique within a single cspf.Edge.
type Tag struct {
	// Key is the unique key string
	Key string
	// Value is a generic value
	Value interface{}
}

// Vertex represents a vertex of the graph.
type Vertex struct {
	// ID identifies uniquely a vertex inside
	// the graph.
	ID string
	//tags map[string]interface{} //TODO: implement constraints eval for vertex
}

// Edge represents a directed edge that connects
// one vertex to another.
// Edge must have a non-negative cost.
// Edge can have a set of generic tags. Each tag must have
// an unique key string. CSPF conditions apply on these tags.
type Edge struct {
	// Source vertex of this edge.
	From Vertex
	// Destination vertex of this edge.
	To Vertex
	// Numeric cost of this edge.
	Cost uint64
	// Set of generic key/value tags.
	// Tag key is an unique string, whereas
	// the value can be of any type.
	Tags map[string]interface{}
}

// Graph represents a directed graph.
type Graph struct {
	// Set of vertices of this graph
	// with associated list of edges originating
	// from every vertex.
	VertexSet map[Vertex][]Edge
	eval      gval.Evaluable
}

func (g *Graph) initGraph() {
	if g.VertexSet == nil {
		g.VertexSet = make(map[Vertex][]Edge)
	}
}

// AddEdge adds a new edge between two vertices with an associated
// cost and possibly a set of tags.
// If the two vertices do not exist, AddEdge adds them
// to the graph automatically.
func (g *Graph) AddEdge(from, to Vertex, cost uint64, tags ...Tag) error {
	edge := Edge{
		From: from,
		To:   to,
		Cost: cost,
	}
	if len(tags) != 0 {
		edge.Tags = make(map[string]interface{})
		for _, tag := range tags {
			if _, ok := edge.Tags[tag.Key]; ok {
				return fmt.Errorf("%w: %s", ErrDuplicateTagKey, tag.Key)
			}
			edge.Tags[tag.Key] = tag.Value
		}
	}
	g.addEdge(edge)
	return nil
}

func (g *Graph) addEdge(e Edge) {
	g.AddNode(e.From)
	g.AddNode(e.To)

	edges := g.VertexSet[e.From]
	edges = append(edges, e)
	g.VertexSet[e.From] = edges
}

// AddNode adds a new vertex to the graph with no edges.
// It is preferable to use AddEdge, given that it
// adds the vertex automatically if missing.
func (g *Graph) AddNode(v Vertex) {
	g.initGraph()

	_, found := g.VertexSet[v]
	if !found {
		g.VertexSet[v] = []Edge{}
	}
}

// SPF runs the Dijkstra algorithm to build a result
// graph only containing the shortest paths from one
// vertex to another.
// All shortest paths with equal cost are part of the
// result graph.
func (g *Graph) SPF(from, to Vertex) (*Graph, error) {
	if g == nil {
		return nil, ErrNilGraph
	}
	unvisitedSet := make(map[Vertex]bool)
	distSet := make(map[Vertex]uint64)
	prevSet := make(map[Vertex][]Edge)

	for v := range g.VertexSet {
		unvisitedSet[v] = true
		distSet[v] = infinity
		prevSet[v] = []Edge{}
	}
	distSet[from] = 0

	for len(unvisitedSet) > 0 {
		setSize := len(unvisitedSet)
		closestVertex := getSmallestDistanceVertex(unvisitedSet, distSet)
		delete(unvisitedSet, closestVertex)
		if setSize == len(unvisitedSet) {
			//No progress on the visited set, some vertex
			//of this graph might not be reachable by <from>
			break
		}

		for _, edge := range g.VertexSet[closestVertex] {
			if stillUnvisited := unvisitedSet[edge.To]; stillUnvisited {
				satisfied, err := g.edgeSatisfiesConstranints(edge)
				if err != nil {
					return nil, err
				}
				if satisfied {
					distFromNeighbor := distSet[closestVertex] + edge.Cost
					if distFromNeighbor <= distSet[edge.To] {
						distSet[edge.To] = distFromNeighbor
						edges := prevSet[edge.To]
						edges = append(edges, edge)
						prevSet[edge.To] = edges
					}
				}
			}
		}
	}

	SPF := Graph{}
	if _, ok := prevSet[to]; ok || to == from {
		for _, edges := range prevSet {
			for _, edge := range edges {
				SPF.addEdge(edge)
			}
		}
	}

	return &SPF, nil
}

func getSmallestDistanceVertex(unvisited map[Vertex]bool, distSet map[Vertex]uint64) Vertex {
	smallestDist := infinity
	closestVertex := Vertex{}
	for v := range unvisited {
		dist := distSet[v]
		if dist < smallestDist {
			smallestDist = dist
			closestVertex = v
		}
	}
	return closestVertex
}

// CSPF runs the Constrained Shortest Path First algorithm
// to perform find the shortest paths which edges all satisfy
// the specified expression.
// An edge cannot be part of the resulting graph if its tags
// do not satisfy the specified expression.
// The tag key/value pairs and the generic expression are
// internally evaluated through github.com/PaesslerAG/gval
// package.
func (g *Graph) CSPF(from, to Vertex, exp string) (*Graph, error) {
	if g == nil {
		return nil, ErrNilGraph
	}
	eval, err := gval.Full().NewEvaluable(exp)
	if err != nil {
		return nil, err
	}
	g.eval = eval
	return g.SPF(from, to)
}

func (g *Graph) edgeSatisfiesConstranints(e Edge) (bool, error) {
	if g.eval == nil {
		return true, nil
	}

	match, err := g.eval.EvalBool(context.Background(), e.Tags)
	if err != nil {
		return false, err
	}
	return match, nil
}

// Paths lists all the possible paths of the graph that
// connect from one vertex to the other.
// Paths are listed through Depth-first search algorithm.
func (g *Graph) Paths(from, to Vertex) (paths [][]Edge) {
	if g == nil {
		return
	}
	//Explore the graph using Depth First Search
	//starting from the <from> object and listing
	//all the paths that reach <to>

	visited := make(map[Vertex]bool)
	path := []Edge{}

	var dfs func(v Vertex, edge *Edge)
	dfs = func(v Vertex, edge *Edge) {
		visited[v] = true
		if edge != nil {
			path = append(path, *edge)
		}

		if v == to {
			paths = append(paths, path)
		} else {
			for _, edge := range g.VertexSet[v] {
				if !visited[edge.To] {
					dfs(edge.To, &edge)
				}
			}
		}

		if len(path) > 0 {
			path = path[1:]
		}
		visited[v] = false
	}

	dfs(from, nil)

	return
}
