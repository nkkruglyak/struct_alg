package custom_io

import (
	"os"
	"bufio"
	"strings"
	g "custom_graph"
	"flag"
	"log"
	"io"
	"io/ioutil"
	"fmt"
)

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