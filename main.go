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
	"github.com/eiannone/keyboard"
	"github.com/openai/openai-go"
	"github.com/spf13/cobra"
)

//go:embed templates/budgie.config.json
var defaultConfigContent string

//go:embed templates/budgie.system.md
var defaultSystemContent string

//go:embed templates/docs-readme.md
var defaultDocsReadmeContent string

//go:embed version.txt
var versionContent string

type Config struct {
	Model          string  `json:"model"`
	EmbeddingModel string  `json:"embedding-model"`
	CosineLimit    float64 `json:"cosine-limit"`
	Temperature    float64 `json:"temperature"`
	BaseURL        string  `json:"baseURL"`
}

func setupEscListener(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		err := keyboard.Open()
		if err != nil {
			return
		}
		defer keyboard.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				char, key, err := keyboard.GetKey()
				if err != nil {
					continue
				}

				// Check for ESC key
				if key == keyboard.KeyEsc || char == 27 {
					fmt.Print("\n🛑 Stream stopped by user (ESC pressed)\n")
					cancel()
					return
				}
			}
		}
	}()
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

	fmt.Println("💡 Press ESC to stop streaming")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupEscListener(ctx, cancel)

	var responseBuilder strings.Builder
	_, err = agent.ChatCompletionStream(ctx, func(self *agents.Agent, content string, err error) error {
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
	fromFile, _ := cmd.Flags().GetString("from")

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

		// Handle --from flag in interactive mode - trigger completion immediately
		if fromFile != "" {
			fileContent, err := os.ReadFile(fromFile)
			if err != nil {
				return fmt.Errorf("error reading from file %s: %w", fromFile, err)
			}
			userInput := string(fileContent)

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
				return fmt.Errorf("error creating agent: %v", err)
			}

			fmt.Println("💡 Press ESC to stop streaming")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			setupEscListener(ctx, cancel)

			var responseBuilder strings.Builder
			_, err = agent.ChatCompletionStream(ctx, func(self *agents.Agent, content string, err error) error {
				if err != nil {
					return err
				}
				fmt.Print(content)
				responseBuilder.WriteString(content)
				return nil
			})

			if err != nil {
				return fmt.Errorf("error during streaming: %v", err)
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

		for {
			var userInput string
			err := huh.NewInput().
				Title("What's your question?").
				Description("Enter your question for the AI agent ('/bye' to exit, '/clear' to reset, '/use <file>' to load file, '/from <file>' to ask from file, '#rag' prefix for RAG search)").
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

			if strings.HasPrefix(userInput, "/from ") {
				filePath := strings.TrimPrefix(userInput, "/from ")
				filePath = strings.TrimSpace(filePath)

				if filePath == "" {
					fmt.Println("❌ Please specify a file path: /from <file-path>")
					fmt.Println()
					continue
				}

				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("❌ Error reading file %s: %v\n", filePath, err)
					fmt.Println()
					continue
				}

				// Process the file content as a user question
				userInput = string(fileContent)
				fmt.Printf("📁 Loaded question from file: %s\n", filePath)
				// Don't continue here - let it fall through to process the question
			}

			if userInput == "" {
				fmt.Println("Please enter a question, '/clear' to reset, '/use <file>' to load file, '/from <file>' to ask question from file, '/bye' to exit, or prefix with '#rag' for RAG search")
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

			fmt.Println("💡 Press ESC to stop streaming")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			setupEscListener(ctx, cancel)

			var responseBuilder strings.Builder
			_, err = agent.ChatCompletionStream(ctx, func(self *agents.Agent, content string, err error) error {
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

	// Handle --from flag for single question mode
	if fromFile != "" {
		fileContent, err := os.ReadFile(fromFile)
		if err != nil {
			return fmt.Errorf("error reading from file %s: %w", fromFile, err)
		}
		question = string(fileContent)
	}

	if question == "" {
		return fmt.Errorf("question is required (either via -q flag or -f flag)")
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
	markdownHierarchy, _ := cmd.Flags().GetBool("markdown-hierarchy")
	markdownSections, _ := cmd.Flags().GetBool("markdown-sections")
	delimiter, _ := cmd.Flags().GetString("delimiter")
	chunkSize, _ := cmd.Flags().GetInt("chunk-size")
	overlap, _ := cmd.Flags().GetInt("overlap")
	extension, _ := cmd.Flags().GetString("extension")
	files, _ := cmd.Flags().GetBool("files")

	// Validate that only one chunking method is selected
	chunkingMethods := 0
	if markdownHierarchy {
		chunkingMethods++
	}
	if markdownSections {
		chunkingMethods++
	}
	if delimiter != "" {
		chunkingMethods++
	}
	if chunkSize > 0 {
		chunkingMethods++
	}
	if files {
		chunkingMethods++
	}
	
	if chunkingMethods > 1 {
		return fmt.Errorf("cannot use multiple chunking methods simultaneously (--markdown-hierarchy, --markdown-sections, --delimiter, --chunk-size, --files)")
	}

	// Validate chunk-size and overlap combination
	if overlap > 0 && chunkSize == 0 {
		return fmt.Errorf("--overlap flag requires --chunk-size to be specified")
	}
	if chunkSize > 0 && overlap >= chunkSize {
		return fmt.Errorf("--overlap (%d) must be less than --chunk-size (%d)", overlap, chunkSize)
	}

	// Validate extension flag usage
	if extension != "" && delimiter == "" && chunkSize == 0 && !files {
		return fmt.Errorf("--extension flag can only be used with --delimiter, --chunk-size, or --files methods")
	}

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
	if files {
		fmt.Println("Using whole-file chunking (each file is one chunk)")
	} else if chunkSize > 0 {
		if overlap > 0 {
			fmt.Printf("Using fixed-size text chunking with size: %d, overlap: %d\n", chunkSize, overlap)
		} else {
			fmt.Printf("Using fixed-size text chunking with size: %d\n", chunkSize)
		}
	} else if delimiter != "" {
		fmt.Printf("Using delimiter-based chunking with delimiter: %q\n", delimiter)
	} else if markdownSections {
		fmt.Println("Using markdown sections chunking")
	} else {
		fmt.Println("Using markdown hierarchy chunking (default)")
	}

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

	// Determine file extension to search for
	fileExtension := ".md" // default
	if extension != "" && (delimiter != "" || chunkSize > 0 || files) {
		fileExtension = extension
		if !strings.HasPrefix(fileExtension, ".") {
			fileExtension = "." + fileExtension
		}
	}

	// Find all files with the specified extension in docs directory
	foundFiles, err := helpers.FindFiles(docsPath, fileExtension)
	if err != nil {
		return fmt.Errorf("error finding files with extension %s: %w", fileExtension, err)
	}

	fmt.Printf("Found %d files with extension %s\n", len(foundFiles), fileExtension)

	chunkCount := 0
	for _, filePath := range foundFiles {
		fmt.Printf("Processing: %s\n", filePath)

		// Read file content
		content, err := helpers.ReadTextFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			continue
		}

		// Create chunks based on selected chunking method
		var chunks []string
		if files {
			// For files method, the entire file content is one chunk
			chunks = []string{content}
		} else if chunkSize > 0 {
			chunks = rag.ChunkText(content, chunkSize, overlap)
		} else if delimiter != "" {
			chunks = rag.SplitTextWithDelimiter(content, delimiter)
		} else if markdownSections {
			chunks = rag.SplitMarkdownBySections(content)
		} else {
			// Default to hierarchy chunking (when markdownHierarchy is true or neither flag is set)
			chunks = rag.ChunkWithMarkdownHierarchy(content)
		}

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

func runVersion(cmd *cobra.Command, args []string) error {
	version := strings.TrimSpace(versionContent)
	fmt.Printf("budgie version %s\n", version)
	return nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:     "budgie",
		Short:   "A CLI tool for AI-powered conversations",
		Long:    "qai is a command-line interface that enables AI-powered conversations using configurable models and system instructions.",
		Version: strings.TrimSpace(versionContent),
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
	askCmd.Flags().StringP("from", "f", "", "Path to file containing the user question/message")

	askCmd.MarkFlagsOneRequired("question", "prompt", "from")

	var generateEmbeddingsCmd = &cobra.Command{
		Use:   "generate-embeddings",
		Short: "Generate embeddings from markdown files in docs directory",
		Long:  "Parse markdown files, create chunks, and generate embeddings for RAG functionality.",
		RunE:  runGenerateEmbeddings,
	}

	generateEmbeddingsCmd.Flags().StringP("config", "c", ".budgie/budgie.config.json", "Path to configuration file")
	generateEmbeddingsCmd.Flags().StringP("docs", "d", ".budgie/docs", "Path to docs directory containing markdown files")
	generateEmbeddingsCmd.Flags().BoolP("markdown-hierarchy", "m", false, "Use markdown hierarchy chunking")
	generateEmbeddingsCmd.Flags().BoolP("markdown-sections", "s", false, "Use markdown sections chunking")
	generateEmbeddingsCmd.Flags().StringP("delimiter", "D", "", "Use delimiter-based chunking with specified delimiter")
	generateEmbeddingsCmd.Flags().IntP("chunk-size", "z", 0, "Use fixed-size text chunking with specified size")
	generateEmbeddingsCmd.Flags().IntP("overlap", "o", 0, "Overlap length for fixed-size chunking (requires --chunk-size)")
	generateEmbeddingsCmd.Flags().StringP("extension", "e", "", "File extension to process (for --delimiter, --chunk-size, and --files methods, default: .md)")
	generateEmbeddingsCmd.Flags().BoolP("files", "f", false, "Use whole-file chunking (each file is one chunk)")

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Budgie CLI project",
		Long:  "Create .budgie directory with default configuration, system instructions, and documentation structure.",
		RunE:  runInit,
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the version of budgie",
		Long:  "Display the current version of the budgie CLI tool.",
		RunE:  runVersion,
	}

	rootCmd.AddCommand(askCmd)
	rootCmd.AddCommand(generateEmbeddingsCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)

	if err := fang.Execute(context.TODO(), rootCmd); err != nil {
		os.Exit(1)
	}
}
