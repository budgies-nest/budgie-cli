package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func main() {
	var name string
	var age int
	var verbose bool

	cmd := &cobra.Command{
		Use:   "greet [message]",
		Short: "A CLI with flags and arguments",
		Long:  "This example demonstrates how to use flags and positional arguments with charmbracelet/fang",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			message := "Hello"
			if len(args) > 0 {
				message = args[0]
			}

			if verbose {
				fmt.Printf("Executing greet command with verbose mode\n")
			}

			if name != "" {
				fmt.Printf("%s, %s!", message, name)
			} else {
				fmt.Printf("%s, there!", message)
			}

			if age > 0 {
				fmt.Printf(" You are %d years old.", age)
			}
			fmt.Println()
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name to greet")
	cmd.Flags().IntVarP(&age, "age", "a", 0, "Age of the person")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	if err := fang.Execute(context.Background(), cmd); err != nil {
		os.Exit(1)
	}
}