package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"os"

	"github.com/budgies-nest/budgie-cli/pkg/config"
	"github.com/budgies-nest/budgie-cli/pkg/utils"
	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/budgies-nest/budgie/rag"
	"github.com/openai/openai-go"
	"github.com/spf13/cobra"
)

// RunGenerateEmbeddings handles the generate-embeddings command execution
func RunGenerateEmbeddings(cmd *cobra.Command, args []string) error {
	configFile, _ := cmd.Flags().GetString("config")
	docsPath, _ := cmd.Flags().GetString("docs")
	markdownHierarchy, _ := cmd.Flags().GetBool("markdown-hierarchy")
	markdownSections, _ := cmd.Flags().GetBool("markdown-sections")
	delimiter, _ := cmd.Flags().GetString("delimiter")
	chunkSize, _ := cmd.Flags().GetInt("chunk-size")
	overlap, _ := cmd.Flags().GetInt("overlap")
	extension, _ := cmd.Flags().GetString("extension")
	files, _ := cmd.Flags().GetBool("files")
	vscodeMode, _ := cmd.Flags().GetBool("vscode")

	// Resolve paths based on vscode mode
	resolvedConfigFile, _, err := utils.ResolveBudgiePaths("", configFile, vscodeMode)
	if err != nil {
		return err
	}
	configFile = resolvedConfigFile

	// Resolve docs path in vscode mode
	if vscodeMode {
		workingDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting working directory: %w", err)
		}
		rootPath, err := utils.FindRootBudgieDir(workingDir)
		if err != nil {
			return err
		}
		docsPath = filepath.Join(rootPath, ".budgie", "docs")
	}

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
	config, err := config.LoadConfig(configFile)
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
		agents.WithMemoryVectorStore(filepath.Join(filepath.Dir(configFile), "embeddings.json")),
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

	embeddingsPath := filepath.Join(filepath.Dir(configFile), "embeddings.json")
	fmt.Printf("Successfully generated %d embeddings and saved to %s\n", chunkCount, embeddingsPath)
	return nil
}