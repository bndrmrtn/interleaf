package ingestion

import (
	"context"
	"io"
)

type Pipeline struct {
	extractors *ExtractorRegistry
	chunker    Chunker
}

func NewPipeline(extractors *ExtractorRegistry, chunker Chunker) *Pipeline {
	return &Pipeline{
		extractors: extractors,
		chunker:    chunker,
	}
}

func (p *Pipeline) Prepare(
	ctx context.Context,
	file io.ReaderAt,
	size int64,
) ([]Page, []Chunk, error) {
	pages, err := p.extractors.Extract(ctx, file, size)
	if err != nil {
		return nil, nil, err
	}

	var chunks []Chunk

	for _, page := range pages {
		parts, err := p.chunker.Split(ctx, page.MarkdownContent)
		if err != nil {
			return nil, nil, err
		}

		for i, part := range parts {
			chunks = append(chunks, Chunk{
				PageNumber: page.Number,
				Index:      i,
				Content:    part,
			})
		}
	}

	return pages, chunks, nil
}
