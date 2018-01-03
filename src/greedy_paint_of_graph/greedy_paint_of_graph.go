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
    PaintedVertices   []Vertex
    NoPaintedVertices []Vertex
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

func (g Graph) GetVertex (id string) (v Vertex){

    for _, v := range  g.PaintedVertices {
        if v.Id == id {
            return v
        }
    }

    for _, v := range  g.NoPaintedVertices {
        if v.Id == id {
            return v
        }
    }

    return  Vertex{}
}

func (graph *Graph) PaintGraph () () {
    Info.Println("Begin Paint Graph")
    color := 1
    var vertex Vertex

    for len(graph.NoPaintedVertices) > 0 {
        for ind := 0; ind < len(graph.NoPaintedVertices); ind++ {
            vertex = graph.NoPaintedVertices[ind]

            neighbour_colors := vertex.GetColorsOfNeighbours(graph)

            if neighbour_colors[color] {
                Info.Printf(
                    "Can't paint %s vertex to color %d",
                    vertex.Id,
                    color,
                )
                continue
            }

            vertex.Paint(color)
            Info.Printf(
                "Paint %s vertex to color %d",
                vertex.Id,
                color,
            )

            graph.NoPaintedVertices = cut_el_from_list(graph.NoPaintedVertices, ind)

            graph.PaintedVertices = put_el_to_list(graph.PaintedVertices, vertex, -1)
        }
        color += 1
    }
}

func AddFirstVertexOfEdgeToMap(graph_as_map map[string][]string, edge []string) {
    neighbours := graph_as_map[edge[0]]
    if len(neighbours) > 0 {
        graph_as_map[edge[0]] = append(graph_as_map[edge[0]], edge[1])
    } else {
        graph_as_map[edge[0]] = []string{edge[1]}
    }
}

func AddEdgeToMap(graph_as_map map[string][]string, edge []string) {
    AddFirstVertexOfEdgeToMap(graph_as_map, edge)
    AddFirstVertexOfEdgeToMap(graph_as_map, []string{edge[1],edge[0]})
}

// сравнить по памяти времени с методом
//function (g Graph) GetVertex (id int) {
//    return g.Vertex[id]
//}
//
// type Graph struct {
//    Vertices []Vertex
// }

// readGraphFromFile reads a whole file into memory
// and returns a slice of its lines.
func readGraphFromFile(path string) (Graph, error) {
    file, err := os.Open(path)
    if err != nil {
        Error.Println("Can't open file")
        return Graph{}, err
    }
    defer file.Close()

    edge := make([]string,2)
    vertices := []Vertex{}
    sCaner := bufio.NewScanner(file)
    graph_as_map := make(map[string][]string)

    for sCaner.Scan() {
        edge = strings.Split(sCaner.Text(), " ")
        AddEdgeToMap(graph_as_map, edge)
        Info.Printf("Added edge %v",edge)
    }

    for vertex, neighbours := range graph_as_map {
        vertex := Vertex{
            Id:vertex,
            NeighboursIds:neighbours,
            }
        vertices = append(vertices, vertex)
        Info.Printf(
            "Added vertex %s with neighbours %v",
            vertex.Id,
            vertex.NeighboursIds,
            )
    }

    graph := Graph{NoPaintedVertices:vertices}

    Info.Printf(
        "Init graph with %d painted and %d no painted vertices.",
        len(graph.PaintedVertices),
        len(graph.NoPaintedVertices),
    )
    return graph, sCaner.Err()
}

//хорошое место для интерфейсов
//метод конечно долэен работать для всех списков
func cut_el_from_list (list []Vertex, ind int) ([]Vertex) {
    n := len(list)
    new_list := make([]Vertex, n - 1)
    at := copy(new_list, list[:ind])
    copy(new_list[at:],list[ind+1:])
    return new_list
}

func put_el_to_list(list []Vertex, v Vertex, ind int) ([]Vertex) {
    switch ind {
    case -1:
        list = append(list, v)
    default:
        new_list := make([]Vertex, len(list) + 1)
        at := copy(new_list, list[ind:])
        new_list[at] = v
        copy(new_list[at+1:], list[ind+1:])
    }
    return list
}

func main () (){
    Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

    filePtr := flag.String(
        "file",
        "testdata_1",
        "path to file with graph description",
    )

    flag.Parse()

    graph,err := readGraphFromFile(*filePtr)
    if err != nil {
        Error.Fatal("Can't init graph")
    }
    graph.PaintGraph()
}
