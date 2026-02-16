package modules

import (
	"fmt"
)

type DockerModule struct{}

func (m *DockerModule) Name() string {
	return "docker"
}

func (m *DockerModule) Apply(ctx *Context) error {
	// Apply docker template
	templatePath := "files/docker"
	if err := ApplyModuleTemplate(templatePath, ctx.ProjectPath, ctx.TemplateData); err != nil {
		return fmt.Errorf("failed to apply docker template: %w", err)
	}

	return nil
}
