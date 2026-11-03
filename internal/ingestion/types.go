package ingestion

import (
	"context"
	"io"
)

// Page holds a document's page with it's page number,
// original and markdown formatted content.
type Page struct {
	// Number is the current document page number.
	// Representation may differ for specific document
	// extractors. (DOCX currently shows 1 page only)
	Number          int
	OriginalContent string
	MarkdownContent string
}

// PageExtractor is an inerface for any document
// extractor program.
type PageExtractor interface {
	Extract(
		ctx context.Context,
		file io.ReaderAt,
		size int64,
	) ([]Page, error)
}

// Chunker splits document content into smaller text segments.
type Chunker interface {
	Split(ctx context.Context, content string) ([]string, error)
}

type Chunk struct {
	PageNumber int
	Index      int
	Content    string
}
