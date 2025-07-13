package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/budgies-nest/budgie/rag"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/openai/openai-go"
	"github.com/spf13/cobra"
)

//go:embed templates/budgie.config.json
var defaultConfigContent string

//go:embed templates/budgie.system.md
var defaultSystemContent string

//go:embed templates/docs-readme.md
var defaultDocsReadmeContent string

type Config struct {
	Model          string  `json:"model"`
	EmbeddingModel string  `json:"embedding-model"`
	CosineLimit    float64 `json:"cosine-limit"`
	Temperature    float64 `json:"temperature"`
	BaseURL        string  `json:"baseURL"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Set default cosine limit if not specified
	if config.CosineLimit == 0 {
		config.CosineLimit = 0.7
	}

	return &config, nil
}

func createSearchAgent(config *Config) (*agents.Agent, error) {
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

func searchSimilarities(question string, searchAgent *agents.Agent, config *Config) ([]string, error) {
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

func displaySimilarities(similarities []string) {
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

func processQuestion(question, systemFile, configFile, outputPath, useFile string, generate bool) error {
	config, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("error loading config file: %w", err)
	}

	systemInstructions, err := os.ReadFile(systemFile)
	if err != nil {
		return fmt.Errorf("error reading system instructions file: %w", err)
	}

	var similarities []string
	var actualQuestion = question
	
	// Check if RAG search is requested
	if strings.HasPrefix(question, "#rag ") {
		actualQuestion = strings.TrimPrefix(question, "#rag ")
		
		// Create search agent and perform similarity search
		fmt.Print("🔍 Searching... ")
		searchAgent, err := createSearchAgent(config)
		if err != nil {
			fmt.Printf("\nWarning: Error creating search agent: %v\n", err)
		} else if searchAgent != nil {
			similarities, err = searchSimilarities(actualQuestion, searchAgent, config)
			if err != nil {
				fmt.Printf("\nWarning: Error searching similarities: %v\n", err)
			} else {
				fmt.Println("✓")
			}
		}

		// Display similarities in green
		displaySimilarities(similarities)
	}

	// Build messages array starting with system message
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(string(systemInstructions)),
	}

	// Add additional file content as system message if specified
	if useFile != "" {
		useFileContent, err := os.ReadFile(useFile)
		if err != nil {
			return fmt.Errorf("error reading use file %s: %w", useFile, err)
		}
		messages = append(messages, openai.SystemMessage(string(useFileContent)))
	}

	// Add similarity results if found
	if len(similarities) > 0 {
		contextMessage := "Relevant context from documentation:\n\n" + strings.Join(similarities, "\n\n")
		messages = append(messages, openai.UserMessage(contextMessage))
	}

	// Add user question (without #rag prefix if it was used)
	messages = append(messages, openai.UserMessage(actualQuestion))

	agent, err := agents.NewAgent("budgie",
		agents.WithDMR(config.BaseURL),
		agents.WithParams(openai.ChatCompletionNewParams{
			Model:       config.Model,
			Temperature: openai.Opt(config.Temperature),
			Messages:    messages,
		}),
	)
	if err != nil {
		return fmt.Errorf("error creating agent: %w", err)
	}

	var responseBuilder strings.Builder
	_, err = agent.ChatCompletionStream(context.Background(),func(self *agents.Agent, content string, err error) error {
		if err != nil {
			return err
		}
		fmt.Print(content)
		responseBuilder.WriteString(content)
		return nil
	})

	if err != nil {
		return fmt.Errorf("error during streaming: %w", err)
	}

	fmt.Println()

	if generate {
		timestamp := time.Now().Format("2006-01-02-15-04-05")
		filename := fmt.Sprintf("result-%s.md", timestamp)
		filepath := filepath.Join(outputPath, filename)

		err = os.WriteFile(filepath, []byte(responseBuilder.String()), 0644)
		if err != nil {
			return fmt.Errorf("error saving result to file: %w", err)
		}

		blueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
		fmt.Println(blueStyle.Render(fmt.Sprintf("💾 Result saved to: %s", filepath)))
	}

	return nil
}

func runAsk(cmd *cobra.Command, args []string) error {
	systemFile, _ := cmd.Flags().GetString("system")
	configFile, _ := cmd.Flags().GetString("config")
	outputPath, _ := cmd.Flags().GetString("output")
	generate, _ := cmd.Flags().GetBool("generate")
	question, _ := cmd.Flags().GetString("question")
	prompt, _ := cmd.Flags().GetBool("prompt")
	useFile, _ := cmd.Flags().GetString("use")

	if prompt {
		fmt.Println("Interactive mode - type '/bye' to exit")
		fmt.Println()
		
		// Load config and system instructions once for the session
		config, err := loadConfig(configFile)
		if err != nil {
			return fmt.Errorf("error loading config file: %w", err)
		}

		systemInstructions, err := os.ReadFile(systemFile)
		if err != nil {
			return fmt.Errorf("error reading system instructions file: %w", err)
		}

		// Initialize conversation history with system message
		messages := []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(string(systemInstructions)),
		}

		// Add additional file content as system message if specified via flag
		if useFile != "" {
			useFileContent, err := os.ReadFile(useFile)
			if err != nil {
				return fmt.Errorf("error reading use file %s: %w", useFile, err)
			}
			messages = append(messages, openai.SystemMessage(string(useFileContent)))
		}

		for {
			var userInput string
			err := huh.NewInput().
				Title("What's your question?").
				Description("Enter your question for the AI agent ('/bye' to exit, '/clear' to reset, '/use <file>' to load file, '#rag' prefix for RAG search)").
				Value(&userInput).
				Run()
			if err != nil {
				return fmt.Errorf("error getting user input: %w", err)
			}

			if userInput == "/bye" {
				fmt.Println("Goodbye!")
				break
			}

			if userInput == "/clear" {
				// Reload system instructions
				systemInstructions, err := os.ReadFile(systemFile)
				if err != nil {
					fmt.Printf("Error reloading system instructions: %v\n", err)
					continue
				}
				
				// Reset conversation history with new system message
				messages = []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(string(systemInstructions)),
				}

				// Re-add initial use file if it was specified via flag
				if useFile != "" {
					useFileContent, err := os.ReadFile(useFile)
					if err != nil {
						fmt.Printf("Error re-reading use file: %v\n", err)
					} else {
						messages = append(messages, openai.SystemMessage(string(useFileContent)))
					}
				}
				
				fmt.Println("✅ Conversation cleared and system instructions reloaded")
				fmt.Println()
				continue
			}

			if strings.HasPrefix(userInput, "/use ") {
				filePath := strings.TrimPrefix(userInput, "/use ")
				filePath = strings.TrimSpace(filePath)
				
				if filePath == "" {
					fmt.Println("❌ Please specify a file path: /use <file-path>")
					fmt.Println()
					continue
				}

				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("❌ Error reading file %s: %v\n", filePath, err)
					fmt.Println()
					continue
				}

				messages = append(messages, openai.SystemMessage(string(fileContent)))
				fmt.Printf("✅ File %s loaded as system message\n", filePath)
				fmt.Println()
				continue
			}

			if userInput == "" {
				fmt.Println("Please enter a question, '/clear' to reset, '/use <file>' to load file, '/bye' to exit, or prefix with '#rag' for RAG search")
				continue
			}

			var similarities []string
			var actualUserInput = userInput
			
			// Check if RAG search is requested
			if strings.HasPrefix(userInput, "#rag ") {
				actualUserInput = strings.TrimPrefix(userInput, "#rag ")
				
				// Create search agent and perform similarity search
				fmt.Print("🔍 Searching... ")
				searchAgent, err := createSearchAgent(config)
				if err != nil {
					fmt.Printf("\nWarning: Error creating search agent: %v\n", err)
				} else if searchAgent != nil {
					similarities, err = searchSimilarities(actualUserInput, searchAgent, config)
					if err != nil {
						fmt.Printf("\nWarning: Error searching similarities: %v\n", err)
					} else {
						fmt.Println("✓")
					}
				}

				// Display similarities in green
				displaySimilarities(similarities)

				// Add similarity results if found
				if len(similarities) > 0 {
					contextMessage := "Relevant context from documentation:\n\n" + strings.Join(similarities, "\n\n")
					messages = append(messages, openai.SystemMessage(contextMessage))
				}
			}

			// Add user message to conversation history (without #rag prefix if it was used)
			messages = append(messages, openai.UserMessage(actualUserInput))

			// Create agent with current conversation history
			agent, err := agents.NewAgent("budgie",
				agents.WithDMR(config.BaseURL),
				agents.WithParams(openai.ChatCompletionNewParams{
					Model:       config.Model,
					Temperature: openai.Opt(config.Temperature),
					Messages:    messages,
				}),
			)
			if err != nil {
				fmt.Printf("Error creating agent: %v\n", err)
				continue
			}

			var responseBuilder strings.Builder
			_, err = agent.ChatCompletionStream(context.Background(), func(self *agents.Agent, content string, err error) error {
				if err != nil {
					return err
				}
				fmt.Print(content)
				responseBuilder.WriteString(content)
				return nil
			})

			if err != nil {
				fmt.Printf("Error during streaming: %v\n", err)
				continue
			}

			fmt.Println()

			// Add assistant response to conversation history
			assistantResponse := responseBuilder.String()
			messages = append(messages, openai.AssistantMessage(assistantResponse))

			if generate {
				timestamp := time.Now().Format("2006-01-02-15-04-05")
				filename := fmt.Sprintf("result-%s.md", timestamp)
				filepath := filepath.Join(outputPath, filename)

				err = os.WriteFile(filepath, []byte(assistantResponse), 0644)
				if err != nil {
					fmt.Printf("Error saving result to file: %v\n", err)
				} else {
					blueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
					fmt.Println(blueStyle.Render(fmt.Sprintf("💾 Result saved to: %s", filepath)))
				}
			}
			
			fmt.Println()
		}
		return nil
	}

	if question == "" {
		return fmt.Errorf("question is required")
	}

	return processQuestion(question, systemFile, configFile, outputPath, useFile, generate)
}

func runInit(cmd *cobra.Command, args []string) error {
	budgieDir := ".budgie"
	docsDir := filepath.Join(budgieDir, "docs")
	
	// Check if .budgie directory already exists
	if _, err := os.Stat(budgieDir); !os.IsNotExist(err) {
		return fmt.Errorf(".budgie directory already exists")
	}
	
	fmt.Println("🚀 Initializing Budgie CLI project...")
	
	// Create .budgie directory
	if err := os.MkdirAll(budgieDir, 0755); err != nil {
		return fmt.Errorf("error creating .budgie directory: %w", err)
	}
	
	// Create docs directory
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return fmt.Errorf("error creating docs directory: %w", err)
	}
	
	// Write budgie.config.json
	configPath := filepath.Join(budgieDir, "budgie.config.json")
	if err := os.WriteFile(configPath, []byte(defaultConfigContent), 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	
	// Write budgie.system.md
	systemPath := filepath.Join(budgieDir, "budgie.system.md")
	if err := os.WriteFile(systemPath, []byte(defaultSystemContent), 0644); err != nil {
		return fmt.Errorf("error writing system file: %w", err)
	}
	
	// Write docs/README.md
	docsReadmePath := filepath.Join(docsDir, "README.md")
	if err := os.WriteFile(docsReadmePath, []byte(defaultDocsReadmeContent), 0644); err != nil {
		return fmt.Errorf("error writing docs README: %w", err)
	}
	
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	
	fmt.Println(greenStyle.Render("✅ Successfully initialized Budgie CLI project!"))
	fmt.Println()
	fmt.Println("Created:")
	fmt.Printf("  📁 %s/\n", budgieDir)
	fmt.Printf("  ⚙️  %s\n", configPath)
	fmt.Printf("  📝 %s\n", systemPath)
	fmt.Printf("  📁 %s/\n", docsDir)
	fmt.Printf("  📖 %s\n", docsReadmePath)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Add your documentation to .budgie/docs/")
	fmt.Println("  2. Generate embeddings: budgie generate-embeddings")
	fmt.Println("  3. Start asking questions: budgie ask -p")
	
	return nil
}

func runGenerateEmbeddings(cmd *cobra.Command, args []string) error {
	configFile, _ := cmd.Flags().GetString("config")
	docsPath, _ := cmd.Flags().GetString("docs")

	// Load config
	config, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("error loading config file: %w", err)
	}

	if config.EmbeddingModel == "" {
		return fmt.Errorf("embedding-model not specified in config file")
	}

	fmt.Printf("Generating embeddings from docs in: %s\n", docsPath)
	fmt.Printf("Using embedding model: %s\n", config.EmbeddingModel)

	// Create budgie-search agent
	agent, err := agents.NewAgent("budgie-search",
		agents.WithDMR(config.BaseURL),
		agents.WithEmbeddingParams(
			openai.EmbeddingNewParams{
				Model: openai.EmbeddingModel(config.EmbeddingModel),
			},
		),
		agents.WithMemoryVectorStore(".budgie/embeddings.json"),
	)
	if err != nil {
		return fmt.Errorf("error creating agent: %w", err)
	}

	// Reset the vector store
	agent.ResetMemoryVectorStore()

	// Find all markdown files in docs directory
	markdownFiles, err := helpers.FindFiles(docsPath, ".md")
	if err != nil {
		return fmt.Errorf("error finding markdown files: %w", err)
	}

	fmt.Printf("Found %d markdown files\n", len(markdownFiles))

	chunkCount := 0
	for _, filePath := range markdownFiles {
		fmt.Printf("Processing: %s\n", filePath)
		
		// Read markdown file content
		content, err := helpers.ReadTextFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			continue
		}

		// Create chunks using ChunkWithMarkdownHierarchy
		chunks := rag.ChunkWithMarkdownHierarchy(content)
		
		fmt.Printf("  Created %d chunks\n", len(chunks))

		// Create embeddings for each chunk
		for idx, chunk := range chunks {
			chunkID := fmt.Sprintf("%s-chunk-%d", filepath.Base(filePath), idx+1)
			_, err = agent.CreateAndSaveEmbeddingFromText(
				context.Background(),
				chunk,
				chunkID,
			)
			if err != nil {
				fmt.Printf("Error creating embedding for chunk %s: %v\n", chunkID, err)
				continue
			}
			chunkCount++
		}
	}

	// Persist embeddings to .budgie/embeddings.json
	err = agent.PersistMemoryVectorStore()
	if err != nil {
		return fmt.Errorf("error persisting embeddings: %w", err)
	}

	fmt.Printf("Successfully generated %d embeddings and saved to .budgie/embeddings.json\n", chunkCount)
	return nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "budgie",
		Short: "A CLI tool for AI-powered conversations",
		Long:  "qai is a command-line interface that enables AI-powered conversations using configurable models and system instructions.",
	}

	var askCmd = &cobra.Command{
		Use:   "ask",
		Short: "Ask a question to the AI agent",
		Long:  "Ask a question to the AI agent using the configured model and system instructions.",
		RunE:  runAsk,
	}

	askCmd.Flags().StringP("system", "s", ".budgie/budgie.system.md", "Path to system instructions file")
	askCmd.Flags().StringP("config", "c", ".budgie/budgie.config.json", "Path to configuration file")
	askCmd.Flags().StringP("output", "o", ".", "Path where to generate result files")
	askCmd.Flags().BoolP("generate", "g", true, "Generate result file")
	askCmd.Flags().StringP("question", "q", "", "User question (required)")
	askCmd.Flags().BoolP("prompt", "p", false, "Interactive TUI prompt mode")
	askCmd.Flags().StringP("use", "u", "", "Path to file to include as additional system message")
	
	askCmd.MarkFlagsOneRequired("question", "prompt")

	var generateEmbeddingsCmd = &cobra.Command{
		Use:   "generate-embeddings",
		Short: "Generate embeddings from markdown files in docs directory",
		Long:  "Parse markdown files, create chunks, and generate embeddings for RAG functionality.",
		RunE:  runGenerateEmbeddings,
	}

	generateEmbeddingsCmd.Flags().StringP("config", "c", ".budgie/budgie.config.json", "Path to configuration file")
	generateEmbeddingsCmd.Flags().StringP("docs", "d", ".budgie/docs", "Path to docs directory containing markdown files")

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Budgie CLI project",
		Long:  "Create .budgie directory with default configuration, system instructions, and documentation structure.",
		RunE:  runInit,
	}

	rootCmd.AddCommand(askCmd)
	rootCmd.AddCommand(generateEmbeddingsCmd)
	rootCmd.AddCommand(initCmd)

	if err := fang.Execute(context.TODO(), rootCmd); err != nil {
		os.Exit(1)
	}
}
