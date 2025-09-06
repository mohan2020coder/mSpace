package main

import (
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func DumpPDFContent(inFile, outDir string) error {
	conf := model.NewDefaultConfiguration()

	// nil = all pages
	if err := api.ExtractContentFile(inFile, outDir, nil, conf); err != nil {
		return fmt.Errorf("extract content failed: %w", err)
	}
	return nil
}
