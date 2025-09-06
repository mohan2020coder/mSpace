package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration parameters
type Config struct {
	Model          string `yaml:"model"`
	EmbeddingModel string `yaml:"embedding_model"`
	MaxChunkTokens int    `yaml:"max_chunk_tokens"`

	Database struct {
		DSN   string `yaml:"dsn"`
		Table string `yaml:"table"`
	} `yaml:"database"`

	Ollama struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"ollama"`
}

// LoadConfig loads YAML config from file path
func LoadConfig(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}

	// Basic validation
	if cfg.Model == "" || cfg.EmbeddingModel == "" || cfg.MaxChunkTokens <= 0 {
		return nil, fmt.Errorf("invalid config: model, embedding_model and max_chunk_tokens must be set")
	}
	if cfg.Database.DSN == "" || cfg.Database.Table == "" {
		return nil, fmt.Errorf("invalid config: database.dsn and database.table must be set")
	}
	if cfg.Ollama.BaseURL == "" {
		return nil, fmt.Errorf("invalid config: ollama.base_url must be set")
	}

	return &cfg, nil
}
func Float32ToFloat64(input []float32) []float64 {
	result := make([]float64, len(input))
	for i, v := range input {
		result[i] = float64(v)
	}
	return result
}
