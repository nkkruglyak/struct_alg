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
    g "custom_graph"
)

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

//custom_io
//1) rerurn list of edges and vertices
//2) define custom InitGraphFromEdgesAndVertices -- greedy
// ReadVerticesFromFile init graph woth edge by line
// if line equals "a b" add graph edge ab
// if line equals "c" add graph edge c
// N.B.!there are not validation of graf like as "a b\na"
func ReadVerticesFromFile(path string) ([][]string,[]string, error) {
    file, err := os.Open(path)
    edges := [][]string{}
    vertices := []string{}
    if err != nil {
        Error.Println("Can't open file")
        return edges, vertices, err
    }

    edge := make([]string, 2)
    defer file.Close()
    sCaner := bufio.NewScanner(file)
    for sCaner.Scan() {
        edge = strings.Split(sCaner.Text(), " ")
        switch len(edge) {
        case 2:
            edges = append(edges, edge)
        case 1:
            vertices = append(vertices, edge[0])
        }
        Info.Printf("Added edge %v", edge)
    }
    return edges, vertices, nil
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

    graph,err := ReadVerticesFromFile(*input_file)
    if err != nil {
        Error.Fatal("Can't init graph")
    }
    graph.PaintGraph()
    data := graph.DumpPaintedGraphToList()
    err = DumpDataToFile(data, *output_file)
    if err != nil {
        Error.Fatal("Can't dump graph")
    }
}
