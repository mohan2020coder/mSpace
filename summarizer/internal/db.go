package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func StoreChunkEmbedding(dsn, table, filename string, chunkID int, chunkText string, embedding []float64) error {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	sql := fmt.Sprintf(`
	INSERT INTO %s (filename, chunk_id, text_chunk, embedding)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (filename, chunk_id)
	DO UPDATE SET text_chunk = EXCLUDED.text_chunk, embedding = EXCLUDED.embedding`, table)

	_, err = conn.Exec(context.Background(), sql, filename, chunkID, chunkText, embedding)
	return err
}
