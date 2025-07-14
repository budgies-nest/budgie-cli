package main

import (
	"context"
	_ "embed"
	"os"
	"strings"

	"github.com/budgies-nest/budgie-cli/cmd"
	"github.com/charmbracelet/fang"
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
		RunE:  cmd.RunAsk,
	}

	askCmd.Flags().StringP("system", "s", ".budgie/budgie.system.md", "Path to system instructions file")
	askCmd.Flags().StringP("config", "c", ".budgie/budgie.config.json", "Path to configuration file")
	askCmd.Flags().StringP("output", "o", ".", "Path where to generate result files")
	askCmd.Flags().BoolP("generate", "g", true, "Generate result file")
	askCmd.Flags().StringP("question", "q", "", "User question (required)")
	askCmd.Flags().BoolP("prompt", "p", false, "Interactive TUI prompt mode")
	askCmd.Flags().StringP("use", "u", "", "Path to file to include as additional system message")
	askCmd.Flags().StringP("from", "f", "", "Path to file containing the user question/message")
	askCmd.Flags().BoolP("rag", "r", false, "Enable RAG (Retrieval-Augmented Generation) mode for enhanced responses with document context")

	askCmd.MarkFlagsOneRequired("question", "prompt", "from")

	var generateEmbeddingsCmd = &cobra.Command{
		Use:   "generate-embeddings",
		Short: "Generate embeddings from markdown files in docs directory",
		Long:  "Parse markdown files, create chunks, and generate embeddings for RAG functionality.",
		RunE:  cmd.RunGenerateEmbeddings,
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
		RunE: func(c *cobra.Command, args []string) error {
			return cmd.RunInit(c, args, defaultConfigContent, defaultSystemContent, defaultDocsReadmeContent)
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the version of budgie",
		Long:  "Display the current version of the budgie CLI tool.",
		RunE: func(c *cobra.Command, args []string) error {
			return cmd.RunVersion(c, args, versionContent)
		},
	}

	rootCmd.AddCommand(askCmd)
	rootCmd.AddCommand(generateEmbeddingsCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)

	if err := fang.Execute(context.TODO(), rootCmd); err != nil {
		os.Exit(1)
	}
}

