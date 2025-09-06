package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"

	"github.com/jackc/pgx/v5"
)

type SearchResult struct {
	ChunkID int
	Text    string
	Score   float64
}

func CosineSimilarity(vec1, vec2 []float64) (float64, error) {
	if len(vec1) != len(vec2) {
		return 0, errors.New("vector length mismatch")
	}
	var dot, magA, magB float64
	for i := range vec1 {
		dot += vec1[i] * vec2[i]
		magA += vec1[i] * vec1[i]
		magB += vec2[i] * vec2[i]
	}
	if magA == 0 || magB == 0 {
		return 0, errors.New("zero magnitude vector")
	}
	return dot / (math.Sqrt(magA) * math.Sqrt(magB)), nil
}

func SearchRelevantChunks(dsn, table, filename string, queryEmbedding []float64, topK int) ([]SearchResult, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Printf("Failed to connect to DB: %v", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	log.Printf("Querying chunks for file: %s", filename)
	rows, err := conn.Query(context.Background(), fmt.Sprintf("SELECT chunk_id, text_chunk, embedding FROM %s WHERE filename=$1", table), filename)
	if err != nil {
		log.Printf("DB query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	results := []SearchResult{}
	for rows.Next() {
		var chunkID int
		var text string
		var embedding []float64
		err := rows.Scan(&chunkID, &text, &embedding)
		if err != nil {
			log.Printf("Row scan failed, skipping chunk: %v", err)
			continue
		}
		score, err := CosineSimilarity(queryEmbedding, embedding)
		if err != nil {
			log.Printf("Cosine similarity failed for chunk %d: %v", chunkID, err)
			continue
		}
		results = append(results, SearchResult{ChunkID: chunkID, Text: text, Score: score})
	}

	if len(results) == 0 {
		log.Printf("No chunks found for file %s", filename)
		return nil, nil
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	log.Printf("Top %d chunks retrieved by similarity:", topK)
	for i, res := range results {
		if i >= topK {
			break
		}
		log.Printf("Chunk %d: score=%.4f, text snippet=%q", res.ChunkID, res.Score, snippet(res.Text, 100))
	}

	if len(results) > topK {
		results = results[:topK]
	}

	return results, nil
}

// snippet returns a shortened snippet of text for logging purposes.
func snippet(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
