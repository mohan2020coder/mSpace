package internal

import (
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractTextFromPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var content strings.Builder
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		content.WriteString(text)
		content.WriteString("\n")
	}

	return normalizeText(content.String()), nil
}

func normalizeText(text string) string {
	// Remove excessive whitespace
	text = strings.Join(strings.Fields(text), " ")

	// Clean up common PDF artifacts
	text = strings.ReplaceAll(text, "  ", " ")
	text = strings.ReplaceAll(text, " .", ".")
	text = strings.ReplaceAll(text, " ,", ",")

	return strings.TrimSpace(text)
}
