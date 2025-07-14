package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// RunInit handles the init command execution
func RunInit(cmd *cobra.Command, args []string, defaultConfigContent, defaultSystemContent, defaultDocsReadmeContent string) error {
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