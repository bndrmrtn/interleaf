package ingestion

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
)

type Ingestion struct {
	extractors []extractor
	config     Config
}

func NewIngestion() *Ingestion {
	return &Ingestion{}
}

func (ig *Ingestion) RegisterExtractor(
	pe PageExtractor,
	contentType string,
	otherContentType ...string,
) {
	ig.extractors = append(ig.extractors, extractor{
		pe:           pe,
		contentTypes: append(otherContentType, contentType),
	})
}

func (ig *Ingestion) Extract(file io.ReaderAt, size int64) ([]Page, error) {
	mime, err := ig.detectMIMEType(file)
	if err != nil {
		return nil, err
	}

	var pe PageExtractor
	for _, ex := range ig.extractors {
		if slices.Contains(ex.contentTypes, mime) {
			pe = ex.pe
			break
		}
	}

	if pe == nil {
		return nil, fmt.Errorf("ingestion: no extractor for mime: %s", mime)
	}

	return pe.Extract(context.TODO(), file, size)
}

func (ig *Ingestion) detectMIMEType(r io.ReaderAt) (string, error) {
	var buf [512]byte

	n, err := r.ReadAt(buf[:], 0)
	if err != nil && err != io.EOF {
		return "", err
	}

	return http.DetectContentType(buf[:n]), nil
}
