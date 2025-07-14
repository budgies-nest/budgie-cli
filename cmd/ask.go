package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/budgies-nest/budgie-cli/pkg/config"
	"github.com/budgies-nest/budgie-cli/pkg/rag"
	"github.com/budgies-nest/budgie-cli/pkg/utils"
	"github.com/budgies-nest/budgie/agents"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/openai/openai-go"
	"github.com/spf13/cobra"
)

// processQuestion handles a single question processing workflow
func processQuestion(question, systemFile, configFile, outputPath, useFile string, generate, ragEnabled bool) error {
	config, err := config.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("error loading config file: %w", err)
	}

	systemInstructions, err := os.ReadFile(systemFile)
	if err != nil {
		return fmt.Errorf("error reading system instructions file: %w", err)
	}

	var similarities []string
	var actualQuestion = question

	// Check if RAG search is requested (either via --rag flag or #rag prefix)
	ragRequested := ragEnabled || strings.HasPrefix(question, "#rag ")
	
	if ragRequested {
		// Remove #rag prefix if present (when using --rag flag, #rag prefix is not needed)
		if strings.HasPrefix(question, "#rag ") {
			actualQuestion = strings.TrimPrefix(question, "#rag ")
		}

		// Create search agent and perform similarity search
		fmt.Print("🔍 Searching... ")
		searchAgent, err := rag.CreateSearchAgent(config)
		if err != nil {
			fmt.Printf("\nWarning: Error creating search agent: %v\n", err)
		} else if searchAgent != nil {
			similarities, err = rag.SearchSimilarities(actualQuestion, searchAgent, config)
			if err != nil {
				fmt.Printf("\nWarning: Error searching similarities: %v\n", err)
			} else {
				fmt.Println("✓")
			}
		}

		// Display similarities in green
		rag.DisplaySimilarities(similarities)
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
	utils.SetupEscListener(ctx, cancel)

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

// RunAsk handles the ask command execution
func RunAsk(cmd *cobra.Command, args []string) error {
	systemFile, _ := cmd.Flags().GetString("system")
	configFile, _ := cmd.Flags().GetString("config")
	outputPath, _ := cmd.Flags().GetString("output")
	generate, _ := cmd.Flags().GetBool("generate")
	question, _ := cmd.Flags().GetString("question")
	prompt, _ := cmd.Flags().GetBool("prompt")
	useFile, _ := cmd.Flags().GetString("use")
	fromFile, _ := cmd.Flags().GetString("from")
	ragEnabled, _ := cmd.Flags().GetBool("rag")
	vscodeMode, _ := cmd.Flags().GetBool("vscode")

	// Resolve paths based on vscode mode
	resolvedSystemFile, resolvedConfigFile, err := utils.ResolveBudgiePaths(systemFile, configFile, vscodeMode)
	if err != nil {
		return err
	}
	systemFile = resolvedSystemFile
	configFile = resolvedConfigFile

	if prompt {
		fmt.Println("Interactive mode - type '/bye' to exit")
		fmt.Println()

		// Load config and system instructions once for the session
		config, err := config.LoadConfig(configFile)
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

			// Check if RAG search is requested (either via --rag flag or #rag prefix)
			ragRequested := ragEnabled || strings.HasPrefix(userInput, "#rag ")
			
			if ragRequested {
				// Remove #rag prefix if present (when using --rag flag, #rag prefix is not needed)
				if strings.HasPrefix(userInput, "#rag ") {
					actualUserInput = strings.TrimPrefix(userInput, "#rag ")
				}

				// Create search agent and perform similarity search
				fmt.Print("🔍 Searching... ")
				searchAgent, err := rag.CreateSearchAgent(config)
				if err != nil {
					fmt.Printf("\nWarning: Error creating search agent: %v\n", err)
				} else if searchAgent != nil {
					similarities, err = rag.SearchSimilarities(actualUserInput, searchAgent, config)
					if err != nil {
						fmt.Printf("\nWarning: Error searching similarities: %v\n", err)
					} else {
						fmt.Println("✓")
					}
				}

				// Display similarities in green
				rag.DisplaySimilarities(similarities)

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
			utils.SetupEscListener(ctx, cancel)

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
				Description("Enter your question for the AI agent ('/bye' to exit, '/clear' to reset, '/use <file>' to load file, '/from <file>' to ask from file, '#rag' prefix for RAG search when --rag flag not used)").
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
				fmt.Println("Please enter a question, '/clear' to reset, '/use <file>' to load file, '/from <file>' to ask question from file, '/bye' to exit, or prefix with '#rag' for RAG search (when --rag flag not used)")
				continue
			}

			var similarities []string
			var actualUserInput = userInput

			// Check if RAG search is requested (either via --rag flag or #rag prefix)
			ragRequested := ragEnabled || strings.HasPrefix(userInput, "#rag ")
			
			if ragRequested {
				// Remove #rag prefix if present (when using --rag flag, #rag prefix is not needed)
				if strings.HasPrefix(userInput, "#rag ") {
					actualUserInput = strings.TrimPrefix(userInput, "#rag ")
				}

				// Create search agent and perform similarity search
				fmt.Print("🔍 Searching... ")
				searchAgent, err := rag.CreateSearchAgent(config)
				if err != nil {
					fmt.Printf("\nWarning: Error creating search agent: %v\n", err)
				} else if searchAgent != nil {
					similarities, err = rag.SearchSimilarities(actualUserInput, searchAgent, config)
					if err != nil {
						fmt.Printf("\nWarning: Error searching similarities: %v\n", err)
					} else {
						fmt.Println("✓")
					}
				}

				// Display similarities in green
				rag.DisplaySimilarities(similarities)

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
			utils.SetupEscListener(ctx, cancel)

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

	return processQuestion(question, systemFile, configFile, outputPath, useFile, generate, ragEnabled)
}