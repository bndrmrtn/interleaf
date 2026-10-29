package ingestion

import (
	"context"
	"fmt"
	"io"
	"mime"

	"github.com/gabriel-vasile/mimetype"
)

type Ingestion struct {
	extractors map[string]PageExtractor
}

func NewIngestion() *Ingestion {
	return &Ingestion{
		extractors: make(map[string]PageExtractor),
	}
}

func (ig *Ingestion) RegisterExtractor(
	pe PageExtractor,
	contentType string,
	contentTypes ...string,
) {
	ig.extractors[contentType] = pe

	for _, contentType := range contentTypes {
		ig.extractors[contentType] = pe
	}
}

func (ig *Ingestion) Extract(
	ctx context.Context,
	file io.ReaderAt,
	size int64,
) ([]Page, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	contentType, err := detectMIMEType(file, size)
	if err != nil {
		return nil, err
	}

	extractor, ok := ig.extractors[contentType]
	if !ok {
		return nil, fmt.Errorf("ingestion: unsupported MIME type: %s", contentType)
	}

	return extractor.Extract(ctx, file, size)
}

func detectMIMEType(file io.ReaderAt, size int64) (string, error) {
	const limit int64 = 1024 * 1024

	data, err := io.ReadAll(io.NewSectionReader(file, 0, min(size, limit)))
	if err != nil {
		return "", err
	}

	contentType, _, err := mime.ParseMediaType(
		mimetype.Detect(data).String(),
	)
	if err != nil {
		return "", fmt.Errorf("ingestion: parse MIME type: %w", err)
	}

	return contentType, nil
}
