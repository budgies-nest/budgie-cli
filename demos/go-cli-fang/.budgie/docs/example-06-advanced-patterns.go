package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "advanced",
		Short: "Advanced CLI patterns with charmbracelet/fang",
		Long:  "This example demonstrates advanced patterns and features for CLI development",
	}

	// Command with custom validation
	var validateCmd = &cobra.Command{
		Use:   "validate [email]",
		Short: "Validate email address",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			email := args[0]
			if len(email) == 0 {
				return fmt.Errorf("email cannot be empty")
			}
			if !contains(email, "@") {
				return fmt.Errorf("invalid email format")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Valid email: %s\n", args[0])
		},
	}

	// Command with progress indication
	var progressCmd = &cobra.Command{
		Use:   "progress",
		Short: "Show progress indication",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Processing...")
			for i := 0; i <= 100; i += 10 {
				fmt.Printf("\rProgress: %d%%", i)
				time.Sleep(200 * time.Millisecond)
			}
			fmt.Println("\nCompleted!")
		},
	}

	// Command with conditional execution
	var conditionalCmd = &cobra.Command{
		Use:   "conditional",
		Short: "Conditional command execution",
		Run: func(cmd *cobra.Command, args []string) {
			if os.Getenv("DEBUG") != "" {
				fmt.Println("Debug mode enabled")
			}
			
			if _, err := os.Stat("config.json"); err == nil {
				fmt.Println("Config file found")
			} else {
				fmt.Println("Config file not found")
			}
		},
	}

	// Command with error handling patterns
	var errorCmd = &cobra.Command{
		Use:   "error [type]",
		Short: "Demonstrate error handling",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errorType := args[0]
			
			switch errorType {
			case "warning":
				fmt.Fprintf(os.Stderr, "Warning: This is a warning message\n")
			case "error":
				fmt.Fprintf(os.Stderr, "Error: This is an error message\n")
				os.Exit(1)
			case "info":
				fmt.Println("Info: This is an informational message")
			default:
				fmt.Fprintf(os.Stderr, "Unknown error type: %s\n", errorType)
				os.Exit(1)
			}
		},
	}

	// Command with persistent flags
	var globalVerbose bool
	rootCmd.PersistentFlags().BoolVarP(&globalVerbose, "verbose", "v", false, "Enable verbose output")

	var infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Show application info",
		Run: func(cmd *cobra.Command, args []string) {
			if globalVerbose {
				fmt.Println("Verbose mode enabled")
				fmt.Println("Application: advanced-cli")
				fmt.Println("Version: 1.0.0")
				fmt.Println("Build: 2024-01-01")
			} else {
				fmt.Println("advanced-cli v1.0.0")
			}
		},
	}

	// Command with custom help
	var helpCmd = &cobra.Command{
		Use:   "custom-help",
		Short: "Command with custom help format",
		Long: `This command demonstrates custom help formatting.

Examples:
  advanced custom-help --option value
  advanced custom-help -o value

For more information, visit: https://example.com`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Custom help command executed")
		},
	}
	
	var option string
	helpCmd.Flags().StringVarP(&option, "option", "o", "", "Custom option")

	// Command with aliases
	var aliasCmd = &cobra.Command{
		Use:     "status",
		Aliases: []string{"stat", "st"},
		Short:   "Show status (aliases: stat, st)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Status: Running")
		},
	}

	// Command with custom completion
	var completeCmd = &cobra.Command{
		Use:   "complete [resource]",
		Short: "Command with custom completion",
		Args:  cobra.ExactArgs(1),
		ValidArgs: []string{"user", "project", "task"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Selected resource: %s\n", args[0])
		},
	}

	// Add all commands
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(progressCmd)
	rootCmd.AddCommand(conditionalCmd)
	rootCmd.AddCommand(errorCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(helpCmd)
	rootCmd.AddCommand(aliasCmd)
	rootCmd.AddCommand(completeCmd)

	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}