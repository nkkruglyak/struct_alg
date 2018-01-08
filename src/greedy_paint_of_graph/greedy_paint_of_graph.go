package main

import (
    "os"
    "flag"
    "log"
    "io"
    "io/ioutil"
    g "custom_graph"
    ci "custom_io"
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


func PaintGraph (graph *g.Graph) () {
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

    edges, vertices, err := ci.ReadVerticesFromFile(*input_file, Error, Info)

    if err != nil {
        Error.Fatal("Can't init graph")
    }

    graph := g.InitGraph(edges, vertices, Info)

    PaintGraph(graph)
    data := graph.DumpPaintedGraphToList()
    Info.Println(data)
    err = ci.DumpDataToFile(data, *output_file, Info)
    if err != nil {
        Error.Fatal("Can't dump graph")
    }
}
