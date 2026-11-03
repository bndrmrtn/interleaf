package ingestion

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type textExtractor struct{}

func (*textExtractor) Extract(
	_ context.Context,
	file io.ReaderAt,
	size int64,
) ([]Page, error) {
	b, err := io.ReadAll(io.NewSectionReader(file, 0, size))
	if err != nil {
		return nil, err
	}

	content := string(b)

	return []Page{
		{
			Number:          1,
			OriginalContent: content,
			MarkdownContent: content,
		},
	}, nil
}

func Test_ExtractString(t *testing.T) {
	ig := NewExtractorRegistry()

	err := ig.RegisterExtractor(&textExtractor{}, "text/plain")
	assert.NoError(t, err, "must be able to register an extractor")

	err = ig.RegisterExtractor(&textExtractor{}, "text/plain")
	assert.ErrorIs(t, err, ErrExtractorAlreadyRegistered, "must return extractor already registered")

	data := []byte("Hello, World!")
	pages, err := ig.Extract(context.Background(), bytes.NewReader(data), int64(len(data)*8))
	assert.NoError(t, err, "plain text extract must not error")

	assert.Equal(t, 1, len(pages), "plain text extract must return exactly one page")

	assert.Equal(t, string(data), pages[0].OriginalContent, "original content must equal to extracted content")
	assert.Equal(t, string(data), pages[0].MarkdownContent, "markdown content must equal to extracted content")
}
