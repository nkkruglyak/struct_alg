package custom_graph

import (
	"os"
	"bufio"
	"strings"
	"flag"
	"log"
	"io"
	"io/ioutil"
	"fmt"
)

type Vertex struct {
	Id string
	NeighboursIds []string
	IsPainted bool
	Color int
}

type Graph struct {
	PaintedVertices   map[string]Vertex
	NoPaintedVertices map[string]Vertex
}


// custom_graph
func (v *Vertex) Paint (color int) {
	v.Color = color
	v.IsPainted = true
}

// custom_graph
func (v *Vertex) GetColorsOfNeighbours (g *Graph)(map[int]bool) {
	colors := make(map[int]bool)
	for _,id := range v.NeighboursIds {
		neighbour := g.GetVertex(id)
		colors[neighbour.Color] = true

	}
	return colors
}

// custom_graph
func (graph Graph) GetVertex (id string) (Vertex){
	//я уверена что это это можно сделать изящнее
	vertex := graph.PaintedVertices[id]
	if vertex.Id == "" {
		vertex = graph.NoPaintedVertices[id]
	}
	return vertex
}

// custom_graph
func (graph *Graph) PaintVertex (vertex Vertex, color int) {
	vertex.Paint(color)
	delete(graph.NoPaintedVertices, vertex.Id)
	graph.PaintedVertices[vertex.Id] = vertex
}

//greedy
func AddFirstVertexOfEdgeToMap(vertex_neighbours_map map[string][]string, edge []string) {
	neighbours := vertex_neighbours_map[edge[0]]
	if len(neighbours) > 0 {
		vertex_neighbours_map[edge[0]] = append(vertex_neighbours_map[edge[0]], edge[1])
	} else {
		vertex_neighbours_map[edge[0]] = []string{edge[1]}
	}
}

//as 2 - greedy
func AddEdgeToMap(vertex_neighbours_map map[string][]string, edge []string) {
	AddFirstVertexOfEdgeToMap(vertex_neighbours_map, edge)
	AddFirstVertexOfEdgeToMap(vertex_neighbours_map, []string{edge[1],edge[0]})
}

func InitGraph(edges [][]string, alone_vertices [] string, Info *log.Logger) (Graph) {
	vertices := []Vertex{}
	vertex_neighbours_map := make(map[string][]string)

	for _, edge := range edges {
		AddEdgeToMap(vertex_neighbours_map, edge)
	}

	for _, v := range alone_vertices {
		vertex_neighbours_map[v] = []string{}
	}

	graph := Graph{
		NoPaintedVertices: make(map[string]Vertex),
		PaintedVertices: make(map[string]Vertex),
	}

	for vertex_id, neighbours_ids := range vertex_neighbours_map {
		vertex := Vertex{
			Id: vertex_id,
			NeighboursIds: neighbours_ids,
		}
		vertices = append(vertices, vertex)
		graph.NoPaintedVertices[vertex.Id] = vertex
		Info.Printf(
			"Added vertex %s with neighbours %v",
			vertex.Id,
			vertex.NeighboursIds,
		)
	}

	Info.Printf(
		"Init graph with %d painted and %d no painted vertices.",
		len(graph.PaintedVertices),
		len(graph.NoPaintedVertices),
	)
	return graph
}


func (graph Graph) DumpPaintedGraphToList() ([]string) {
	dump_data := make([]string, len(graph.PaintedVertices))
	i := 0
	for _, vertex := range graph.PaintedVertices {
		dump_vertex := fmt.Sprintf(
			"vertex %s color %d",
			vertex.Id,
			vertex.Color,
		)
		dump_data[i] = dump_vertex
		i += 1
	}
	return dump_data
}
