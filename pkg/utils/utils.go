package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/eiannone/keyboard"
)

// SetupEscListener sets up a goroutine to listen for ESC key presses to cancel streaming
func SetupEscListener(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		err := keyboard.Open()
		if err != nil {
			return
		}
		defer keyboard.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				char, key, err := keyboard.GetKey()
				if err != nil {
					continue
				}

				// Check for ESC key
				if key == keyboard.KeyEsc || char == 27 {
					fmt.Print("\n🛑 Stream stopped by user (ESC pressed)\n")
					cancel()
					return
				}
			}
		}
	}()
}

// FindRootBudgieDir traverses up the directory tree to find the root .budgie directory
func FindRootBudgieDir(startPath string) (string, error) {
	currentPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", err
	}

	for {
		budgiePath := filepath.Join(currentPath, ".budgie")
		if info, err := os.Stat(budgiePath); err == nil && info.IsDir() {
			return currentPath, nil
		}

		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			return "", fmt.Errorf("no .budgie directory found in current path or any parent directories")
		}
		currentPath = parentPath
	}
}

// ResolveBudgiePaths resolves the system and config file paths, optionally finding root .budgie directory
func ResolveBudgiePaths(systemFile, configFile string, vscodeMode bool) (string, string, error) {
	if !vscodeMode {
		return systemFile, configFile, nil
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("error getting working directory: %w", err)
	}

	rootPath, err := FindRootBudgieDir(workingDir)
	if err != nil {
		return "", "", err
	}

	resolvedSystemFile := filepath.Join(rootPath, ".budgie", "budgie.system.md")
	resolvedConfigFile := filepath.Join(rootPath, ".budgie", "budgie.config.json")

	return resolvedSystemFile, resolvedConfigFile, nil
}