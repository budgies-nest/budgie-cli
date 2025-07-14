package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "tool",
		Short: "A CLI tool with subcommands",
		Long:  "This example demonstrates how to create subcommands using charmbracelet/fang",
	}

	// Create subcommand
	var createCmd = &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new item",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			fmt.Printf("Creating item: %s\n", name)
		},
	}

	// List subcommand
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all items",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing all items:")
			fmt.Println("- Item 1")
			fmt.Println("- Item 2")
			fmt.Println("- Item 3")
		},
	}

	// Delete subcommand with flags
	var force bool
	var deleteCmd = &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete an item",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if force {
				fmt.Printf("Force deleting item: %s\n", name)
			} else {
				fmt.Printf("Deleting item: %s\n", name)
			}
		},
	}
	deleteCmd.Flags().BoolVarP(&force, "force", "f", false, "Force deletion without confirmation")

	// Status subcommand with nested subcommands
	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show status information",
	}

	var statusSystemCmd = &cobra.Command{
		Use:   "system",
		Short: "Show system status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("System Status: OK")
		},
	}

	var statusServiceCmd = &cobra.Command{
		Use:   "service",
		Short: "Show service status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Service Status: Running")
		},
	}

	// Add nested subcommands
	statusCmd.AddCommand(statusSystemCmd)
	statusCmd.AddCommand(statusServiceCmd)

	// Add all subcommands to root
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(statusCmd)

	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}