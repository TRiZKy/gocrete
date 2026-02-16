package modules

import (
	"fmt"
)

type MongoModule struct{}

func (m *MongoModule) Name() string {
	return "mongo"
}

func (m *MongoModule) Apply(ctx *Context) error {
	// Apply mongo template
	templatePath := "files/db/mongo"
	if err := ApplyModuleTemplate(templatePath, ctx.ProjectPath, ctx.TemplateData); err != nil {
		return fmt.Errorf("failed to apply mongo template: %w", err)
	}

	return nil
}
