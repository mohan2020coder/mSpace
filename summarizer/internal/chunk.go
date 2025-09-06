package internal

import (
	"strings"
)

// Split text into chunks by max tokens (approx token ~ word here)
func ChunkText(text string, maxTokens int) []string {
	words := strings.Fields(text)
	var chunks []string
	for i := 0; i < len(words); i += maxTokens {
		end := i + maxTokens
		if end > len(words) {
			end = len(words)
		}
		chunks = append(chunks, strings.Join(words[i:end], " "))
	}
	return chunks
}
