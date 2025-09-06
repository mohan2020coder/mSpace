package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"` // changed from Input to Prompt
}

func GenerateEmbedding(baseURL, model, prompt string) ([]float32, error) {
	req := EmbeddingRequest{Model: model, Prompt: prompt}
	body, _ := json.Marshal(req)

	resp, err := http.Post(fmt.Sprintf("%s/api/embeddings", baseURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var embResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embResp); err != nil {
		return nil, err
	}

	if len(embResp.Embedding) == 0 {
		return nil, fmt.Errorf("empty embedding response")
	}

	return embResp.Embedding, nil
}
