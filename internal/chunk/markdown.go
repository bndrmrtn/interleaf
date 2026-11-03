package chunk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var wordPattern = regexp.MustCompile(`\S+`)

type MarkdownChunker struct {
	config Config
}

func NewMarkdownChunker() *MarkdownChunker {
	c, _ := NewMarkdownChunkerWithConfig(DefaultConfig())
	return c
}

func NewMarkdownChunkerWithConfig(cfg Config) (*MarkdownChunker, error) {
	if cfg.maxWords <= 0 ||
		cfg.overlapWords < 0 ||
		cfg.overlapWords >= cfg.maxWords ||
		cfg.maxChunks < 0 {
		return nil, errors.New("invalid chunker configuration")
	}

	return &MarkdownChunker{config: cfg}, nil
}

func (c *MarkdownChunker) Split(
	ctx context.Context,
	content string,
) ([]string, error) {
	words := wordPattern.FindAllStringIndex(content, -1)
	var chunks []string

	for start := 0; start < len(words); {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if c.config.maxChunks > 0 && len(chunks) >= c.config.maxChunks {
			return nil, fmt.Errorf("maximum chunk count exceeded: %d", c.config.maxChunks)
		}

		end := min(start+c.config.maxWords, len(words))

		text := content[words[start][0]:words[end-1][1]]
		chunks = append(chunks, strings.TrimSpace(text))

		if end == len(words) {
			break
		}

		start = end - c.config.overlapWords
	}

	return chunks, nil
}
