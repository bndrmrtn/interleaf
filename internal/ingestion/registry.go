package ingestion

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"

	"github.com/gabriel-vasile/mimetype"
)

var ErrExtractorAlreadyRegistered = errors.New("extractor already registered")

type ExtractorRegistry struct {
	extractors map[string]PageExtractor
}

func NewExtractorRegistry() *ExtractorRegistry {
	return &ExtractorRegistry{
		extractors: make(map[string]PageExtractor),
	}
}

func (ig *ExtractorRegistry) RegisterExtractor(
	pe PageExtractor,
	contentType string,
	contentTypes ...string,
) error {
	for _, ct := range append([]string{contentType}, contentTypes...) {
		if _, exists := ig.extractors[ct]; exists {
			return fmt.Errorf("%w: %s", ErrExtractorAlreadyRegistered, ct)
		}
	}

	ig.extractors[contentType] = pe

	for _, contentType := range contentTypes {
		ig.extractors[contentType] = pe
	}

	return nil
}

func (ig *ExtractorRegistry) Extract(
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
