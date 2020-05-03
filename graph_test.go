package cspf_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bigmikes/cspf"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSPF(t *testing.T) {
	a := cspf.Vertex{ID: "a"}
	b := cspf.Vertex{ID: "b"}
	c := cspf.Vertex{ID: "c"}
	d := cspf.Vertex{ID: "d"}

	graph := cspf.Graph{}

	Convey("Populate thegraph with no error", t, func() {
		err := graph.AddEdge(a, b, 1)
		So(err, ShouldBeNil)
		err = graph.AddEdge(a, c, 1)
		So(err, ShouldBeNil)
		err = graph.AddEdge(b, d, 1)
		So(err, ShouldBeNil)
		err = graph.AddEdge(c, d, 1)
		So(err, ShouldBeNil)
	})

	Convey("Run the SPF algorithm", t, func() {
		spfGraph, err := graph.SPF(a, d)
		So(err, ShouldBeNil)
		So(spfGraph, ShouldNotBeNil)
		paths := spfGraph.Paths(a, d)
		So(paths, ShouldNotBeNil)
		//Check paths lengths
		So(len(paths), ShouldEqual, 2)
		So(len(paths[0]), ShouldEqual, 2)
		So(len(paths[1]), ShouldEqual, 2)
		//Check if both paths are correct
		//Both paths must start from vertex A
		So(paths[0][0].From, ShouldResemble, a)
		So(paths[1][0].From, ShouldResemble, a)
		//Then, depending on the next vertex, we follow
		//two different paths with equal cost:
		//1) A -> B -> D
		//2) A -> C -> D
		if paths[0][0].To.ID == "b" {
			So(paths[0][0].To, ShouldResemble, b)
			So(paths[0][1].From, ShouldResemble, b)
			So(paths[1][0].To, ShouldResemble, c)
			So(paths[1][1].From, ShouldResemble, c)

		} else {
			So(paths[0][0].To, ShouldResemble, c)
			So(paths[0][1].From, ShouldResemble, c)
			So(paths[1][0].To, ShouldResemble, b)
			So(paths[1][1].From, ShouldResemble, b)
		}
		//Destination vertex must be D for both paths
		So(paths[0][1].To, ShouldResemble, d)
		So(paths[1][1].To, ShouldResemble, d)
	})
}

func TestSPFWithCycle(t *testing.T) {
	a := cspf.Vertex{ID: "a"}
	b := cspf.Vertex{ID: "b"}
	c := cspf.Vertex{ID: "c"}
	d := cspf.Vertex{ID: "d"}
	e := cspf.Vertex{ID: "e"}

	graph := cspf.Graph{}

	Convey("Populate thegraph with no error", t, func() {
		err := graph.AddEdge(a, b, 1)
		So(err, ShouldBeNil)
		err = graph.AddEdge(b, c, 1)
		So(err, ShouldBeNil)
		err = graph.AddEdge(c, d, 1)
		So(err, ShouldBeNil)
		err = graph.AddEdge(d, a, 1)
		So(err, ShouldBeNil)
		err = graph.AddEdge(d, e, 1)
		So(err, ShouldBeNil)
	})

	Convey("Run the SPF algorithm", t, func() {
		spfGraph, err := graph.SPF(a, e)
		So(err, ShouldBeNil)
		So(spfGraph, ShouldNotBeNil)
		paths := spfGraph.Paths(a, e)
		So(paths, ShouldNotBeNil)
		//Check paths lengths
		So(len(paths), ShouldEqual, 1)
		So(len(paths[0]), ShouldEqual, 4)
		//Check that the only path is actually correct:
		//1) A -> B -> C -> D -> E
		So(paths[0][0].From, ShouldResemble, a)
		So(paths[0][0].To, ShouldResemble, b)
		So(paths[0][1].From, ShouldResemble, b)
		So(paths[0][1].To, ShouldResemble, c)
		So(paths[0][2].From, ShouldResemble, c)
		So(paths[0][2].To, ShouldResemble, d)
		So(paths[0][3].From, ShouldResemble, d)
		So(paths[0][3].To, ShouldResemble, e)
	})
}

func TestCSPF(t *testing.T) {
	tagBlue := cspf.Tag{
		Key:   "link",
		Value: "blue",
	}
	tagRed := cspf.Tag{
		Key:   "link",
		Value: "red",
	}
	tagRedBlue := cspf.Tag{
		Key:   "link",
		Value: "redblue",
	}

	a := cspf.Vertex{ID: "a"}
	b := cspf.Vertex{ID: "b"}
	c := cspf.Vertex{ID: "c"}
	d := cspf.Vertex{ID: "d"}
	e := cspf.Vertex{ID: "e"}

	graph := cspf.Graph{}

	Convey("Populate thegraph with no error", t, func() {
		err := graph.AddEdge(a, b, 1, tagBlue)
		So(err, ShouldBeNil)
		err = graph.AddEdge(a, c, 1, tagRed)
		So(err, ShouldBeNil)
		err = graph.AddEdge(b, d, 1, tagBlue)
		So(err, ShouldBeNil)
		err = graph.AddEdge(c, d, 1, tagRed)
		So(err, ShouldBeNil)
		err = graph.AddEdge(d, e, 1, tagRedBlue)
		So(err, ShouldBeNil)
	})

	Convey("Run the CSPF algorithm", t, func() {
		spfGraph, err := graph.CSPF(a, e, `link == "blue" || link == "redblue"`)
		So(err, ShouldBeNil)
		So(spfGraph, ShouldNotBeNil)
		paths := spfGraph.Paths(a, e)
		So(paths, ShouldNotBeNil)
		//Check paths lengths
		So(len(paths), ShouldEqual, 1)
		So(len(paths[0]), ShouldEqual, 3)
		//Check if both paths are correct
		So(paths[0][0].From, ShouldResemble, a)
		So(paths[0][0].To, ShouldResemble, b)
		So(paths[0][1].From, ShouldResemble, b)
		So(paths[0][1].To, ShouldResemble, d)
		So(paths[0][2].From, ShouldResemble, d)
		So(paths[0][2].To, ShouldResemble, e)
	})
}

func TestCallsOnNilGraph(t *testing.T) {
	a := cspf.Vertex{ID: "a"}
	b := cspf.Vertex{ID: "b"}
	var graph *cspf.Graph = nil

	Convey("Call all APIs on a nil graph and check if they panic", t, func() {
		So(func() {
			graph.SPF(a, b)
		}, ShouldNotPanic)
		So(func() {
			graph.CSPF(a, b, `link != "blue"`)
		}, ShouldNotPanic)
		So(func() {
			graph.Paths(a, b)
		}, ShouldNotPanic)
	})

	Convey("Call all APIs on a nil graph and check the errors", t, func() {
		g, err := graph.SPF(a, b)
		So(g, ShouldBeNil)
		So(err, ShouldBeError, cspf.ErrNilGraph)
		So(errors.Is(err, cspf.ErrNilGraph), ShouldEqual, true)
		g, err = graph.CSPF(a, b, `link != "blue"`)
		So(g, ShouldBeNil)
		So(err, ShouldBeError, cspf.ErrNilGraph)
		So(errors.Is(err, cspf.ErrNilGraph), ShouldEqual, true)
		p := graph.Paths(a, b)
		So(p, ShouldBeNil)
		So(len(p), ShouldEqual, 0)
	})
}

func TestDuplicateKeyError(t *testing.T) {
	tagBlue := cspf.Tag{
		Key:   "link",
		Value: "blue",
	}
	tagRed := cspf.Tag{
		Key:   "link",
		Value: "red",
	}

	a := cspf.Vertex{ID: "a"}
	b := cspf.Vertex{ID: "b"}

	graph := cspf.Graph{}

	Convey("Assign two equal tag keys and read the error", t, func() {
		err := graph.AddEdge(a, b, 1, tagBlue, tagRed)
		So(err, ShouldBeError, fmt.Errorf("%w: %s", cspf.ErrDuplicateTagKey, "link"))
		So(errors.Is(err, cspf.ErrDuplicateTagKey), ShouldEqual, true)
	})
}

func TestCSPFInvalidCondition(t *testing.T) {
	tagBlue := cspf.Tag{
		Key:   "link",
		Value: "blue",
	}

	a := cspf.Vertex{ID: "a"}
	b := cspf.Vertex{ID: "b"}

	graph := cspf.Graph{}

	Convey("Populate the graph with no error", t, func() {
		err := graph.AddEdge(a, b, 1, tagBlue)
		So(err, ShouldBeNil)
	})

	Convey("Run the CSPF algorithm with invalid condition (or instead of ||)", t, func() {
		spfGraph, err := graph.CSPF(a, b, `link == "blue" or link == "redblue"`)
		So(err, ShouldNotBeNil)
		So(spfGraph, ShouldBeNil)
		paths := spfGraph.Paths(a, b)
		So(paths, ShouldBeNil)
	})
}

func generateFullyConnectedGraph(nVertex int, tag bool) (*cspf.Graph, []cspf.Vertex) {
	graph := cspf.Graph{}
	vertices := make([]cspf.Vertex, 0, nVertex)
	for i := 0; i < nVertex; i++ {
		v := cspf.Vertex{ID: fmt.Sprintf("%d", i)}
		vertices = append(vertices, v)
	}
	for i := 0; i < nVertex; i++ {
		for j := 0; j < nVertex; j++ {
			if tag {
				graph.AddEdge(vertices[i], vertices[j], 1, cspf.Tag{
					Key:   "key",
					Value: "value",
				})
			} else {
				graph.AddEdge(vertices[i], vertices[j], 1)
			}
		}
	}
	return &graph, vertices
}

func BenchmarkSPF(b *testing.B) {
	graph, vertices := generateFullyConnectedGraph(100, false)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		spfGraph, err := graph.SPF(vertices[0], vertices[len(vertices)-1])
		if err != nil {
			b.Fatal(err)
		}
		_ = spfGraph
	}
}

func BenchmarkPaths(b *testing.B) {
	graph, vertices := generateFullyConnectedGraph(100, false)
	spfGraph, err := graph.SPF(vertices[0], vertices[len(vertices)-1])
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		paths := spfGraph.Paths(vertices[0], vertices[len(vertices)-1])
		_ = paths
	}
}

func BenchmarkCSPF(b *testing.B) {
	graph, vertices := generateFullyConnectedGraph(100, true)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		spfGraph, err := graph.CSPF(vertices[0], vertices[len(vertices)-1], `key == "value"`)
		if err != nil {
			b.Fatal(err)
		}
		_ = spfGraph
	}
}
