package modules

import (
	"fmt"
	"path/filepath"
)

type OpenAPIGenModule struct{}

func (m *OpenAPIGenModule) Name() string {
	return "openapi-gen"
}

func (m *OpenAPIGenModule) Apply(ctx *Context) error {
	// Apply openapi gen template
	templatePath := "files/openapi/gen"
	if err := ApplyModuleTemplate(templatePath, ctx.ProjectPath, ctx.TemplateData); err != nil {
		return fmt.Errorf("failed to apply openapi gen template: %w", err)
	}

	// Copy spec file if provided
	if ctx.Options.SpecPath != "" {
		// Create a simple example spec if file doesn't exist
		specContent := `openapi: 3.0.0
info:
  title: ` + ctx.Options.ProjectName + ` API
  version: 1.0.0
paths:
  /health:
    get:
      summary: Health check
      responses:
        '200':
          description: Service is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
  /api/v1/users:
    get:
      summary: List users
      responses:
        '200':
          description: List of users
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id:
                      type: integer
                    email:
                      type: string
`
		specPath := filepath.Join(ctx.ProjectPath, "api", "openapi.yaml")
		if err := WriteFile(specPath, specContent); err != nil {
			return err
		}
	}

	// Create Makefile for code generation
	makefileContent := `.PHONY: generate
generate:
	go generate ./...

.PHONY: api-gen
api-gen:
	oapi-codegen -package generated -generate types,chi-server,spec api/openapi.yaml > internal/api/generated/api.gen.go
`
	if err := WriteFile(filepath.Join(ctx.ProjectPath, "Makefile"), makefileContent); err != nil {
		return err
	}

	return nil
}

type OpenAPIManualModule struct{}

func (m *OpenAPIManualModule) Name() string {
	return "openapi-manual"
}

func (m *OpenAPIManualModule) Apply(ctx *Context) error {
	// Apply openapi manual template
	templatePath := "files/openapi/manual"
	if err := ApplyModuleTemplate(templatePath, ctx.ProjectPath, ctx.TemplateData); err != nil {
		return fmt.Errorf("failed to apply openapi manual template: %w", err)
	}

	return nil
}
