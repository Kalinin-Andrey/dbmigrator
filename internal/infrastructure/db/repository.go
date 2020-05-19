package db

import (
	"context"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/dbmigrator/internal/pkg/dbx"

	"github.com/Kalinin-Andrey/dbmigrator/internal/app"
	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
)

// IRepository is an interface of repository
type IRepository interface {
	SetLogger(logger app.Logger)
	Query(ctx context.Context, offset, limit uint) ([]migration.MigrationLog, error)
}

// repository persists albums in database
type repository struct {
	db                dbx.DBx
	logger            app.Logger
	//defaultConditions map[string]interface{}
}

// MaxLIstLimit const
const MaxLIstLimit = 1000

// GetRepository return a repository
func GetRepository(dbase dbx.DBx, logger app.Logger, entity string) (repo IRepository, err error) {
	if logger == nil {
		logger = log.New(os.Stdout, "sqlmigrator", log.LstdFlags)
	}
	r := &repository{
		db:     dbase,
		logger: logger,
	}

	switch entity {
	case migration.TableName:
		repo, err = NewMigrationRepository(r)
	default:
		err = errors.Errorf("Repository for entity %q not found", entity)
	}
	return repo, err
}

