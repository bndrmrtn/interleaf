package extract

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/bndrmrtn/interleaf/internal/ingestion"
	"github.com/tenebris-tech/x2md/docx2md"
)

type DocxExtractor struct{}

func NewDocxExtractor() *DocxExtractor {
	return &DocxExtractor{}
}

// Extract extracts document content from a DOCX file.
func (d *DocxExtractor) Extract(
	ctx context.Context,
	file io.ReaderAt,
	size int64,
) ([]ingestion.Page, error) {
	return d.extract(ctx, file, size)
}

func (d *DocxExtractor) extract(
	ctx context.Context,
	file io.ReaderAt,
	size int64,
) ([]ingestion.Page, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(io.NewSectionReader(file, 0, size))
	if err != nil {
		return nil, fmt.Errorf("read DOCX: %w", err)
	}

	converter := docx2md.New(
		docx2md.WithPreserveFormatting(true),
		docx2md.WithPreserveImages(false),
	)

	markdown, err := converter.Convert(data)
	if err != nil {
		return nil, fmt.Errorf("convert DOCX: %w", err)
	}

	original, err := d.extractDOCXText(data)
	if err != nil {
		return nil, fmt.Errorf("extract DOCX text: %w", err)
	}

	return []ingestion.Page{
		{
			Number:          1,
			OriginalContent: cleanText(original),
			MarkdownContent: cleanText(markdown),
		},
	}, nil
}

func (d *DocxExtractor) extractDOCXText(data []byte) (string, error) {
	archive, err := zip.NewReader(
		bytes.NewReader(data),
		int64(len(data)),
	)
	if err != nil {
		return "", err
	}

	const namespace = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"

	for _, entry := range archive.File {
		if entry.Name != "word/document.xml" {
			continue
		}

		file, err := entry.Open()
		if err != nil {
			return "", err
		}
		defer mustClose(file)

		decoder := xml.NewDecoder(file)
		var text strings.Builder

		for {
			token, err := decoder.Token()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return "", err
			}

			switch element := token.(type) {
			case xml.StartElement:
				if element.Name.Space != namespace {
					continue
				}

				switch element.Name.Local {
				case "t":
					var value string
					if err := decoder.DecodeElement(&value, &element); err != nil {
						return "", err
					}
					text.WriteString(value)

				case "tab":
					text.WriteByte('\t')

				case "br", "cr":
					text.WriteByte('\n')
				}

			case xml.EndElement:
				if element.Name.Space != namespace {
					continue
				}

				switch element.Name.Local {
				case "p", "tr":
					text.WriteByte('\n')
				case "tc":
					text.WriteByte('\t')
				}
			}
		}

		return text.String(), nil
	}

	return "", fmt.Errorf("DOCX document.xml not found")
}
