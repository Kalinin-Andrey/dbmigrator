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
	NewEntity() *Log
	// Get returns an entity with given ID
	//Get(ctx context.Context, id uint) (*Log, error)
	//First(ctx context.Context, entity *Event) (*Event, error)
	// Query returns a list with pagination
	Query(ctx context.Context, offset, limit uint) ([]Log, error)
	// List entity
	List(ctx context.Context) ([]Log, error)
	//Count(ctx context.Context) (uint, error)
	// Create entity
	//Create(ctx context.Context, entity *Log) error
	// Update entity
	//Update(ctx context.Context, entity *Log) error
	// Delete entity
	//Delete(ctx context.Context, id uint) error
	// CreateTable creates table for migration
	CreateTable(ctx context.Context) error
	// Up a list of migrations
	Up(ctx context.Context, ms MigrationsList, quantity int) error
	// Down a list of migrations
	Down(ctx context.Context, ms MigrationsList, quantity int) error
	// Redo a last migration
	Redo(ctx context.Context, ms MigrationsList) error
	// Last returns a last Log
	Last(ctx context.Context) (*Log, error)
	// Create creates a file for migration
	Create(ctx context.Context, wr io.Writer, p CreateParams) (err error)
	// Create creates a main file for migrations execution
	CreateMainFile(ctx context.Context, wr io.Writer) (err error)
}

// Service stgruct
type Service struct {
	repo   IRepository
	logger app.Logger
}

var _ IService = (*Service)(nil)

// DefaultDownQuantity const
const DefaultDownQuantity = 1

// NewService creates a new Service.
func NewService(repo IRepository, logger app.Logger) *Service {
	s := &Service{repo, logger}
	return s
}

// NewEntity returns a new empty entity
func (s Service) NewEntity() *Log {
	return &Log{}
}

// Get returns the entity with the specified ID.
/*func (s Service) Get(ctx context.Context, id uint) (*Log, error) {
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
func (s Service) Query(ctx context.Context, offset, limit uint) ([]Log, error) {
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
func (s Service) List(ctx context.Context) ([]Log, error) {
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
/*func (s Service) Create(ctx context.Context, entity *Log) error {
	return s.repo.Create(ctx, entity)
}

// Update entity
func (s Service) Update(ctx context.Context, entity *Log) error {
	return s.repo.Update(ctx, entity)
}

// Delete entity
func (s Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}*/

// CreateTable creates table for migration
func (s Service) CreateTable(ctx context.Context) error {
	return s.repo.ExecSQL(ctx, SQLCreateTable)
}

// Last returns a last Log
func (s Service) Last(ctx context.Context) (*Log, error) {
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

// Up a list of migrations
func (s Service) Up(ctx context.Context, ms MigrationsList, quantity int) error {
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
	ids			:= migrations.IDs()
	sort.Ints(ids)

	if quantity < 1 {
		quantity = len(ids)
	}

	if len(ids) < quantity {
		quantity = len(ids)
	}

	if quantity == 0 {
		return apperror.ErrNotFound
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

// Down a list of migrations
func (s Service) Down(ctx context.Context, ms MigrationsList, quantity int) error {
	t, err := s.repo.BeginTx(ctx)
	if err != nil {
		return errors.Wrapf(err, "migration.Service.Down: transaction begin error")
	}

	list, err := s.repo.QueryTx(ctx, t, nil, 0, 0)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return errors.Wrapf(apperror.ErrInternal, "migration.Service.Down: get list logs of migrations error: %v", err)
	}

	gl			:= GroupLogsByStatus(list)
	migrations	:= MigrationsListFilterExistsByKeys(ms, gl[StatusApplied])
	ids			:= migrations.IDs()
	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	if quantity < 1 {
		quantity = DefaultDownQuantity
	}

	if len(ids) < quantity {
		quantity = len(ids)
	}

	if quantity == 0 {
		return apperror.ErrNotFound
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
func (s Service) Redo(ctx context.Context, ms MigrationsList) error {
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


func (s Service) upProceed(ctx context.Context, ms MigrationsList, ids []int) (appliedMigrationsLogs LogsList, idErr uint, err error) {
	appliedMigrationsLogs = make(LogsList, len(ids))

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


func (s Service) downProceed(ctx context.Context, ms MigrationsList, ids []int) (downMigrationsLogs LogsList, idErr uint, err error) {
	downMigrationsLogs = make(LogsList, len(ids))

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



func (s Service) actionExec(ctx context.Context, in interface{}) (err error) {

	switch i := in.(type) {
	case string:
		err = s.repo.ExecSQL(ctx, i)
	case Func:
		err = s.repo.ExecFunc(ctx, i)
	default:
		err = apperror.ErrUndefinedTypeOfAction
	}

	return err
}


func (s Service) actionExecTx(ctx context.Context, t Transaction, in interface{}) (err error) {

	switch i := in.(type) {
	case string:
		err = s.repo.ExecSQLTx(ctx, t, i)
	case Func:
		err = s.repo.ExecFuncTx(ctx, t, i)
	default:
		err = apperror.ErrUndefinedTypeOfAction
	}

	return err
}

// Create creates a file for migration
func (s Service) Create(ctx context.Context, wr io.Writer, p CreateParams) (err error) {
	if err = p.Validate(); err != nil {
		return errors.Wrapf(err, "Invalid create params")
	}

	return s.getTemplate().ExecuteTemplate(wr, "tpl", p)
}

func (s Service) getTemplate() (*template.Template) {
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

// CreateMainFile here is dummy. Redefined in ServiceTool.
func (s Service) CreateMainFile(ctx context.Context, wr io.Writer) (err error) {
	return nil
}


