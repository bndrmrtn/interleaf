package ingestion

import (
	"context"
	"io"
)

type Page struct {
	Number          int
	OriginalContent string
	MarkdownContent string
}

type PageExtractor interface {
	Extract(
		ctx context.Context,
		file io.ReaderAt,
		size int64,
	) ([]Page, error)
}

type extractor struct {
	contentTypes []string
	pe           PageExtractor
}
