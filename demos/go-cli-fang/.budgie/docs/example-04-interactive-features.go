package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "interactive",
		Short: "A CLI with interactive features",
		Long:  "This example demonstrates interactive input handling with charmbracelet/fang",
	}

	// Interactive prompt command
	var promptCmd = &cobra.Command{
		Use:   "prompt",
		Short: "Interactive prompt for user input",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			
			fmt.Print("Enter your name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			
			fmt.Print("Enter your age: ")
			age, _ := reader.ReadString('\n')
			age = strings.TrimSpace(age)
			
			fmt.Printf("Hello %s! You are %s years old.\n", name, age)
		},
	}

	// Confirm command with yes/no prompt
	var confirmCmd = &cobra.Command{
		Use:   "confirm",
		Short: "Confirmation dialog example",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			
			fmt.Print("Are you sure you want to proceed? (y/n): ")
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			
			if response == "y" || response == "yes" {
				fmt.Println("Proceeding with the operation...")
			} else {
				fmt.Println("Operation cancelled.")
			}
		},
	}

	// Menu selection command
	var menuCmd = &cobra.Command{
		Use:   "menu",
		Short: "Interactive menu selection",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			
			fmt.Println("Select an option:")
			fmt.Println("1. Option A")
			fmt.Println("2. Option B")
			fmt.Println("3. Option C")
			fmt.Print("Enter your choice (1-3): ")
			
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)
			
			switch choice {
			case "1":
				fmt.Println("You selected Option A")
			case "2":
				fmt.Println("You selected Option B")
			case "3":
				fmt.Println("You selected Option C")
			default:
				fmt.Println("Invalid choice")
			}
		},
	}

	// Multi-step wizard command
	var wizardCmd = &cobra.Command{
		Use:   "wizard",
		Short: "Multi-step configuration wizard",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			
			fmt.Println("=== Configuration Wizard ===")
			
			// Step 1
			fmt.Print("Step 1/3 - Enter project name: ")
			projectName, _ := reader.ReadString('\n')
			projectName = strings.TrimSpace(projectName)
			
			// Step 2
			fmt.Print("Step 2/3 - Enter project description: ")
			description, _ := reader.ReadString('\n')
			description = strings.TrimSpace(description)
			
			// Step 3
			fmt.Print("Step 3/3 - Choose language (go/python/javascript): ")
			language, _ := reader.ReadString('\n')
			language = strings.TrimSpace(language)
			
			// Summary
			fmt.Println("\n=== Configuration Summary ===")
			fmt.Printf("Project Name: %s\n", projectName)
			fmt.Printf("Description: %s\n", description)
			fmt.Printf("Language: %s\n", language)
			
			fmt.Print("\nSave configuration? (y/n): ")
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(strings.ToLower(confirm))
			
			if confirm == "y" || confirm == "yes" {
				fmt.Println("Configuration saved successfully!")
			} else {
				fmt.Println("Configuration discarded.")
			}
		},
	}

	rootCmd.AddCommand(promptCmd)
	rootCmd.AddCommand(confirmCmd)
	rootCmd.AddCommand(menuCmd)
	rootCmd.AddCommand(wizardCmd)

	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}