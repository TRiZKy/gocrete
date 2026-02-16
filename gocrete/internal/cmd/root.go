package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gocrete",
	Short: "Gocrete - Modular Go project scaffolding CLI",
	Long: `Gocrete is a modular project scaffolding CLI that generates 
backend Go projects based on user-selected capabilities.

It provides a clean foundation with optional modules for databases,
OpenAPI, Docker, and more.`,
	SilenceUsage: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
}
