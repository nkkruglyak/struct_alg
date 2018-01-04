package main

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

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

// custom_log
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


//custom_io
//1) rerurn list of edges and vertices
//2) define custom InitGraphFromEdgesAndVertices -- greedy
// InitGraphFromFile init graph woth edge by line
// if line equals "a b" add graph edge ab
// if line equals "c" add graph edge c
// N.B.!there are not validation of graf like as "a b\na"
func InitGraphFromFile(path string) (Graph, error) {
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

//greedy
func DumpPaintedGraphToList(graph Graph) ([]string) {
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

//custom_io
func DumpDataToFile (dump_data []string, path string) (error) {
    file, err := os.Open(path)
    
    if err != nil {
        Error.Println("Can't open file")
        file, err = os.Create(path)
        if err != nil {
            Error.Println("Can't create file")
        }
        return err
    }
    
    defer file.Close()
    
    writer := bufio.NewWriter(file)
    
    for _, data := range dump_data {
        writer.WriteString(data)
    }
    writer.Flush()
    return nil
}

func main () (){
    Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

    input_file := flag.String(
        "input",
        "testdata_1",
        "path to file with graph description",
    )

    output_file := flag.String(
        "output",
        "result_1",
        "path to file with painted graph",
    )

    flag.Parse()

    graph,err := InitGraphFromFile(*input_file)
    if err != nil {
        Error.Fatal("Can't init graph")
    }
    graph.PaintGraph()
    data := DumpPaintedGraphToList(graph)
    err = DumpDataToFile(data, *output_file)
    if err != nil {
        Error.Fatal("Can't dump graph")
    }
}
