# Gocrete Architecture

This document explains how Gocrete works internally and how to extend it.

## Overview

Gocrete is a **code generator** that scaffolds Go projects. It follows a modular architecture with three main layers:

```
┌─────────────────────┐
│    CLI Layer        │  Cobra commands (init, add)
│  (Presentation)     │
└──────────┬──────────┘
           │
┌──────────▼──────────┐
│  Generator Engine   │  Orchestrates generation
│    (Business)       │  Template rendering
└──────────┬──────────┘  Post-processing
           │
┌──────────▼──────────┐
│  Module Registry    │  Pluggable modules
│     (Domain)        │  Template application
└─────────────────────┘
```

## Components

### 1. CLI Layer (`internal/cmd/`)

**Responsibility:** User interface and command parsing

**Files:**
- `root.go` - Root Cobra command and initialization
- `init.go` - Project initialization command
- `add.go` - Module addition command

**Key Functions:**
```go
// init.go
func (cmd *initCmd) RunE(...) error {
    // 1. Validate flags
    // 2. Create InitOptions
    // 3. Call engine.InitProject()
    // 4. Display success message
}
```

**Design Principles:**
- Thin layer - minimal logic
- Delegate to engine
- Clear error messages
- User-friendly output

### 2. Generator Engine (`internal/engine/`)

**Responsibility:** Core generation logic

**Files:**
- `engine.go` - Main engine implementation

**Key Structures:**
```go
type Engine struct {
    registry *modules.Registry
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

type Context struct {
    ProjectPath  string
    Options      InitOptions
    TemplateData map[string]interface{}
}
```

**Workflow:**

```
InitProject()
    │
    ├─> Validate options
    │
    ├─> Create project directory
    │
    ├─> Apply base template
    │   └─> Walk template files
    │       └─> Render .tmpl files
    │           └─> Copy others
    │
    ├─> Apply database module (if selected)
    │
    ├─> Apply OpenAPI module (if selected)
    │
    ├─> Apply Docker module (if selected)
    │
    └─> Run post-generation steps
        ├─> go mod init
        ├─> go mod tidy
        └─> go fmt
```

**Key Methods:**

```go
func (e *Engine) InitProject(path string, opts InitOptions) error
func (e *Engine) AddModule(path string, opts AddOptions) error
func (e *Engine) applyTemplate(templatePath, destPath string, data map[string]interface{}) error
func (e *Engine) renderTemplate(content string, data map[string]interface{}) ([]byte, error)
func (e *Engine) runPostSteps(projectPath string, opts InitOptions) error
```

### 3. Module Registry (`internal/modules/`)

**Responsibility:** Module management and application

**Files:**
- `registry.go` - Module registry and interface
- `postgres.go` - PostgreSQL module
- `mongo.go` - MongoDB module
- `openapi.go` - OpenAPI modules (gen/manual)
- `docker.go` - Docker module

**Module Interface:**
```go
type Module interface {
    Name() string
    Apply(ctx *Context) error
}
```

**Registry Pattern:**
```go
type Registry struct {
    modules map[string]map[string]Module
}

// Usage
registry.Register("db", "postgres", &PostgresModule{})
mod := registry.GetModule("db", "postgres")
```

**Module Implementation Example:**
```go
type PostgresModule struct{}

func (m *PostgresModule) Name() string {
    return "postgres"
}

func (m *PostgresModule) Apply(ctx *Context) error {
    // 1. Apply templates
    templatePath := "templates/db/postgres"
    if err := ApplyModuleTemplate(templatePath, ctx.ProjectPath, ctx.TemplateData); err != nil {
        return err
    }
    
    // 2. Additional logic (migrations, etc.)
    if ctx.Options.Migrations == "goose" {
        // Create migration files
    }
    
    return nil
}
```

### 4. Template System (`templates/`)

**Responsibility:** Embedded project templates

**Structure:**
```
templates/
├── base/                      # Base project (always applied)
│   ├── cmd/server/
│   │   └── main.go.tmpl       # Templated file
│   ├── internal/
│   │   ├── config/
│   │   ├── logger/
│   │   ├── http/
│   │   └── errors/
│   ├── .gitignore             # Static file
│   └── README.md.tmpl         # Templated file
│
├── db/
│   ├── postgres/              # PostgreSQL overlay
│   │   └── internal/db/postgres/
│   │       ├── postgres.go.tmpl
│   │       └── repository.go.tmpl
│   └── mongo/                 # MongoDB overlay
│       └── internal/db/mongo/
│           ├── mongo.go.tmpl
│           └── repository.go.tmpl
│
├── openapi/
│   ├── gen/                   # OpenAPI code generation
│   │   └── internal/api/
│   │       ├── generated/
│   │       └── handlers/
│   └── manual/                # OpenAPI manual mode
│       └── internal/api/handlers/
│
└── docker/
    ├── Dockerfile
    └── docker-compose.yml.tmpl
```

**Template Rules:**

1. **Files ending in `.tmpl`** are rendered with Go's `text/template`
2. **Other files** are copied as-is
3. **Templates have access to:**
   ```go
   type TemplateData map[string]interface{}{
       "ProjectName": "my-service",
       "ModulePath":  "github.com/user/my-service",
       "Router":      "chi",
       "Database":    "postgres",
       "OpenAPI":     "gen",
       "HasDocker":   true,
   }
   ```

4. **Template syntax:**
   ```go
   // Conditionals
   {{- if eq .Database "postgres"}}
   import "github.com/jackc/pgx/v5"
   {{- end}}
   
   // Variables
   package {{.ProjectName}}
   module {{.ModulePath}}
   
   // Loops (if needed)
   {{- range .Items}}
   - {{.}}
   {{- end}}
   ```

**Embedding:**
```go
//go:embed all:../../templates
var templatesFS embed.FS
```

## Data Flow

### Init Command Flow

```
User Command
    │
    ├─> gocrete init my-service --module github.com/user/my-service --db postgres
    │
    ▼
CLI Layer (cmd/init.go)
    │
    ├─> Parse flags
    ├─> Validate required flags
    ├─> Create InitOptions
    │
    ▼
Engine (engine/engine.go)
    │
    ├─> Validate options
    ├─> Create project directory
    │
    ├─> Create Context with template data
    │
    ├─> Apply base template
    │   │
    │   ├─> Walk embedded template files
    │   ├─> For each .tmpl file:
    │   │   ├─> Parse template
    │   │   ├─> Execute with data
    │   │   └─> Write to destination
    │   └─> For each other file:
    │       └─> Copy to destination
    │
    ├─> Get PostgresModule from registry
    │   │
    │   ▼
    │   PostgresModule (modules/postgres.go)
    │       │
    │       ├─> Apply templates/db/postgres
    │       └─> Create migration files
    │
    ├─> Run post-generation steps
    │   ├─> go mod init
    │   ├─> go mod tidy
    │   └─> go fmt
    │
    ▼
Success message to user
```

### Add Command Flow

```
User Command
    │
    ├─> gocrete add db --type postgres
    │
    ▼
CLI Layer (cmd/add.go)
    │
    ├─> Verify we're in a Go project
    ├─> Parse flags
    ├─> Create AddOptions
    │
    ▼
Engine (engine/engine.go)
    │
    ├─> Read go.mod for module path
    ├─> Create Context
    │
    ├─> Get module from registry
    │   │
    │   ▼
    │   Module.Apply(ctx)
    │
    ├─> go mod tidy
    │
    ▼
Success message to user
```

## Extensibility

### Adding a New Module

Example: Adding a Redis module

**1. Create module file:**
```go
// internal/modules/redis.go
package modules

type RedisModule struct{}

func (m *RedisModule) Name() string {
    return "redis"
}

func (m *RedisModule) Apply(ctx *Context) error {
    templatePath := "templates/cache/redis"
    return ApplyModuleTemplate(templatePath, ctx.ProjectPath, ctx.TemplateData)
}
```

**2. Register module:**
```go
// internal/modules/registry.go
func NewRegistry() *Registry {
    r := &Registry{modules: make(map[string]map[string]Module)}
    
    // Existing registrations...
    
    // Add Redis
    r.Register("cache", "redis", &RedisModule{})
    
    return r
}
```

**3. Create templates:**
```
templates/cache/redis/
└── internal/cache/redis/
    ├── redis.go.tmpl
    └── client.go.tmpl
```

**4. Update CLI:**
```go
// internal/cmd/init.go
var cache string
initCmd.Flags().StringVar(&cache, "cache", "none", "Cache (none|redis)")

// Add to InitOptions
type InitOptions struct {
    // ...
    Cache string
}
```

**5. Update engine:**
```go
// internal/engine/engine.go
func (e *Engine) InitProject(...) {
    // ...
    
    // Apply cache module
    if opts.Cache != "none" {
        mod := e.registry.GetModule("cache", opts.Cache)
        if err := mod.Apply(ctx); err != nil {
            return err
        }
    }
}
```

### Module Categories

Modules are organized by category:

- `db` - Database modules (postgres, mongo, mysql, etc.)
- `openapi` - OpenAPI modes (gen, manual)
- `docker` - Container configuration
- `cache` - Caching layers (redis, memcached)
- `auth` - Authentication (jwt, oauth, session)
- `queue` - Message queues (rabbitmq, kafka)
- `search` - Search engines (elasticsearch, typesense)
- `storage` - Object storage (s3, gcs)

### Template Best Practices

1. **Keep templates focused** - One responsibility per template
2. **Use conditionals sparingly** - Prefer separate templates for complex variations
3. **Provide defaults** - Templates should work with minimal data
4. **Document variables** - Comment expected template data
5. **Test templates** - Ensure they render correctly

### Validation Strategy

Validation happens at multiple levels:

```go
// CLI level - basic validation
if openapi == "gen" && specPath == "" {
    return fmt.Errorf("--spec required when using --openapi gen")
}

// Engine level - logical validation
func (e *Engine) validateInitOptions(opts InitOptions) error {
    // Validate router choices
    // Validate database choices
    // Check incompatibilities
}

// Module level - module-specific validation
func (m *Module) Validate(opts InitOptions) error {
    // Module-specific checks
}
```

## Design Decisions

### Why Not a Framework?

**Frameworks require:**
- Runtime dependency
- Version upgrades
- Breaking changes
- Learning curve
- Lock-in

**Generators provide:**
- Code that's yours
- No dependencies
- Immediate understanding
- Full control
- Easy customization

### Why Go Templates?

**Alternatives considered:**
- String replacement - Too fragile
- AST manipulation - Too complex
- External template engine - Additional dependency

**Why `text/template`:**
- Standard library
- Familiar to Go developers
- Powerful enough for our needs
- Deterministic
- Fast

### Why Embedded Templates?

**Alternatives:**
- External template directory - Deployment complexity
- Download from GitHub - Network dependency
- Bundled templates - Distribution complexity

**Why `go:embed`:**
- Single binary
- Version-locked templates
- Fast access
- No network required
- Easy to distribute

### Module Registry Pattern

**Benefits:**
- Decoupled modules
- Easy to add new modules
- Clear dependencies
- Testable in isolation
- Version module independently (future)

## Testing Strategy

### Unit Tests

Test individual components:
```go
func TestPostgresModule(t *testing.T) {
    ctx := &Context{...}
    mod := &PostgresModule{}
    err := mod.Apply(ctx)
    // Assert expected files exist
    // Assert content is correct
}
```

### Integration Tests

Test complete workflows:
```go
func TestInitProject(t *testing.T) {
    tmpDir := t.TempDir()
    engine := NewEngine()
    opts := InitOptions{...}
    
    err := engine.InitProject(tmpDir, opts)
    
    // Assert project structure
    // Assert go.mod exists and is valid
    // Assert can build: go build ./...
}
```

### E2E Tests

Test CLI commands:
```bash
# test/e2e_test.sh
#!/bin/bash
gocrete init test-project --module github.com/test/project
cd test-project
go build ./...
go test ./...
```

## Performance Considerations

### Template Rendering

- Templates are parsed once per file
- Rendering is O(n) where n = template size
- Embedding is done at compile-time
- No runtime template compilation

### File Operations

- Batch writes where possible
- Use buffered I/O
- Create directories once
- Avoid redundant stat calls

### Go Module Operations

- Run `go mod tidy` once at the end
- Use `-e` flag to continue on errors
- Cache module downloads

## Error Handling

### Strategy

```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to apply postgres template: %w", err)
}

// Provide actionable messages
return fmt.Errorf("directory %s already exists (use --force to overwrite)", name)

// Don't swallow errors
if err := doSomething(); err != nil {
    log.Error("non-critical error", "error", err)
    // Continue
}
```

### Error Types

- Validation errors - User input problems
- I/O errors - File system issues
- Template errors - Template syntax problems
- Execution errors - Command failures

## Future Architecture

### Planned Features

1. **Config File Support**
   ```yaml
   # gocrete.yaml
   modules:
     - postgres
     - docker
   router: chi
   ```

2. **Template Overrides**
   ```bash
   gocrete init my-service --templates ./my-templates
   ```

3. **Plugin System**
   ```bash
   gocrete plugin install github.com/user/gocrete-redis
   ```

4. **Interactive Mode**
   ```bash
   gocrete init --interactive
   # Wizard-style prompts
   ```

5. **Presets**
   ```bash
   gocrete init my-service --preset microservice
   ```

### Architectural Considerations

- Keep backward compatibility
- Maintain single-binary distribution
- No required external dependencies
- Templates versioned with binary
- Clear upgrade path

## Contributing to Architecture

When proposing architectural changes:

1. **Open an issue first** - Discuss the change
2. **Explain the problem** - What's not working?
3. **Propose solution** - Why is this the best approach?
4. **Consider impact** - Breaking changes? Migration path?
5. **Update docs** - Keep architecture docs current

## Resources

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Cobra Documentation](https://github.com/spf13/cobra)
- [text/template Documentation](https://pkg.go.dev/text/template)
- [embed Package](https://pkg.go.dev/embed)
