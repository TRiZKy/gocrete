//go:build integration
// +build integration

package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEngineValidateInitOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    InitOptions
		wantErr bool
	}{
		{
			name: "valid chi router",
			opts: InitOptions{
				ProjectName: "test",
				ModulePath:  "github.com/test/test",
				Router:      "chi",
				Database:    "none",
				OpenAPI:     "none",
				Migrations:  "none",
			},
			wantErr: false,
		},
		{
			name: "invalid router",
			opts: InitOptions{
				ProjectName: "test",
				ModulePath:  "github.com/test/test",
				Router:      "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid database",
			opts: InitOptions{
				ProjectName: "test",
				ModulePath:  "github.com/test/test",
				Router:      "chi",
				Database:    "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid openapi",
			opts: InitOptions{
				ProjectName: "test",
				ModulePath:  "github.com/test/test",
				Router:      "chi",
				Database:    "none",
				OpenAPI:     "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEngine()
			err := e.validateInitOptions(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInitOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEngineRenderTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     map[string]interface{}
		want     string
		wantErr  bool
	}{
		{
			name:     "simple variable",
			template: "Hello {{.Name}}",
			data:     map[string]interface{}{"Name": "World"},
			want:     "Hello World",
			wantErr:  false,
		},
		{
			name:     "conditional",
			template: "{{if .Show}}Visible{{end}}",
			data:     map[string]interface{}{"Show": true},
			want:     "Visible",
			wantErr:  false,
		},
		{
			name:     "conditional false",
			template: "{{if .Show}}Visible{{end}}",
			data:     map[string]interface{}{"Show": false},
			want:     "",
			wantErr:  false,
		},
		{
			name:     "invalid template",
			template: "{{.InvalidSyntax",
			data:     map[string]interface{}{},
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEngine()
			got, err := e.renderTemplate(tt.template, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("renderTemplate() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestEngineGetModulePath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test go.mod file
	goModContent := `module github.com/test/project

go 1.22
`
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		t.Fatal(err)
	}

	e := NewEngine()
	modulePath, err := e.getModulePath(tmpDir)
	if err != nil {
		t.Fatalf("getModulePath() error = %v", err)
	}

	want := "github.com/test/project"
	if modulePath != want {
		t.Errorf("getModulePath() = %v, want %v", modulePath, want)
	}
}

func TestEngineGetModulePathNoFile(t *testing.T) {
	tmpDir := t.TempDir()

	e := NewEngine()
	_, err := e.getModulePath(tmpDir)
	if err == nil {
		t.Error("getModulePath() expected error for missing go.mod, got nil")
	}
}

// Integration test - requires templates to be embedded
// Run with: go test -tags integration

func TestInitProjectIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")

	e := NewEngine()
	opts := InitOptions{
		ProjectName: "test-project",
		ModulePath:  "github.com/test/project",
		Router:      "chi",
		Database:    "none",
		OpenAPI:     "none",
		Docker:      false,
		Migrations:  "none",
		Force:       false,
	}

	// Note: This test would require Go to be installed
	// and network access for go mod tidy
	err := e.InitProject(projectPath, opts)
	if err != nil {
		t.Fatalf("InitProject() error = %v", err)
	}

	// Verify structure
	expectedFiles := []string{
		"cmd/server/main.go",
		"internal/config/config.go",
		"internal/logger/logger.go",
		"internal/http/server.go",
		"internal/errors/errors.go",
		"go.mod",
		"README.md",
		".gitignore",
	}

	for _, file := range expectedFiles {
		path := filepath.Join(projectPath, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s does not exist", file)
		}
	}
}
