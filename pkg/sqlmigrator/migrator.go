package sqlmigrator

import (
	"context"
	"github.com/Kalinin-Andrey/dbmigrator/internal/pkg/dbx"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"
	"github.com/pkg/errors"
	"log"
	"os"

	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
	dbrep "github.com/Kalinin-Andrey/dbmigrator/internal/infrastructure/db"
)

const Dialect string = "postgres"


type SQLMigrator struct {
	ctx		context.Context
	config	api.Configuration
	logger	api.Logger
	domain	Domain
	ms		api.MigrationsList
	err		error
}

const DefaultQuantityDown = 1

var ms		api.MigrationsList	= make(api.MigrationsList, 0)
var errs	[]error				= make([]error, 0)


// Domain is a Domain Layer Entry Point
type Domain struct {
	Migration struct {
		Repository migration.IRepository
		Service    migration.IService
	}
}

var sqlMigrator *SQLMigrator


func Add(item api.Migration) {

	if _, ok := ms[item.ID]; ok {
		errs = append(errs, errors.Wrapf(api.ErrDuplicate, "Duplicate migration ID: %v", item.ID))
		return
	}

	if err := item.Validate(); err != nil {
		errs = append(errs, errors.Wrapf(err, "Invalid migration #%v", item.ID))
		return
	}

	ms[item.ID] = item
}


// Init initialises SQLMigrator instance
func Init(ctx context.Context, config api.Configuration, logger api.Logger) error {
	if len(errs) > 0 {
		return errors.Errorf("SQLMigrator.Init errors: \n%v", errs)
	}

	if sqlMigrator == nil {


		dbx, err := dbx.New(config, nil)
		if err != nil {
			return err
		}

		rep, err := dbrep.GetRepository(dbx, nil, migration.TableName)
		if err != nil {
			return err
		}

		repository, ok := rep.(migration.IRepository)
		if !ok {
			return errors.Errorf("Can not cast DB repository for entity %q to %v.IRepository. Repo: %v", migration.TableName, migration.TableName, rep)
		}

		sqlMigrator, err = NewSQLMigrator(ctx, config, nil, repository, ms)
		if err != nil {
			return err
		}
	}

	return nil
}


func NewSQLMigrator(ctx context.Context, config api.Configuration, logger api.Logger, repository migration.IRepository, ms api.MigrationsList) (*SQLMigrator, error) {
	config.Dialect = Dialect
	if logger == nil {
		logger = log.New(os.Stdout, "sqlmigrator", log.LstdFlags)
	}
	repository.SetLogger(logger)

	domain := Domain{}
	domain.Migration.Repository	= repository
	domain.Migration.Service	= migration.NewService(domain.Migration.Repository, logger)

	err := domain.Migration.Service.CreateTable(ctx)
	if err != nil {
		return nil, err
	}

	return &SQLMigrator{
		ctx:    ctx,
		config: config,
		logger: logger,
		domain: domain,
		ms:		ms,
	}, nil
}


func Up(quantity int) (err error) {
	if sqlMigrator == nil {
		return api.ErrNotInitialised
	}
	return sqlMigrator.Up(quantity)
}


func (m *SQLMigrator) Up(quantity int) (err error) {
	return m.domain.Migration.Service.Up(m.ctx, m.ms, quantity)
}


func Down(quantity int) (err error) {
	if sqlMigrator == nil {
		return api.ErrNotInitialised
	}
	return sqlMigrator.Down(quantity)
}


func (m *SQLMigrator) Down(quantity int) (err error) {
	return m.domain.Migration.Service.Down(m.ctx, m.ms, quantity)
}


func Redo() (err error) {
	if sqlMigrator == nil {
		return api.ErrNotInitialised
	}
	return sqlMigrator.Redo()
}


func (m *SQLMigrator) Redo() (err error) {
	return m.domain.Migration.Service.Redo(m.ctx, m.ms)
}


func Status() ([]migration.MigrationLog, error) {
	if sqlMigrator == nil {
		return nil, api.ErrNotInitialised
	}
	return sqlMigrator.Status()
}


func (m *SQLMigrator) Status() ([]migration.MigrationLog, error) {
	list, err := m.domain.Migration.Service.List(m.ctx)
	if err != nil && errors.Is(err, api.ErrNotFound) {
		err = nil
	}
	return list, err
}


func DBVersion() (uint, error) {
	if sqlMigrator == nil {
		return 0, api.ErrNotInitialised
	}
	return sqlMigrator.DBVersion()
}


func (m *SQLMigrator) DBVersion() (uint, error) {
	lm, err := m.domain.Migration.Service.Last(m.ctx)
	if err != nil {
		if errors.Is(err, api.ErrNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return lm.ID, nil
}


func Create(p api.MigrationCreateParams) (err error) {
	if sqlMigrator == nil {
		return api.ErrNotInitialised
	}
	return sqlMigrator.Create(p)
}


func (m *SQLMigrator) Create(p api.MigrationCreateParams) (err error) {
	return m.domain.Migration.Service.Create(m.ctx, m.ms, m.config.Dir, p)
}

