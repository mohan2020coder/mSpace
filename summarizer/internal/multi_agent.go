package internal

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// LegalSummarizerState manages the summarization workflow state
type LegalSummarizerState struct {
	Chunks           []SearchResult    `json:"chunks"`
	Query            string            `json:"query"`
	ChunkSummaries   []string          `json:"chunk_summaries"`
	SectionSummaries map[string]string `json:"section_summaries"`
	FinalSummary     string            `json:"final_summary"`
	Errors           []error           `json:"errors"`
}

// MultiAgentLegalSummarizer orchestrates the summarization process
type MultiAgentLegalSummarizer struct {
	llm       llms.Model
	baseURL   string
	modelName string
}

func NewMultiAgentLegalSummarizer(baseURL, modelName string) (*MultiAgentLegalSummarizer, error) {
	llm, err := ollama.New(
		ollama.WithServerURL(baseURL),
		ollama.WithModel(modelName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}

	return &MultiAgentLegalSummarizer{
		llm:       llm,
		baseURL:   baseURL,
		modelName: modelName,
	}, nil
}

// generateText is a helper function to generate text with proper error handling
func (m *MultiAgentLegalSummarizer) generateText(ctx context.Context, prompt string, maxTokens int) (string, error) {
	completion, err := m.llm.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, prompt),
		},
		llms.WithTemperature(0.1),
		llms.WithMaxTokens(maxTokens),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	// Extract text from the completion
	return completion.Choices[0].Content, nil
}

// ChunkSummarizerAgent summarizes individual chunks
func (m *MultiAgentLegalSummarizer) ChunkSummarizerAgent(ctx context.Context, chunk SearchResult, query string) (string, error) {
	prompt := fmt.Sprintf(`As a legal expert, analyze this document excerpt and extract key information relevant to the query: "%s"

Document Excerpt:
%s

Focus on:
1. Legal principles and precedents
2. Relevant facts and arguments
3. Conclusions and holdings
4. Key citations and references

Provide a concise summary focusing only on information relevant to the query:`, query, chunk.Text)

	return m.generateText(ctx, prompt, 512)
}

// SectionOrganizerAgent identifies and groups chunks by legal sections
func (m *MultiAgentLegalSummarizer) SectionOrganizerAgent(ctx context.Context, chunkSummaries []string, query string) ([]string, error) {
	prompt := fmt.Sprintf(`As a legal expert, analyze these chunk summaries and identify legal sections/topics relevant to the query: "%s"

Chunk Summaries:
%s

Identify 3-5 main legal sections (e.g., "Procedural History", "Legal Analysis", "Findings of Fact", "Conclusions of Law", "Precedents", "Arguments").

Return ONLY the section names, one per line, without any additional text:`, query, strings.Join(chunkSummaries, "\n---\n"))

	response, err := m.generateText(ctx, prompt, 256)
	if err != nil {
		return nil, err
	}

	// Parse section names from response
	var sections []string
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		cleanLine := strings.TrimSpace(line)
		if cleanLine != "" &&
			!strings.HasPrefix(cleanLine, "-") &&
			!strings.HasPrefix(cleanLine, "*") &&
			!strings.HasPrefix(cleanLine, "•") &&
			len(cleanLine) > 3 { // Minimum reasonable section name length
			sections = append(sections, cleanLine)
		}
	}

	// If no sections found, use defaults
	if len(sections) == 0 {
		sections = []string{"Legal Analysis", "Key Findings", "Relevant Precedents"}
	}

	return sections, nil
}

// SectionSummarizerAgent creates comprehensive summaries for each legal section
func (m *MultiAgentLegalSummarizer) SectionSummarizerAgent(ctx context.Context, sectionName string, chunkSummaries []string, query string) (string, error) {
	prompt := fmt.Sprintf(`As a legal expert, synthesize this information for the "%s" section to address the query: "%s"

Relevant information:
%s

Create a comprehensive yet concise summary that:
1. Integrates all relevant points
2. Highlights key legal principles
3. Identifies contradictions or consistencies
4. Notes important precedents or citations

Section Summary:`, sectionName, query, strings.Join(chunkSummaries, "\n---\n"))

	return m.generateText(ctx, prompt, 1024)
}

// FinalSynthesisAgent creates the final comprehensive summary
func (m *MultiAgentLegalSummarizer) FinalSynthesisAgent(ctx context.Context, sectionSummaries map[string]string, query string) (string, error) {
	var sectionsContent strings.Builder
	for section, summary := range sectionSummaries {
		sectionsContent.WriteString(fmt.Sprintf("## %s\n%s\n\n", section, summary))
	}

	prompt := fmt.Sprintf(`As a senior legal expert, synthesize these section summaries to provide a comprehensive answer to the query: "%s"

Section Summaries:
%s

Provide a final comprehensive analysis that:
1. Directly answers the query
2. Integrates insights from all relevant sections
3. Highlights key legal conclusions
4. Notes any limitations or uncertainties
5. Provides practical legal guidance

Final Comprehensive Analysis:`, query, sectionsContent.String())

	return m.generateText(ctx, prompt, 2048)
}

// SummarizeLegalDocument orchestrates the multi-agent summarization process
func (m *MultiAgentLegalSummarizer) SummarizeLegalDocument(ctx context.Context, chunks []SearchResult, query string) (string, error) {
	log.Printf("Starting multi-agent legal summarization for query: %s", query)
	log.Printf("Processing %d chunks", len(chunks))

	// Phase 1: Summarize individual chunks in parallel
	chunkSummaries := make([]string, len(chunks))
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make([]error, 0)

	for i, chunk := range chunks {
		wg.Add(1)
		go func(index int, chunk SearchResult) {
			defer wg.Done()

			summary, err := m.ChunkSummarizerAgent(ctx, chunk, query)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("chunk %d: %w", index, err))
				mu.Unlock()
				return
			}

			mu.Lock()
			chunkSummaries[index] = summary
			log.Printf("Completed chunk %d/%d summarization", index+1, len(chunks))
			mu.Unlock()
		}(i, chunk)
	}
	wg.Wait()

	if len(errors) > 0 {
		log.Printf("Warning: %d chunk summarization errors occurred", len(errors))
		if len(errors) == len(chunks) {
			return "", fmt.Errorf("all chunk summarizations failed: %v", errors[0])
		}
	}

	// Filter out empty summaries
	var validSummaries []string
	for _, summary := range chunkSummaries {
		if summary != "" {
			validSummaries = append(validSummaries, summary)
		}
	}

	if len(validSummaries) == 0 {
		return "", fmt.Errorf("no valid chunk summaries generated")
	}

	log.Printf("Generated %d valid chunk summaries", len(validSummaries))

	// Phase 2: Organize by legal sections
	sections, err := m.SectionOrganizerAgent(ctx, validSummaries, query)
	if err != nil {
		return "", fmt.Errorf("section organization failed: %w", err)
	}
	log.Printf("Identified %d legal sections: %v", len(sections), sections)

	// Phase 3: Summarize each section
	sectionSummaries := make(map[string]string)
	for _, section := range sections {
		summary, err := m.SectionSummarizerAgent(ctx, section, validSummaries, query)
		if err != nil {
			log.Printf("Warning: failed to summarize section %s: %v", section, err)
			continue
		}
		sectionSummaries[section] = summary
		log.Printf("Completed section '%s' summarization", section)
	}

	if len(sectionSummaries) == 0 {
		return "", fmt.Errorf("no section summaries generated")
	}

	// Phase 4: Final synthesis
	log.Printf("Starting final synthesis with %d section summaries", len(sectionSummaries))
	finalSummary, err := m.FinalSynthesisAgent(ctx, sectionSummaries, query)
	if err != nil {
		return "", fmt.Errorf("final synthesis failed: %w", err)
	}

	log.Printf("Multi-agent summarization completed successfully")
	return finalSummary, nil
}

// Simple fallback summarizer for comparison
func (m *MultiAgentLegalSummarizer) SimpleSummarize(ctx context.Context, chunks []SearchResult, query string) (string, error) {
	var combinedText strings.Builder
	for _, chunk := range chunks {
		combinedText.WriteString(chunk.Text)
		combinedText.WriteString("\n---\n")
	}

	prompt := fmt.Sprintf(`As a legal expert, using the following document excerpts, answer the question:

%s

Question: %s
Answer:`, combinedText.String(), query)

	return m.generateText(ctx, prompt, 2048)
}
