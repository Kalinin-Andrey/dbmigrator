package dbmigrator

import (
	"context"
	"github.com/Kalinin-Andrey/dbmigrator/internal/infrastructure/gomigration"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator/api"
)

const (
	// MainFileName const
	MainFileName	= "main.go"
	// actionUp const
	actionUp		= "up"
	// actionDown const
	actionDown		= "down"
	// actionRedo const
	actionRedo		= "redo"
)

// DBMigratorTool is DBMigrator as a tool
// nolint
type DBMigratorTool struct { // nolint
	*DBMigrator
}

// InitTool is the func for initialisation DBMigrator as a tool
func InitTool(ctx context.Context, config api.Configuration, logger api.Logger) (err error) {
	if err := Init(ctx, config, logger); err != nil {
		return err
	}
	dbMigrator, err = NewDBMigratorTool(dbMigrator.(*DBMigrator))
	if err != nil {
		return err
	}
	return nil
}

// NewDBMigratorTool returns new DBMigratorTool
func NewDBMigratorTool(m *DBMigrator) (*DBMigratorTool, error) {
	service, err := migration.NewServiceTool(m.domain.Migration.Service)
	if err != nil {
		return nil, err
	}
	m.domain.Migration.Service = service
	mt := &DBMigratorTool{m}
	return mt, nil
}

// Up migrations
func (m *DBMigratorTool) Up(quantity int) (err error) {
	return m.Exec(actionUp)
}

// Down migrations
func (m *DBMigratorTool) Down(quantity int) (err error) {
	return m.Exec(actionDown)
}

// Redo a one last migration
func (m *DBMigratorTool) Redo() (err error) {
	return m.Exec(actionRedo)
}

// Exec migrations
func (m *DBMigratorTool) Exec(action string) (err error) {
	dir := gomigration.Dir{
		Path: m.config.Dir,
	}
	if err = dir.Validate(); err != nil {
		return err
	}

	output, err := dir.Run(gomigration.Args{
		DSN:    m.config.DSN,
		Action: action,
	})
	m.logger.Print(output)
	return err
}

// Create a migration
func (m *DBMigratorTool) Create(p api.MigrationCreateParams) (err error) {
	if err = m.checkMainFile(); err != nil {
		return err
	}
	return m.DBMigrator.Create(p)
}

// checkMainFile checks if main fille exists
func (m *DBMigratorTool) checkMainFile() (err error) {
	filePath := filepath.Join(m.config.Dir, MainFileName)

	if _, err := os.Stat(filePath); err != nil {
		err = m.createMainFile(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

// createMainFile creates main file
func (m *DBMigratorTool) createMainFile(fileName string) (err error) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return errors.Wrapf(err, "Error while creating a main file")
	}
	defer f.Close()

	err = m.domain.Migration.Service.CreateMainFile(m.ctx, f)
	return api.AppErrorConv(err)
}


