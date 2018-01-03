package main

import (
    "os"
    "bufio"
    "strings"
    "flag"
    "log"
    "io"
    "io/ioutil"
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

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func Init(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func (v *Vertex) Paint (color int) {
    v.Color = color
    v.IsPainted = true
}

func (v *Vertex) GetColorsOfNeighbours (g *Graph)(map[int]bool) {
    colors := make(map[int]bool)
    for _,id := range v.NeighboursIds {
          neighbour := g.GetVertex(id)
          colors[neighbour.Color] = true

    }
    return colors
}

func (graph Graph) GetVertex (id string) (Vertex){
    //я уверена что это это можно сделать изящнее
    vertex := graph.PaintedVertices[id]
    if vertex.Id == "" {
        vertex = graph.NoPaintedVertices[id]
    }
    return vertex
}

func (graph *Graph) PaintVertex (vertex Vertex, color int) {
    vertex.Paint(color)
    delete(graph.NoPaintedVertices, vertex.Id)
    graph.PaintedVertices[vertex.Id] = vertex
}

func (graph *Graph) PaintGraph () () {
    Info.Println("Begin Paint Graph")
    color := 1
    for len(graph.NoPaintedVertices) > 0 {
        for _, vertex := range graph.NoPaintedVertices {
            neighbour_colors := vertex.GetColorsOfNeighbours(graph)
            if neighbour_colors[color] {
                continue
            }
            graph.PaintVertex(vertex, color)
            Info.Printf(
                "Paint %s vertex to color %d",
                vertex.Id,
                color,
            )
        }
        color += 1
    }
}

func AddFirstVertexOfEdgeToMap(vertex_neighbours_map map[string][]string, edge []string) {
    neighbours := vertex_neighbours_map[edge[0]]
    if len(neighbours) > 0 {
        vertex_neighbours_map[edge[0]] = append(vertex_neighbours_map[edge[0]], edge[1])
    } else {
        vertex_neighbours_map[edge[0]] = []string{edge[1]}
    }
}

func AddEdgeToMap(vertex_neighbours_map map[string][]string, edge []string) {
    AddFirstVertexOfEdgeToMap(vertex_neighbours_map, edge)
    AddFirstVertexOfEdgeToMap(vertex_neighbours_map, []string{edge[1],edge[0]})
}

// initGraphFromFile init graph woth edge by line
// if line equals "a b" add graph edge ab
// if line equals "c" add graph edge c
// N.B.!there are not validation of graf like as "a b\na"
func initGraphFromFile(path string) (Graph, error) {
    file, err := os.Open(path)
    if err != nil {
        Error.Println("Can't open file")
        return Graph{}, err
    }
    defer file.Close()

    edge := make([]string,2)
    vertices := []Vertex{}
    sCaner := bufio.NewScanner(file)
    vertex_neighbours_map := make(map[string][]string)

    for sCaner.Scan() {
        edge = strings.Split(sCaner.Text(), " ")
        switch len(edge) {
        case 2:
            AddEdgeToMap(vertex_neighbours_map, edge)
        case 1:
            vertex_neighbours_map[edge[0]] = []string{}
        }
        Info.Printf("Added edge %v",edge)
    }

    graph := Graph{
        NoPaintedVertices: make(map[string]Vertex),
        PaintedVertices: make(map[string]Vertex),
    }

    for vertex_id, neighbours_ids := range vertex_neighbours_map {
        vertex := Vertex{
            Id:vertex_id,
            NeighboursIds:neighbours_ids,
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
    return graph, sCaner.Err()
}

func main () (){
    Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

    filePtr := flag.String(
        "file",
        "testdata_1",
        "path to file with graph description",
    )

    flag.Parse()

    graph,err := initGraphFromFile(*filePtr)
    if err != nil {
        Error.Fatal("Can't init graph")
    }
    graph.PaintGraph()
}
