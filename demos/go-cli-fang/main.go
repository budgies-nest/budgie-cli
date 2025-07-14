package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func main() {
	var firstName string
	var lastName string

	cmd := &cobra.Command{
		Use:   "hello",
		Short: "A simple hello world CLI",
		Long:  `This is a basic CLI application that demonstrates the simplest use of charmbracelet/fang`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hello, %s %s!\n", firstName, lastName)
		},
	}

	cmd.Flags().StringVar(&firstName, "first-name", "", "First name of the user")
	cmd.Flags().StringVar(&lastName, "last-name", "", "Last name of the user")

	if err := fang.Execute(context.Background(), cmd); err != nil {
		os.Exit(1)
	}
}
