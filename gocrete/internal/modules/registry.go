package modules

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/TRiZKy/gocrete/pkg/templates"
)

var templatesFS = templates.FS

type Module interface {
	Name() string
	Apply(ctx *Context) error
}

type Context struct {
	ProjectPath  string
	Options      InitOptions
	TemplateData map[string]interface{}
}

type InitOptions struct {
	ProjectName string
	ModulePath  string
	Router      string
	Database    string
	OpenAPI     string
	SpecPath    string
	Docker      bool
	Migrations  string
	Force       bool
}

type Registry struct {
	modules map[string]map[string]Module
}

func NewRegistry() *Registry {
	r := &Registry{
		modules: make(map[string]map[string]Module),
	}

	// Register database modules
	r.Register("db", "postgres", &PostgresModule{})
	r.Register("db", "mongo", &MongoModule{})

	// Register OpenAPI modules
	r.Register("openapi", "gen", &OpenAPIGenModule{})
	r.Register("openapi", "manual", &OpenAPIManualModule{})

	// Register Docker module
	r.Register("docker", "", &DockerModule{})

	return r
}

func (r *Registry) Register(category, name string, module Module) {
	if r.modules[category] == nil {
		r.modules[category] = make(map[string]Module)
	}
	r.modules[category][name] = module
}

func (r *Registry) GetModule(category, name string) Module {
	if mods, ok := r.modules[category]; ok {
		return mods[name]
	}
	return nil
}

// Helper functions for modules

func ApplyModuleTemplate(templatePath, destPath string, data map[string]interface{}) error {
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
			content, err = renderTemplate(string(content), data)
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

func renderTemplate(content string, data map[string]interface{}) ([]byte, error) {
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

func WriteFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}
