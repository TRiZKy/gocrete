# Contributing to Gocrete

Thank you for considering contributing to Gocrete! This document provides guidelines for contributing to the project.

## Code of Conduct

Be respectful, inclusive, and constructive in all interactions.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in Issues
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Go version and OS
   - Generated project example (if applicable)

### Suggesting Features

1. Check existing issues and discussions
2. Create a new issue describing:
   - The problem you're trying to solve
   - Your proposed solution
   - Alternative solutions considered
   - Examples of how it would work

### Contributing Code

1. **Fork the repository**

2. **Create a feature branch**
   ```bash
   git checkout -b feature/my-feature
   ```

3. **Make your changes**
   - Follow Go best practices
   - Add tests for new functionality
   - Update documentation as needed
   - Ensure all tests pass: `go test ./...`
   - Format code: `go fmt ./...`

4. **Commit your changes**
   - Use clear, descriptive commit messages
   - Reference issues when applicable

5. **Push to your fork**
   ```bash
   git push origin feature/my-feature
   ```

6. **Create a Pull Request**
   - Provide a clear description of changes
   - Link to related issues
   - Explain the reasoning behind your approach

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/gocrete.git
cd gocrete

# Install dependencies
go mod download

# Build
go build -o gocrete cmd/gocrete/main.go

# Run tests
go test ./...

# Test the CLI
./gocrete init test-project --module github.com/test/project
```

## Project Structure

```
gocrete/
├── cmd/gocrete/         # CLI entry point
├── internal/
│   ├── cmd/             # Cobra commands
│   ├── engine/          # Generator engine
│   └── modules/         # Module implementations
└── templates/           # Embedded templates (go:embed)
```

## Adding New Modules

To add a new module (e.g., Redis):

1. **Create module implementation**
   ```go
   // internal/modules/redis.go
   package modules

   type RedisModule struct{}

   func (m *RedisModule) Name() string {
       return "redis"
   }

   func (m *RedisModule) Apply(ctx *Context) error {
       templatePath := "templates/redis"
       return ApplyModuleTemplate(templatePath, ctx.ProjectPath, ctx.TemplateData)
   }
   ```

2. **Register the module**
   ```go
   // internal/modules/registry.go
   func NewRegistry() *Registry {
       r := &Registry{...}
       // ... existing registrations
       r.Register("cache", "redis", &RedisModule{})
       return r
   }
   ```

3. **Create templates**
   ```
   templates/redis/
   └── internal/cache/redis/
       ├── redis.go.tmpl
       └── client.go.tmpl
   ```

4. **Update CLI flags** (if needed)

5. **Add tests**

6. **Update documentation**

## Testing Guidelines

- Write tests for new functionality
- Ensure tests are deterministic
- Use table-driven tests where appropriate
- Mock external dependencies
- Test both success and error cases

Example test:
```go
func TestPostgresModule(t *testing.T) {
    tests := []struct {
        name    string
        ctx     *Context
        wantErr bool
    }{
        {
            name: "valid context",
            ctx:  &Context{...},
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := &PostgresModule{}
            err := m.Apply(tt.ctx)
            if (err != nil) != tt.wantErr {
                t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Template Guidelines

Templates must:
- Use `.tmpl` extension for files requiring templating
- Use meaningful variable names in template data
- Include comments explaining complex logic
- Be idempotent when possible
- Follow Go project layout conventions

## Documentation Guidelines

- Update README.md for new features
- Add examples for new modules
- Keep documentation concise and clear
- Include code examples where helpful
- Update architecture docs if structure changes

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `go vet` before committing
- Keep functions focused and small
- Write clear error messages
- Use meaningful variable names

## Commit Message Format

```
type(scope): brief description

Longer description if needed, explaining:
- Why the change was made
- What problem it solves
- Any breaking changes

Fixes #123
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

## Pull Request Process

1. Update documentation
2. Add tests for new code
3. Ensure all tests pass
4. Update CHANGELOG if applicable
5. Request review from maintainers
6. Address review feedback
7. Squash commits if requested

## Questions?

- Open an issue for questions
- Join discussions for broader topics
- Reach out to maintainers directly

## License

By contributing to Gocrete, you agree that your contributions will be licensed under the MIT License.
