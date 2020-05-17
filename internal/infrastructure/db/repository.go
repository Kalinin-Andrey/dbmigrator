package db

import (
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"
	"github.com/pkg/errors"
	"log"
	"os"

	"github.com/Kalinin-Andrey/dbmigrator/internal/pkg/dbx"

	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
)

// IRepository is an interface of repository
type IRepository interface {}

// repository persists albums in database
type repository struct {
	db                dbx.DBx
	logger            api.Logger
	//defaultConditions map[string]interface{}
}

// Limit is default limit
const Limit = 100

// GetRepository return a repository
func GetRepository(dbase dbx.DBx, logger api.Logger, entity string) (repo IRepository, err error) {
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

