// ileaf is a basic command line interface
// to run and debug specific parts of the
// application
package main

import (
	"context"
	"encoding/json"
	"io"
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

	// #nosec allow any file for development
	file, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer mustClose(file)

	info, err := file.Stat()
	if err != nil {
		panic(err)
	}

	ig := ingestion.NewIngestion()
	registerExtractors(ig)

	pages, err := ig.Extract(context.Background(), file, info.Size())
	if err != nil {
		panic(err)
	}

	// #nosec allow any file for development
	if err := os.MkdirAll(filepath.Dir(output), os.ModePerm); err != nil {
		panic(err)
	}

	// #nosec allow any file for development
	outFile, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer mustClose(outFile)

	if err := json.NewEncoder(outFile).Encode(pages); err != nil {
		panic(err)
	}
}

func registerExtractors(ig *ingestion.Ingestion) {
	pdf := extract.NewPdfExtractor()
	ig.RegisterExtractor(pdf, "application/pdf")

	docx := extract.NewDocxExtractor()
	ig.RegisterExtractor(docx, "application/vnd.openxmlformats-officedocument.wordprocessingml.document")

	plain := extract.NewPlainExtractor()
	ig.RegisterExtractor(plain,
		"text/plain", "text/markdown", "text/x-markdown", "text/csv",
		"text/tab-separated-values", "application/json", "application/yaml",
		"application/x-yaml", "text/yaml", "application/toml",
	)
}

func mustClose(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}
