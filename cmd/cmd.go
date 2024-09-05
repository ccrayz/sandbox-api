package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ccrayz/sandbox-api/cmd/server"
)

var rootCmd = &cobra.Command{
	Use:   "ccrazy-cli",
	Short: "App is a CLI application",
	Long:  `A longer description of your CLI application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from your CLI application!")
	},
}

func Execute() {
	rootCmd.AddCommand(server.ServerCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
