package extract

import (
	"context"
	"fmt"
	"io"

	"github.com/bndrmrtn/interleaf/internal/ingestion"
	"github.com/ledongthuc/pdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/tenebris-tech/x2md/pdf2md"
)

// PdfExtractor extracts text content from
// a PDF file and tries to extract markdown too.
type PdfExtractor struct{}

func NewPdfExtractor() *PdfExtractor {
	return &PdfExtractor{}
}

// Extract extracts a list of ingestion.Page from
// a PDF document.
func (p *PdfExtractor) Extract(
	ctx context.Context,
	file io.ReaderAt,
	size int64,
) ([]ingestion.Page, error) {
	return p.extract(ctx, file, size)
}

func (p *PdfExtractor) extract(
	ctx context.Context,
	file io.ReaderAt,
	size int64,
) ([]ingestion.Page, error) {
	pdfReader, err := pdf.NewReader(file, size)
	if err != nil {
		return nil, fmt.Errorf("open PDF: %w", err)
	}

	converter := pdf2md.New(
		pdf2md.WithPreserveFormatting(true),
		pdf2md.WithDetectHeadings(true),
		pdf2md.WithDetectLists(true),
		pdf2md.WithExtractImages(false),
		pdf2md.WithScanMode(false),
		pdf2md.WithStrip(),
	)

	var pages []ingestion.Page

	reader := io.NewSectionReader(file, 0, size)

	err = api.ExtractPages(
		ctx,
		reader,
		nil,
		func(page io.Reader, number int) error {
			data, err := io.ReadAll(page)
			if err != nil {
				return fmt.Errorf("read page %d: %w", number, err)
			}

			original, err := pdfReader.Page(number).GetPlainText(nil)
			if err != nil {
				return fmt.Errorf("extract text from page %d: %w", number, err)
			}

			markdown, err := converter.Convert(data)
			if err != nil {
				return fmt.Errorf("convert page %d: %w", number, err)
			}

			pages = append(pages, ingestion.Page{
				Number:          number,
				OriginalContent: cleanText(original),
				MarkdownContent: cleanText(markdown),
			})

			return nil
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("extract PDF pages: %w", err)
	}

	return pages, nil
}
