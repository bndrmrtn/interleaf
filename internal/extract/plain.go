package extract

import (
	"context"
	"io"

	"github.com/bndrmrtn/interleaf/internal/ingestion"
)

type PlainExtractor struct{}

func NewPlainExtractor() *PlainExtractor {
	return &PlainExtractor{}
}

// Extract extracts a list of ingestion.Page from
// a PDF document.
func (p *PlainExtractor) Extract(
	_ context.Context,
	file io.ReaderAt,
	size int64,
) ([]ingestion.Page, error) {
	content, err := io.ReadAll(io.NewSectionReader(file, 0, size))
	if err != nil {
		return nil, err
	}

	var (
		contentStr = cleanText(string(content))
		pages      = make([]ingestion.Page, 1)
	)

	pages[0] = ingestion.Page{
		Number:          1,
		OriginalContent: contentStr,
		MarkdownContent: contentStr,
	}

	return pages, nil
}
