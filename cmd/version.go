package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// RunVersion handles the version command execution
func RunVersion(cmd *cobra.Command, args []string, versionContent string) error {
	version := strings.TrimSpace(versionContent)
	fmt.Printf("budgie version %s\n", version)
	return nil
}