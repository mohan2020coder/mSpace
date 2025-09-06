package internal

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type SummarizationRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type SummarizationResponse struct {
	Text string `json:"text"`
}

// BuildContextPrompt remains the same
func BuildContextPrompt(chunks []SearchResult, query string) string {
	var sb strings.Builder
	sb.WriteString("You are a legal expert. Using the following relevant document excerpts, answer the question:\n\n")

	for _, c := range chunks {
		sb.WriteString(c.Text)
		sb.WriteString("\n---\n")
	}
	sb.WriteString("\nQuestion: ")
	sb.WriteString(query)
	sb.WriteString("\nAnswer:")

	prompt := sb.String()

	const maxPromptLength = 4096
	if len(prompt) > maxPromptLength {
		prompt = prompt[:maxPromptLength]
	}

	log.Printf("Constructed prompt (truncated if necessary):\n%s\n", prompt)
	return prompt
}

func SummarizeWithLangChain(ctx context.Context, baseURL, model, prompt string) (string, error) {
	// Create Ollama LLM instance
	llm, err := ollama.New(
		ollama.WithServerURL(baseURL),
		ollama.WithModel(model),
	)
	if err != nil {
		log.Printf("Failed to create Ollama client: %v", err)
		return "", fmt.Errorf("failed to create Ollama client: %w", err)
	}

	log.Printf("Sending summarization request to %s with model %s", baseURL, model)

	// Generate the completion - USE THE CORRECT CONSTANT
	completion, err := llm.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, prompt), // FIXED: Use llms.ChatMessageTypeHuman
		},
		llms.WithTemperature(0.1),
		llms.WithMaxTokens(2048),
	)
	if err != nil {
		log.Printf("Failed to generate content: %v", err)
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(completion.Choices) == 0 {
		log.Printf("No completion choices returned")
		return "", fmt.Errorf("no completion choices returned")
	}

	response := completion.Choices[0].Content
	log.Printf("Summarization response: %s", response)

	return response, nil
}
