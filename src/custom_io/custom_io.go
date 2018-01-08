package custom_io

import (
	"os"
	"bufio"
	"strings"
	"log"
	//"fmt"
)

//custom_io
func DumpDataToFile (dump_data []string, path string, Error *log.Logger) (error) {
	//file, err := os.Open(path)
	//
	//if err != nil {
	//	Error.Println("Can't open file")
	//	file, err = os.Create(path)
	//	if err != nil {
	//		Error.Println("Can't create file")
	//	}
	//	return err
	//}

	file, err := os.Create(path)
	if err != nil {
		Error.Println("Can't create file")
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, data := range dump_data {
		writer.WriteString(data)
	}
	writer.Flush()
	return nil
}


// ReadVerticesFromFile init graph: edge by line
// if line equals "a b" add graph edge ab
// if line equals "c" add graph edge c
// N.B.!there are not validation of graf like as "a b\na"
func ReadVerticesFromFile(path string, Error, Info *log.Logger) ([][]string,[]string, error) {
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
