package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TRiZKy/gocrete/internal/engine"
	"github.com/TRiZKy/gocrete/internal/modules"
	"github.com/spf13/cobra"
)

var (
	modulePath string
	router     string
	database   string
	openapi    string
	specPath   string
	docker     bool
	migrations string
	force      bool
)

var initCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Initialize a new Go project",
	Long: `Initialize a new Go backend project with selected modules and capabilities.

Example:
  gocrete init my-service \
    --module github.com/user/my-service \
    --router chi \
    --db postgres \
    --openapi gen \
    --spec ./api.yaml \
    --docker \
    --migrations goose`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		// Validate required flags
		if modulePath == "" {
			return fmt.Errorf("--module flag is required")
		}

		if openapi == "gen" && specPath == "" {
			return fmt.Errorf("--spec flag is required when using --openapi gen")
		}

		// Create project directory
		projectPath := filepath.Join(".", projectName)

		// Check if directory exists
		if _, err := os.Stat(projectPath); err == nil && !force {
			return fmt.Errorf("directory %s already exists (use --force to overwrite)", projectName)
		}

		// Create engine and initialize project
		eng := engine.NewEngine()

		opts := modules.InitOptions{
			ProjectName: projectName,
			ModulePath:  modulePath,
			Router:      router,
			Database:    database,
			OpenAPI:     openapi,
			SpecPath:    specPath,
			Docker:      docker,
			Migrations:  migrations,
			Force:       force,
		}

		fmt.Printf("Initializing project: %s\n", projectName)
		fmt.Printf("Module path: %s\n", modulePath)

		if err := eng.InitProject(projectPath, opts); err != nil {
			return fmt.Errorf("failed to initialize project: %w", err)
		}

		fmt.Printf("\n✓ Project %s created successfully!\n", projectName)
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Printf("  go mod download\n")
		fmt.Printf("  go run cmd/server/main.go\n")

		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&modulePath, "module", "", "Go module path (required)")
	initCmd.Flags().StringVar(&router, "router", "chi", "HTTP router (chi|gin|fiber)")
	initCmd.Flags().StringVar(&database, "db", "none", "Database type (none|postgres|mongo)")
	initCmd.Flags().StringVar(&openapi, "openapi", "none", "OpenAPI mode (none|gen|manual)")
	initCmd.Flags().StringVar(&specPath, "spec", "", "OpenAPI spec path (required if openapi=gen)")
	initCmd.Flags().BoolVar(&docker, "docker", false, "Include Docker configuration")
	initCmd.Flags().StringVar(&migrations, "migrations", "none", "Migration tool (none|goose)")
	initCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing directory")
}
