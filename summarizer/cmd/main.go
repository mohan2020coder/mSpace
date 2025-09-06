package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"summarizer/internal"
)

func main() {
	pdfPath := flag.String("input", "", "Path to PDF judgment")
	query := flag.String("query", "", "Question or 'summary' to generate summary")
	configPath := flag.String("config", "config.yaml", "Path to config")
	flag.Parse()

	if *pdfPath == "" || *query == "" {
		log.Fatal("Usage: --input file.pdf --query 'summary' or any question")
	}

	cfg, err := internal.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed loading config: %v", err)
	}

	text, err := internal.ExtractTextFromPDF(*pdfPath)
	if err != nil {
		log.Fatalf("PDF extraction failed: %v", err)
	}

	chunks := internal.ChunkText(text, cfg.MaxChunkTokens)

	filename := filepath.Base(*pdfPath)

	// Store embeddings for each chunk (idempotent: skip if exists)
	for i, chunk := range chunks {
		emb32, err := internal.GenerateEmbedding(cfg.Ollama.BaseURL, cfg.EmbeddingModel, chunk)
		if err != nil {
			log.Fatalf("Embedding generation failed: %v", err)
		}
		emb64 := internal.Float32ToFloat64(emb32)
		err = internal.StoreChunkEmbedding(cfg.Database.DSN, cfg.Database.Table, filename, i, chunk, emb64)
		if err != nil {
			log.Fatalf("DB insert failed: %v", err)
		}
	}

	// Generate embedding for user query
	queryEmb32, err := internal.GenerateEmbedding(cfg.Ollama.BaseURL, cfg.EmbeddingModel, *query)
	if err != nil {
		log.Fatalf("Query embedding failed: %v", err)
	}
	queryEmb64 := internal.Float32ToFloat64(queryEmb32)

	// Initialize the multi-agent summarizer
	summarizer, err := internal.NewMultiAgentLegalSummarizer(
		cfg.Ollama.BaseURL,
		cfg.Model,
	)
	if err != nil {
		log.Fatalf("Failed to create summarizer: %v", err)
	}

	// Retrieve relevant chunks
	topChunks, err := internal.SearchRelevantChunks(
		cfg.Database.DSN,
		cfg.Database.Table,
		filename,
		queryEmb64,
		8, // Get more chunks for better context
	)
	if err != nil {
		log.Fatalf("Retrieval failed: %v", err)
	}

	ctx := context.Background()

	// Use multi-agent summarization
	answer, err := summarizer.SummarizeLegalDocument(ctx, topChunks, *query)
	if err != nil {
		log.Printf("Multi-agent summarization failed, trying simple approach: %v", err)
		// Fallback to simple summarization
		answer, err = summarizer.SimpleSummarize(ctx, topChunks, *query)
		if err != nil {
			log.Fatalf("Both summarization methods failed: %v", err)
		}
	}

	fmt.Println("=== COMPREHENSIVE LEGAL ANALYSIS ===")
	fmt.Println(strings.TrimSpace(answer))
}
