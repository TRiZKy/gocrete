package modules

import (
	"fmt"
	"path/filepath"
)

type PostgresModule struct{}

func (m *PostgresModule) Name() string {
	return "postgres"
}

func (m *PostgresModule) Apply(ctx *Context) error {
	// Apply postgres template
	templatePath := "files/db/postgres"
	if err := ApplyModuleTemplate(templatePath, ctx.ProjectPath, ctx.TemplateData); err != nil {
		return fmt.Errorf("failed to apply postgres template: %w", err)
	}

	// Create migrations directory if goose is enabled
	if ctx.Options.Migrations == "goose" {
		migrationsDir := filepath.Join(ctx.ProjectPath, "migrations")
		if err := WriteFile(filepath.Join(migrationsDir, ".gitkeep"), ""); err != nil {
			return err
		}

		// Create initial migration
		migration := `-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
`
		if err := WriteFile(filepath.Join(migrationsDir, "00001_initial.sql"), migration); err != nil {
			return err
		}
	}

	return nil
}
