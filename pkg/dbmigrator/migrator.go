package dbmigrator

import (
	"context"
	"fmt"
	"github.com/Kalinin-Andrey/dbmigrator/internal/pkg/dbx"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator/api"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"

	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
	dbrep "github.com/Kalinin-Andrey/dbmigrator/internal/infrastructure/db"
)

const Dialect string = "postgres"


type DBMigrator struct {
	ctx		context.Context
	config	api.Configuration
	logger	api.Logger
	domain	Domain
	ms		migration.MigrationsList
}

var ms		= make(migration.MigrationsList)
var errs	= make([]error, 0)


// Domain is a Domain Layer Entry Point
type Domain struct {
	Migration struct {
		Repository migration.IRepository
		Service    migration.IService
	}
}

var dbMigrator *DBMigrator


func Add(item migration.Migration) {

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


// Init initialises DBMigrator instance
func Init(ctx context.Context, config api.Configuration, logger api.Logger) error {
	if len(errs) > 0 {
		return errors.Errorf("DBMigrator.Init errors: \n%v", errs)
	}

	if dbMigrator == nil {


		dbx, err := dbx.New(*config.DBxConf(), nil)
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

		dbMigrator, err = NewSQLMigrator(ctx, config, logger, repository, ms)
		if err != nil {
			return err
		}
	}

	return nil
}


func NewSQLMigrator(ctx context.Context, config api.Configuration, logger api.Logger, repository migration.IRepository, ms migration.MigrationsList) (*DBMigrator, error) {
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
		return nil, api.AppErrorConv(err)
	}

	return &DBMigrator{
		ctx:    ctx,
		config: config,
		logger: logger,
		domain: domain,
		ms:		ms,
	}, nil
}


func Up(quantity int) (err error) {
	if dbMigrator == nil {
		return api.ErrNotInitialised
	}
	return dbMigrator.Up(quantity)
}


func (m *DBMigrator) Up(quantity int) (err error) {
	err = m.domain.Migration.Service.Up(m.ctx, m.ms, quantity)
	return api.AppErrorConv(err)
}


func Down(quantity int) (err error) {
	if dbMigrator == nil {
		return api.ErrNotInitialised
	}
	return dbMigrator.Down(quantity)
}


func (m *DBMigrator) Down(quantity int) (err error) {
	err = m.domain.Migration.Service.Down(m.ctx, m.ms, quantity)
	return api.AppErrorConv(err)
}


func Redo() (err error) {
	if dbMigrator == nil {
		return api.ErrNotInitialised
	}
	return dbMigrator.Redo()
}


func (m *DBMigrator) Redo() (err error) {
	err = m.domain.Migration.Service.Redo(m.ctx, m.ms)
	return api.AppErrorConv(err)
}


func Status() ([]migration.MigrationLog, error) {
	if dbMigrator == nil {
		return nil, api.ErrNotInitialised
	}
	return dbMigrator.Status()
}


func (m *DBMigrator) Status() ([]migration.MigrationLog, error) {
	list, err := m.domain.Migration.Service.List(m.ctx)
	err = api.AppErrorConv(err)
	if err != nil && errors.Is(err, api.ErrNotFound) {
		err = nil
	}
	return list, err
}


func DBVersion() (uint, error) {
	if dbMigrator == nil {
		return 0, api.ErrNotInitialised
	}
	return dbMigrator.DBVersion()
}


func (m *DBMigrator) DBVersion() (uint, error) {
	lm, err := m.domain.Migration.Service.Last(m.ctx)
	err = api.AppErrorConv(err)
	if err != nil {
		if errors.Is(err, api.ErrNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return lm.ID, nil
}


func Create(p api.MigrationCreateParams) (err error) {
	if dbMigrator == nil {
		return api.ErrNotInitialised
	}
	return dbMigrator.Create(p)
}


func (m *DBMigrator) Create(p api.MigrationCreateParams) (err error) {
	cp := p.CoreParams()
	if err = cp.Validate(); err != nil {
		return errors.Wrapf(err, "Invalid create params")
	}

	if _, ok := m.ms[cp.ID]; ok {
		return errors.Wrapf(api.ErrBadRequest, "Migration #%v already exists", cp.ID)
	}
	fileName := fmt.Sprintf("%03d", cp.ID) + "_" + cp.Name +".go"
	fileName = filepath.Join(m.config.Dir, fileName)

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return errors.Wrapf(err, "Error while creating a file")
	}
	defer f.Close()

	err = m.domain.Migration.Service.Create(m.ctx, f, *cp)
	return api.AppErrorConv(err)
}

