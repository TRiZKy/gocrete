package cmd

import (
	"fmt"
	"os"

	"github.com/TRiZKy/gocrete/internal/engine"
	"github.com/spf13/cobra"
)

var (
	addType string
	addMode string
	addSpec string
)

var addCmd = &cobra.Command{
	Use:   "add <module>",
	Short: "Add a module to an existing project",
	Long: `Add a module to an existing Gocrete project.

Examples:
  gocrete add db --type postgres
  gocrete add openapi --mode gen --spec api.yaml
  gocrete add docker`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		moduleName := args[0]

		// Check if we're in a project directory
		if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
			return fmt.Errorf("not in a Go project directory (go.mod not found)")
		}

		// Create engine and add module
		eng := engine.NewEngine()

		opts := engine.AddOptions{
			Module: moduleName,
			Type:   addType,
			Mode:   addMode,
			Spec:   addSpec,
		}

		fmt.Printf("Adding module: %s\n", moduleName)

		if err := eng.AddModule(".", opts); err != nil {
			return fmt.Errorf("failed to add module: %w", err)
		}

		fmt.Printf("\n✓ Module %s added successfully!\n", moduleName)

		return nil
	},
}

func init() {
	addCmd.Flags().StringVar(&addType, "type", "", "Module type (for db: postgres|mongo)")
	addCmd.Flags().StringVar(&addMode, "mode", "", "Module mode (for openapi: gen|manual)")
	addCmd.Flags().StringVar(&addSpec, "spec", "", "Spec path (for openapi gen)")
}
