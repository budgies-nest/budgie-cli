package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/budgies-nest/budgie/agents"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/huh"
	"github.com/openai/openai-go"
	"github.com/spf13/cobra"
)

type Config struct {
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	BaseURL     string  `json:"baseURL"`
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

	return &config, nil
}

func processQuestion(question, systemFile, configFile, outputPath string, generate bool) error {
	config, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("error loading config file: %w", err)
	}

	systemInstructions, err := os.ReadFile(systemFile)
	if err != nil {
		return fmt.Errorf("error reading system instructions file: %w", err)
	}

	agent, err := agents.NewAgent("budgie",
		agents.WithDMR(config.BaseURL),
		agents.WithParams(openai.ChatCompletionNewParams{
			Model:       config.Model,
			Temperature: openai.Opt(config.Temperature),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(string(systemInstructions)),
				openai.UserMessage(question),
			},
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

		fmt.Printf("Result saved to: %s\n", filepath)
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

		for {
			var userInput string
			err := huh.NewInput().
				Title("What's your question?").
				Description("Enter your question for the AI agent (or '/bye' to exit)").
				Value(&userInput).
				Run()
			if err != nil {
				return fmt.Errorf("error getting user input: %w", err)
			}

			if userInput == "/bye" {
				fmt.Println("Goodbye!")
				break
			}

			if userInput == "" {
				fmt.Println("Please enter a question or '/bye' to exit")
				continue
			}

			// Add user message to conversation history
			messages = append(messages, openai.UserMessage(userInput))

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
					fmt.Printf("Result saved to: %s\n", filepath)
				}
			}
			
			fmt.Println()
		}
		return nil
	}

	if question == "" {
		return fmt.Errorf("question is required")
	}

	return processQuestion(question, systemFile, configFile, outputPath, generate)
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
	
	askCmd.MarkFlagsOneRequired("question", "prompt")

	rootCmd.AddCommand(askCmd)

	if err := fang.Execute(context.TODO(), rootCmd); err != nil {
		os.Exit(1)
	}
}
