package modules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRegistryGetModule(t *testing.T) {
	r := NewRegistry()

	tests := []struct {
		name     string
		category string
		modName  string
		wantNil  bool
	}{
		{
			name:     "postgres module exists",
			category: "db",
			modName:  "postgres",
			wantNil:  false,
		},
		{
			name:     "mongo module exists",
			category: "db",
			modName:  "mongo",
			wantNil:  false,
		},
		{
			name:     "openapi gen module exists",
			category: "openapi",
			modName:  "gen",
			wantNil:  false,
		},
		{
			name:     "docker module exists",
			category: "docker",
			modName:  "",
			wantNil:  false,
		},
		{
			name:     "non-existent module",
			category: "db",
			modName:  "nonexistent",
			wantNil:  true,
		},
		{
			name:     "non-existent category",
			category: "nonexistent",
			modName:  "test",
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := r.GetModule(tt.category, tt.modName)
			if (mod == nil) != tt.wantNil {
				t.Errorf("GetModule() = %v, wantNil %v", mod, tt.wantNil)
			}
		})
	}
}

func TestPostgresModuleName(t *testing.T) {
	mod := &PostgresModule{}
	if mod.Name() != "postgres" {
		t.Errorf("PostgresModule.Name() = %v, want postgres", mod.Name())
	}
}

func TestMongoModuleName(t *testing.T) {
	mod := &MongoModule{}
	if mod.Name() != "mongo" {
		t.Errorf("MongoModule.Name() = %v, want mongo", mod.Name())
	}
}

func TestOpenAPIGenModuleName(t *testing.T) {
	mod := &OpenAPIGenModule{}
	if mod.Name() != "openapi-gen" {
		t.Errorf("OpenAPIGenModule.Name() = %v, want openapi-gen", mod.Name())
	}
}

func TestOpenAPIManualModuleName(t *testing.T) {
	mod := &OpenAPIManualModule{}
	if mod.Name() != "openapi-manual" {
		t.Errorf("OpenAPIManualModule.Name() = %v, want openapi-manual", mod.Name())
	}
}

func TestDockerModuleName(t *testing.T) {
	mod := &DockerModule{}
	if mod.Name() != "docker" {
		t.Errorf("DockerModule.Name() = %v, want docker", mod.Name())
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()

	testPath := filepath.Join(tmpDir, "subdir", "test.txt")
	testContent := "test content"

	err := WriteFile(testPath, testContent)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Error("WriteFile() did not create file")
	}

	// Verify content
	content, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("file content = %v, want %v", string(content), testContent)
	}
}

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     map[string]interface{}
		want     string
		wantErr  bool
	}{
		{
			name:     "simple substitution",
			template: "module {{.ModulePath}}",
			data:     map[string]interface{}{"ModulePath": "github.com/test/project"},
			want:     "module github.com/test/project",
			wantErr:  false,
		},
		{
			name: "conditional - true",
			template: `{{if eq .Database "postgres"}}
DATABASE_URL=postgres://localhost:5432/db
{{end}}`,
			data:    map[string]interface{}{"Database": "postgres"},
			want:    "\nDATABASE_URL=postgres://localhost:5432/db\n",
			wantErr: false,
		},
		{
			name: "conditional - false",
			template: `{{if eq .Database "postgres"}}
DATABASE_URL=postgres://localhost:5432/db
{{end}}`,
			data:    map[string]interface{}{"Database": "none"},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderTemplate(tt.template, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("renderTemplate() = %q, want %q", string(got), tt.want)
			}
		})
	}
}

// Benchmark tests
func BenchmarkRenderTemplate(b *testing.B) {
	template := "module {{.ModulePath}}\n{{if eq .Database \"postgres\"}}DATABASE_URL=postgres://localhost{{end}}"
	data := map[string]interface{}{
		"ModulePath": "github.com/test/project",
		"Database":   "postgres",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := renderTemplate(template, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteFile(b *testing.B) {
	tmpDir := b.TempDir()
	content := "test content"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path := filepath.Join(tmpDir, "test", "file", "path", "file.txt")
		if err := WriteFile(path, content); err != nil {
			b.Fatal(err)
		}
	}
}
