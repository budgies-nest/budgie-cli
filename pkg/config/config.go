package config

import (
	"encoding/json"
	"os"
)

// Config represents the application configuration
type Config struct {
	Model          string  `json:"model"`
	EmbeddingModel string  `json:"embedding-model"`
	CosineLimit    float64 `json:"cosine-limit"`
	Temperature    float64 `json:"temperature"`
	BaseURL        string  `json:"baseURL"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
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