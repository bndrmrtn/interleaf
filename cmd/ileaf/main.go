// ileaf is a basic command line interface
// to run and debug specific parts of the
// application
package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/bndrmrtn/interleaf/internal/extract"
	"github.com/bndrmrtn/interleaf/internal/ingestion"
)

func main() {
	if len(os.Args) != 3 {
		panic("Usage: program input.pdf output.json")
	}

	input := os.Args[1]
	output := os.Args[2]

	file, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		panic(err)
	}

	ig := ingestion.NewIngestion()
	ex := extract.NewPdfExtractor()
	ig.RegisterExtractor(ex, "application/pdf")
	pages, err := ig.Extract(file, info.Size())
	if err != nil {
		panic(err)
	}

	os.MkdirAll(filepath.Dir(output), os.ModePerm)

	outFile, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	json.NewEncoder(outFile).Encode(pages)
}
