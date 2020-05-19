package migration

import (
	"context"
	"io"
	"sort"
	"text/template"

	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/dbmigrator/internal/pkg/apperror"

	"github.com/Kalinin-Andrey/dbmigrator/internal/app"
)

// IService encapsulates usecase logic for event.
type IService interface {
	// NewEntity returns new empty entity
	NewEntity() *MigrationLog
	// Get returns an entity with given ID
	//Get(ctx context.Context, id uint) (*MigrationLog, error)
	//First(ctx context.Context, entity *Event) (*Event, error)
	// Query returns a list with pagination
	Query(ctx context.Context, offset, limit uint) ([]MigrationLog, error)
	// List entity
	List(ctx context.Context) ([]MigrationLog, error)
	//Count(ctx context.Context) (uint, error)
	// Create entity
	//Create(ctx context.Context, entity *MigrationLog) error
	// Update entity
	//Update(ctx context.Context, entity *MigrationLog) error
	// Delete entity
	//Delete(ctx context.Context, id uint) error
	// CreateTable creates table for migration
	CreateTable(ctx context.Context) error
	// Up migration
	Up(ctx context.Context, ms MigrationsList, quantity int) error
	// Down migration
	Down(ctx context.Context, ms MigrationsList, quantity int) error
	// Redo a last migration
	Redo(ctx context.Context, ms MigrationsList) error
	Last(ctx context.Context) (*MigrationLog, error)
	Create(ctx context.Context, wr io.Writer, p MigrationCreateParams) (err error)
}

type service struct {
	//Domain     Domain
	repo   IRepository
	logger app.Logger
}

const(
	DefaultDownQuantity = 1
)

// NewService creates a new service.
func NewService(repo IRepository, logger app.Logger) IService {
	s := &service{repo, logger}
	return s
}

// NewEntity returns a new empty entity
func (s service) NewEntity() *MigrationLog {
	return &MigrationLog{}
}

// Get returns the entity with the specified ID.
/*func (s service) Get(ctx context.Context, id uint) (*MigrationLog, error) {
	entity, err := s.repo.Get(ctx, id)
	if err != nil {
		if err == apperror.ErrNotFound {
			return nil, err
		}
		return nil, errors.Wrapf(err, "Can not get a event by id: %v", id)
	}
	return entity, nil
}*/

// Query returns the items with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit uint) ([]MigrationLog, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
		return nil, errors.Wrapf(err, "error has occurred")
	}
	return items, nil
}


// List returns the items list.
func (s service) List(ctx context.Context) ([]MigrationLog, error) {
	items, err := s.repo.Query(ctx, 0, 0)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
		return nil, errors.Wrapf(err, "error has occurred")
	}
	return items, nil
}

// Create entity
/*func (s service) Create(ctx context.Context, entity *MigrationLog) error {
	return s.repo.Create(ctx, entity)
}

// Update entity
func (s service) Update(ctx context.Context, entity *MigrationLog) error {
	return s.repo.Update(ctx, entity)
}

// Delete entity
func (s service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}*/

// CreateTable creates table for migration
func (s service) CreateTable(ctx context.Context) error {
	return s.repo.ExecSQL(ctx, SQLCreateTable)
}

func (s service) Last(ctx context.Context) (*MigrationLog, error) {
	mLog, err := s.repo.Last(ctx, &QueryCondition{
		Where:	&WhereCondition{
			Status:	StatusApplied,
		},
	})

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
		return nil, errors.Wrapf(apperror.ErrInternal, "migration.Service.Last error: %v", err)
	}
	return mLog, nil
}

// Up list of migrations
func (s service) Up(ctx context.Context, ms MigrationsList, quantity int) error {
	t, err := s.repo.BeginTx(ctx)
	if err != nil {
		return errors.Wrapf(err, "migration.Service.Up: transaction begin error")
	}

	list, err := s.repo.QueryTx(ctx, t, nil, 0, 0)
	if err != nil {
		return errors.Wrapf(apperror.ErrInternal, "migration.Service.Down: get list logs of migrations error: %v", err)
	}

	gl			:= GroupLogsByStatus(list)
	migrations	:= MigrationsListFilterExceptByKeys(ms, gl[StatusApplied])
	ids			:= migrations.GetIDs()
	sort.Ints(ids)

	if quantity < 1 {
		quantity = len(ids)
	}

	appliedMigrationsLogs, idErr, err := s.upProceed(ctx, migrations, ids[:quantity])

	migrationsLogsForUpdate	:= MigrationsLogsFilterExistsByKeys(appliedMigrationsLogs, gl[StatusNotApplied])
	migrationsLogsForCreate := MigrationsLogsFilterExceptByKeys(appliedMigrationsLogs, gl[StatusNotApplied])

	if err != nil {
		if mLog, ok := gl[StatusNotApplied][idErr]; ok {
			mLog.Status = StatusError
			migrationsLogsForUpdate[idErr] = mLog
		} else {
			mLog = *migrations[idErr].Log(StatusError)
			migrationsLogsForCreate[idErr] = mLog
		}
	}
	err = s.repo.BatchUpdateTx(ctx, t, migrationsLogsForUpdate)
	if err != nil {
		err = t.Rollback()
		if err != nil {
			return errors.Wrapf(err, "migration.Service.Up: transaction rollback error")
		}
		return errors.Wrapf(err, "migration.Service.Up: batch update error")
	}

	err = s.repo.BatchCreateTx(ctx, t, migrationsLogsForCreate)
	if err != nil {
		err = t.Rollback()
		if err != nil {
			return errors.Wrapf(err, "migration.Service.Up: transaction rollback error")
		}
		return errors.Wrapf(err, "migration.Service.Up: batch create error")
	}

	err = t.Commit()
	if err != nil {
		return errors.Wrapf(err, "migration.Service.Up: transaction commit error")
	}

	return nil
}

// Down list of migrations
func (s service) Down(ctx context.Context, ms MigrationsList, quantity int) error {
	t, err := s.repo.BeginTx(ctx)
	if err != nil {
		return errors.Wrapf(err, "migration.Service.Down: transaction begin error")
	}

	list, err := s.repo.QueryTx(ctx, t, nil, 0, 0)
	if err != nil {
		return errors.Wrapf(apperror.ErrInternal, "migration.Service.Down: get list logs of migrations error: %v", err)
	}

	gl			:= GroupLogsByStatus(list)
	migrations	:= MigrationsListFilterExistsByKeys(ms, gl[StatusApplied])
	ids			:= migrations.GetIDs()
	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	if quantity < 1 {
		quantity = DefaultDownQuantity
	}

	migrationsLogsForUpdate, _, er := s.downProceed(ctx, migrations, ids[:quantity])

	err = s.repo.BatchUpdateTx(ctx, t, migrationsLogsForUpdate)
	if err != nil {
		err = t.Rollback()
		if err != nil {
			return errors.Wrapf(err, "migration.Service.Down: transaction rollback error")
		}
		return errors.Wrapf(err, "migration.Service.Down: batch update error")
	}

	err = t.Commit()
	if err != nil {
		return errors.Wrapf(err, "migration.Service.Down: transaction commit error")
	}

	return er
}

// Redo a last migration
func (s service) Redo(ctx context.Context, ms MigrationsList) error {
	t, err := s.repo.BeginTx(ctx)
	if err != nil {
		return errors.Wrapf(err, "migration.Service.Redo: transaction begin error")
	}

	mLog, err := s.repo.LastTx(ctx, t, &QueryCondition{
		Where:	&WhereCondition{
			Status:	StatusApplied,
		},
	})
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return errors.Wrapf(apperror.ErrInternal, "migration.Service.Redo: get last migration error: %v", err)
	}

	m, ok := ms[mLog.ID]
	if !ok {
		return errors.Wrapf(apperror.ErrNotFound, "migration.Service.Redo: can not find last migration #%v", mLog.ID)
	}

	err = s.actionExecTx(ctx, t, m.Down)
	if err != nil {
		err = t.Rollback()
		if err != nil {
			return errors.Wrapf(err, "migration.Service.Redo: transaction rollback error")
		}
		s.logger.Print("down #", m.ID, " - error: ", err)
		return errors.Wrapf(err, "migration.Service.Redo: error on down a migration #%v", mLog.ID)
	}
	s.logger.Print("down #", m.ID, " - done")

	err = s.actionExecTx(ctx, t, m.Up)
	if err != nil {
		err = t.Rollback()
		if err != nil {
			return errors.Wrapf(err, "migration.Service.Redo: transaction rollback error")
		}
		s.logger.Print("up #", m.ID, " - error: ", err)
		return errors.Wrapf(err, "migration.Service.Redo: error on up a migration #%v", mLog.ID)
	}
	s.logger.Print("up #", m.ID, " - done")

	err = t.Commit()
	if err != nil {
		return errors.Wrapf(err, "migration.Service.Redo: transaction commit error")
	}

	return nil
}


func (s service) upProceed(ctx context.Context, ms MigrationsList, ids []int) (appliedMigrationsLogs MigrationsLogsList, idErr uint, err error) {
	appliedMigrationsLogs = make(MigrationsLogsList, len(ids))

	for _, i := range ids {
		id := uint(i)
		err = s.actionExec(ctx, ms[id].Up)
		if err != nil {
			s.logger.Print("up #", id, " - error: ", err)
			return appliedMigrationsLogs, id, errors.Wrapf(err, "up error on migration #%v", id)
		}
		s.logger.Print("up #", id, " - done")
		appliedMigrationsLogs[id] = *ms[id].Log(StatusApplied)
	}
	return appliedMigrationsLogs, 0, nil
}


func (s service) downProceed(ctx context.Context, ms MigrationsList, ids []int) (downMigrationsLogs MigrationsLogsList, idErr uint, err error) {
	downMigrationsLogs = make(MigrationsLogsList, len(ids))

	for _, i := range ids {
		id := uint(i)
		err = s.actionExec(ctx, ms[id].Down)
		if err != nil {
			s.logger.Print("down #", id, " - error: ", err)
			return downMigrationsLogs, id, errors.Wrapf(err, "up error on migration #%v", id)
		}
		s.logger.Print("down #", id, " - done")
		downMigrationsLogs[id] = *ms[id].Log(StatusNotApplied)
	}
	return downMigrationsLogs, 0, nil
}



func (s service) actionExec(ctx context.Context, in interface{}) (err error) {

	switch i := in.(type) {
	case string:
		err = s.repo.ExecSQL(ctx, i)
	case MigrationFunc:
		err = s.repo.ExecFunc(ctx, i)
	default:
		err = apperror.ErrUndefinedTypeOfAction
	}

	return err
}


func (s service) actionExecTx(ctx context.Context, t Transaction, in interface{}) (err error) {

	switch i := in.(type) {
	case string:
		err = s.repo.ExecSQLTx(ctx, t, i)
	case MigrationFunc:
		err = s.repo.ExecFuncTx(ctx, t, i)
	default:
		err = apperror.ErrUndefinedTypeOfAction
	}

	return err
}

func (s service) Create(ctx context.Context, wr io.Writer, p MigrationCreateParams) (err error) {
	if err = p.Validate(); err != nil {
		return errors.Wrapf(err, "Invalid create params")
	}

	return getTemplate().ExecuteTemplate(wr, "tpl", p)
}

func getTemplate() (*template.Template) {
	return template.Must(template.New("tpl").Parse(`
package migration

import (
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"{{if (eq .Type "go")}}
	"github.com/jmoiron/sqlx"{{end}}
)

func init() {
	sqlmigrator.Add(api.Migration{
		ID:		{{.ID}},
		Name:	"{{.Name}}",
		Up:		{{if (eq .Type "sql")}}""{{else}}api.MigrationFunc(func(tx *sqlx.Tx) error {
			_, err := tx.Exec("CREATE TABLE IF NOT EXISTS public.test01(id int4)")	// for example
			return err
		}){{end}},
		Down:	{{if (eq .Type "sql")}}""{{else}}api.MigrationFunc(func(tx *sqlx.Tx) error {
			_, err := tx.Exec("DROP TABLE public.test01")							// for example
			return err
		}){{end}},
	})
}

`))
}



