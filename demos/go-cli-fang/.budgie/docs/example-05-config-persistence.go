package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

type Config struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Theme    string `json:"theme"`
	Debug    bool   `json:"debug"`
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".myapp", "config.json")
}

func loadConfig() (*Config, error) {
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Theme: "default",
			Debug: false,
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	return &config, err
}

func saveConfig(config *Config) error {
	configPath := getConfigPath()
	configDir := filepath.Dir(configPath)
	
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "config-app",
		Short: "A CLI with configuration and persistence",
		Long:  "This example demonstrates configuration management and data persistence with charmbracelet/fang",
	}

	// Config command group
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage application configuration",
	}

	// Set configuration
	var setCmd = &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value := args[1]

			config, err := loadConfig()
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				return
			}

			switch key {
			case "username":
				config.Username = value
			case "email":
				config.Email = value
			case "theme":
				config.Theme = value
			case "debug":
				config.Debug = value == "true"
			default:
				fmt.Printf("Unknown configuration key: %s\n", key)
				return
			}

			if err := saveConfig(config); err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				return
			}

			fmt.Printf("Set %s = %s\n", key, value)
		},
	}

	// Get configuration
	var getCmd = &cobra.Command{
		Use:   "get [key]",
		Short: "Get a configuration value",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := loadConfig()
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				return
			}

			if len(args) == 0 {
				// Show all config
				fmt.Println("Current configuration:")
				fmt.Printf("  Username: %s\n", config.Username)
				fmt.Printf("  Email: %s\n", config.Email)
				fmt.Printf("  Theme: %s\n", config.Theme)
				fmt.Printf("  Debug: %t\n", config.Debug)
				return
			}

			key := args[0]
			switch key {
			case "username":
				fmt.Println(config.Username)
			case "email":
				fmt.Println(config.Email)
			case "theme":
				fmt.Println(config.Theme)
			case "debug":
				fmt.Printf("%t\n", config.Debug)
			default:
				fmt.Printf("Unknown configuration key: %s\n", key)
			}
		},
	}

	// Reset configuration
	var resetCmd = &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to defaults",
		Run: func(cmd *cobra.Command, args []string) {
			config := &Config{
				Theme: "default",
				Debug: false,
			}

			if err := saveConfig(config); err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				return
			}

			fmt.Println("Configuration reset to defaults")
		},
	}

	// Data persistence example
	var dataCmd = &cobra.Command{
		Use:   "data",
		Short: "Data persistence operations",
	}

	var addDataCmd = &cobra.Command{
		Use:   "add [item]",
		Short: "Add an item to the data store",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			item := args[0]
			dataPath := filepath.Join(filepath.Dir(getConfigPath()), "data.json")

			var items []string
			if data, err := os.ReadFile(dataPath); err == nil {
				json.Unmarshal(data, &items)
			}

			items = append(items, item)
			
			data, _ := json.MarshalIndent(items, "", "  ")
			os.WriteFile(dataPath, data, 0644)

			fmt.Printf("Added item: %s\n", item)
		},
	}

	var listDataCmd = &cobra.Command{
		Use:   "list",
		Short: "List all stored items",
		Run: func(cmd *cobra.Command, args []string) {
			dataPath := filepath.Join(filepath.Dir(getConfigPath()), "data.json")

			var items []string
			if data, err := os.ReadFile(dataPath); err == nil {
				json.Unmarshal(data, &items)
			}

			if len(items) == 0 {
				fmt.Println("No items stored")
				return
			}

			fmt.Println("Stored items:")
			for i, item := range items {
				fmt.Printf("  %d. %s\n", i+1, item)
			}
		},
	}

	// Add subcommands
	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(getCmd)
	configCmd.AddCommand(resetCmd)

	dataCmd.AddCommand(addDataCmd)
	dataCmd.AddCommand(listDataCmd)

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(dataCmd)

	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}