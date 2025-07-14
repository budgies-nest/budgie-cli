package rag

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/budgies-nest/budgie-cli/pkg/config"
	"github.com/budgies-nest/budgie/agents"
	"github.com/charmbracelet/lipgloss"
	"github.com/openai/openai-go"
)

// CreateSearchAgent creates and configures a search agent for RAG functionality
func CreateSearchAgent(config *config.Config) (*agents.Agent, error) {
	embeddingsPath := ".budgie/embeddings.json"

	// Check if embeddings file exists
	if _, err := os.Stat(embeddingsPath); os.IsNotExist(err) {
		return nil, nil // No embeddings file, return nil
	}

	if config.EmbeddingModel == "" {
		return nil, fmt.Errorf("embedding-model not specified in config file")
	}

	// Create budgie-search agent for similarity search
	searchAgent, err := agents.NewAgent("budgie-search",
		agents.WithDMR(config.BaseURL),
		agents.WithEmbeddingParams(
			openai.EmbeddingNewParams{
				Model: openai.EmbeddingModel(config.EmbeddingModel),
			},
		),
		agents.WithMemoryVectorStore(embeddingsPath),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating search agent: %w", err)
	}

	// Load existing embeddings
	err = searchAgent.LoadMemoryVectorStore()
	if err != nil {
		return nil, fmt.Errorf("error loading vector store: %w", err)
	}

	return searchAgent, nil
}

// SearchSimilarities searches for similar content using the search agent
func SearchSimilarities(question string, searchAgent *agents.Agent, config *config.Config) ([]string, error) {
	if searchAgent == nil {
		return nil, nil // No search agent available
	}

	// Search for similarities
	similarities, err := searchAgent.RAGMemorySearchSimilaritiesWithText(
		context.Background(),
		question,
		config.CosineLimit,
	)
	if err != nil {
		return nil, fmt.Errorf("error searching similarities: %w", err)
	}

	return similarities, nil
}

// DisplaySimilarities displays the found similarities in a formatted way
func DisplaySimilarities(similarities []string) {
	if len(similarities) == 0 {
		fmt.Println("📚 No relevant documentation found")
		fmt.Println()
		return
	}

	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Background(lipgloss.Color("0"))
	contentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	fmt.Println(headerStyle.Render(fmt.Sprintf("📚 Found %d relevant documentation chunks:", len(similarities))))
	fmt.Println()

	for i, similarity := range similarities {
		lines := strings.Split(strings.TrimSpace(similarity), "\n")

		fmt.Printf("%s %d. ", greenStyle.Render("  "), i+1)

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "TITLE:") {
				fmt.Println(greenStyle.Render(line))
			} else if strings.HasPrefix(line, "HIERARCHY:") {
				fmt.Printf("     %s\n", contentStyle.Render(line))
			} else if strings.HasPrefix(line, "CONTENT:") {
				fmt.Printf("     %s\n", contentStyle.Render(line))
			} else {
				// Content continuation
				fmt.Printf("     %s\n", contentStyle.Render(line))
			}
		}
		fmt.Println()
	}
}