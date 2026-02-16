package engine

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/TRiZKy/gocrete/internal/modules"
	"github.com/TRiZKy/gocrete/pkg/templates"
)

var templatesFS = templates.FS

type AddOptions struct {
	Module string
	Type   string
	Mode   string
	Spec   string
}

type Engine struct {
	registry *modules.Registry
}

func NewEngine() *Engine {
	return &Engine{
		registry: modules.NewRegistry(),
	}
}

func (e *Engine) InitProject(projectPath string, opts modules.InitOptions) error {
	// Validate options
	if err := e.validateInitOptions(opts); err != nil {
		return err
	}

	// Create or clean project directory
	if opts.Force {
		if err := os.RemoveAll(projectPath); err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create context
	ctx := &modules.Context{
		ProjectPath: projectPath,
		Options:     opts,
		TemplateData: map[string]interface{}{
			"ProjectName": opts.ProjectName,
			"ModulePath":  opts.ModulePath,
			"Router":      opts.Router,
			"Database":    opts.Database,
			"OpenAPI":     opts.OpenAPI,
			"Migrations":  opts.Migrations,
			"HasDocker":   opts.Docker,
		},
	}

	// Apply base template
	fmt.Println("→ Applying base template...")
	if err := e.applyTemplate("files/base", projectPath, ctx.TemplateData); err != nil {
		return fmt.Errorf("failed to apply base template: %w", err)
	}

	// Apply database module
	if opts.Database != "none" {
		fmt.Printf("→ Adding %s database...\n", opts.Database)
		mod := e.registry.GetModule("db", opts.Database)
		if mod == nil {
			return fmt.Errorf("database module %s not found", opts.Database)
		}
		if err := mod.Apply(ctx); err != nil {
			return fmt.Errorf("failed to apply database module: %w", err)
		}
	}

	// Apply OpenAPI module
	if opts.OpenAPI != "none" {
		fmt.Printf("→ Adding OpenAPI (%s mode)...\n", opts.OpenAPI)
		mod := e.registry.GetModule("openapi", opts.OpenAPI)
		if mod == nil {
			return fmt.Errorf("openapi module %s not found", opts.OpenAPI)
		}
		if err := mod.Apply(ctx); err != nil {
			return fmt.Errorf("failed to apply openapi module: %w", err)
		}
	}

	// Apply Docker module
	if opts.Docker {
		fmt.Println("→ Adding Docker configuration...")
		mod := e.registry.GetModule("docker", "")
		if mod == nil {
			return fmt.Errorf("docker module not found")
		}
		if err := mod.Apply(ctx); err != nil {
			return fmt.Errorf("failed to apply docker module: %w", err)
		}
	}

	// Run post-generation steps
	fmt.Println("→ Running post-generation steps...")
	if err := e.runPostSteps(projectPath, opts); err != nil {
		return fmt.Errorf("post-generation steps failed: %w", err)
	}

	return nil
}

func (e *Engine) AddModule(projectPath string, opts AddOptions) error {
	// Read existing go.mod to get module path
	modPath, err := e.getModulePath(projectPath)
	if err != nil {
		return err
	}

	// Create context
	ctx := &modules.Context{
		ProjectPath: projectPath,
		Options: modules.InitOptions{
			ModulePath: modPath,
		},
		TemplateData: map[string]interface{}{
			"ModulePath": modPath,
		},
	}

	// Get and apply module
	var mod modules.Module
	switch opts.Module {
	case "db":
		if opts.Type == "" {
			return fmt.Errorf("--type flag is required for db module")
		}
		mod = e.registry.GetModule("db", opts.Type)
		ctx.Options.Database = opts.Type
		ctx.TemplateData["Database"] = opts.Type
	case "openapi":
		if opts.Mode == "" {
			return fmt.Errorf("--mode flag is required for openapi module")
		}
		mod = e.registry.GetModule("openapi", opts.Mode)
		ctx.Options.OpenAPI = opts.Mode
		ctx.Options.SpecPath = opts.Spec
		ctx.TemplateData["OpenAPI"] = opts.Mode
	case "docker":
		mod = e.registry.GetModule("docker", "")
		ctx.Options.Docker = true
		ctx.TemplateData["HasDocker"] = true
	default:
		return fmt.Errorf("unknown module: %s", opts.Module)
	}

	if mod == nil {
		return fmt.Errorf("module not found: %s", opts.Module)
	}

	if err := mod.Apply(ctx); err != nil {
		return fmt.Errorf("failed to apply module: %w", err)
	}

	// Run go mod tidy
	fmt.Println("→ Running go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Warning: go mod tidy failed: %s\n", output)
	}

	return nil
}

func (e *Engine) applyTemplate(templatePath, destPath string, data map[string]interface{}) error {
	return fs.WalkDir(templatesFS, templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == templatePath {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(templatePath, path)
		if err != nil {
			return err
		}

		destFilePath := filepath.Join(destPath, relPath)

		if d.IsDir() {
			return os.MkdirAll(destFilePath, 0755)
		}

		// Read file content
		content, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Check if file should be templated
		if strings.HasSuffix(path, ".tmpl") {
			destFilePath = strings.TrimSuffix(destFilePath, ".tmpl")
			content, err = e.renderTemplate(string(content), data)
			if err != nil {
				return fmt.Errorf("failed to render template %s: %w", path, err)
			}
		}

		// Write file
		if err := os.MkdirAll(filepath.Dir(destFilePath), 0755); err != nil {
			return err
		}

		return os.WriteFile(destFilePath, content, 0644)
	})
}

func (e *Engine) renderTemplate(content string, data map[string]interface{}) ([]byte, error) {
	tmpl, err := template.New("template").Parse(content)
	if err != nil {
		return nil, err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}

func (e *Engine) runPostSteps(projectPath string, opts modules.InitOptions) error {
	// Initialize go module
	cmd := exec.Command("go", "mod", "init", opts.ModulePath)
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod init failed: %s", output)
	}

	// Run go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %s", output)
	}

	// Format code
	cmd = exec.Command("go", "fmt", "./...")
	cmd.Dir = projectPath
	cmd.Run() // Ignore errors for formatting

	return nil
}

func (e *Engine) validateInitOptions(opts modules.InitOptions) error {
	// Validate router
	validRouters := map[string]bool{"chi": true, "gin": true, "fiber": true}
	if !validRouters[opts.Router] {
		return fmt.Errorf("invalid router: %s (must be chi, gin, or fiber)", opts.Router)
	}

	// Validate database
	validDatabases := map[string]bool{"none": true, "postgres": true, "mongo": true}
	if !validDatabases[opts.Database] {
		return fmt.Errorf("invalid database: %s (must be none, postgres, or mongo)", opts.Database)
	}

	// Validate OpenAPI
	validOpenAPI := map[string]bool{"none": true, "gen": true, "manual": true}
	if !validOpenAPI[opts.OpenAPI] {
		return fmt.Errorf("invalid openapi: %s (must be none, gen, or manual)", opts.OpenAPI)
	}

	// Validate migrations
	validMigrations := map[string]bool{"none": true, "goose": true}
	if !validMigrations[opts.Migrations] {
		return fmt.Errorf("invalid migrations: %s (must be none or goose)", opts.Migrations)
	}

	return nil
}

func (e *Engine) getModulePath(projectPath string) (string, error) {
	modFile := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(modFile)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "module ") {
			return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "module ")), nil
		}
	}

	return "", fmt.Errorf("module path not found in go.mod")
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
