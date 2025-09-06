package main

import (
	"fmt"
	"log"
	"summarizer/internal" // replace with your actual module path
)

func main() {
	pdfPath := "OP_1_2021.pdf" // replace with your PDF file path

	text, err := internal.ExtractTextFromPDF(pdfPath)
	if err != nil {
		log.Fatalf("Failed to extract text: %v", err)
	}

	fmt.Println("Extracted Text:")
	fmt.Println(text)
}
