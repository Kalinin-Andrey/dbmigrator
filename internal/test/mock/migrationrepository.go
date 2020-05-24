package mock

import (
	"context"
	"sort"

	"github.com/Kalinin-Andrey/dbmigrator/internal/pkg/apperror"

	"github.com/Kalinin-Andrey/dbmigrator/internal/app"
	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
	"github.com/Kalinin-Andrey/dbmigrator/internal/test/fixture"
)

// Transaction mock
type Transaction struct {
}

// MigrationRepository mock
type MigrationRepository struct {
	ExecutionLogs	[]MigrationRepositoryLog
}

// MigrationRepositoryLog struct
type MigrationRepositoryLog struct {
	MethodName	string
	Params		map[string]interface{}
}

var _ migration.IRepository = (*MigrationRepository)(nil)
var _ migration.Transaction = (*Transaction)(nil)

// Commit a transaction
func (t Transaction) Commit() error {
	return nil
}

// Rollback a transaction
func (t Transaction) Rollback() error {
	return nil
}

// NewMigrationRepository returns a new MigrationRepository mock
func NewMigrationRepository() *MigrationRepository {
	return &MigrationRepository{
		ExecutionLogs:	make([]MigrationRepositoryLog, 0, 10),
	}
}

// Reset MigrationRepository.ExecutionLogs
func (r *MigrationRepository) Reset(logger app.Logger) {
	r.ExecutionLogs = make([]MigrationRepositoryLog, 0, 10)
}

// SetLogger mock
func (r *MigrationRepository) SetLogger(logger app.Logger) {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"SetLogger",
		Params:		map[string]interface{}{
			"logger":	logger,
		},
	})
}

// Query mock
func (r *MigrationRepository) Query(ctx context.Context, offset, limit uint) ([]migration.Log, error) {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"Query",
		Params:		map[string]interface{}{
			"ctx":		ctx,
			"offset":	offset,
			"limit":	limit,
		},
	})

	/*if limit == 0 {
		limit = db.MaxLIstLimit
	}*/
	sl := fixture.MigrationsLogsList.Slice()
	sort.Sort(migration.LogsSlice(sl))

	return sl, nil
}

// QueryTx mock
func (r *MigrationRepository) QueryTx(ctx context.Context, t migration.Transaction, query *migration.QueryCondition, offset, limit uint) ([]migration.Log, error) {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"QueryTx",
		Params:		map[string]interface{}{
			"ctx":		ctx,
			"t":		t,
			"query":	query,
			"offset":	offset,
			"limit":	limit,
		},
	})

	/*if limit == 0 {
		limit = db.MaxLIstLimit
	}*/
	mls := fixture.MigrationsLogsList
	if query != nil && query.Where != nil {
		tmls := FilterMigrationsLogsByStatus(*mls, query.Where.Status)
		mls = &tmls
	}

	if len(*mls) == 0 {
		return nil, apperror.ErrNotFound
	}

	sl := mls.Slice()
	sort.Sort(migration.LogsSlice(sl))

	return sl, nil
}

// FilterMigrationsLogsByStatus mock
func FilterMigrationsLogsByStatus(sourceMigrationsLogsList migration.LogsList, status uint) migration.LogsList {
	ml := make(migration.LogsList)

	for id, i := range sourceMigrationsLogsList {
		if i.Status == status {
			ml[id] = i
		}
	}

	return ml
}

// Last mock
func (r *MigrationRepository) Last(ctx context.Context, query *migration.QueryCondition) (*migration.Log, error) {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"Last",
		Params:		map[string]interface{}{
			"ctx":		ctx,
			"query":	query,
		},
	})

	mls := fixture.MigrationsLogsList
	if query != nil && query.Where != nil {
		tmls := FilterMigrationsLogsByStatus(*mls, query.Where.Status)
		mls = &tmls
	}

	if len(*mls) == 0 {
		return nil, apperror.ErrNotFound
	}

	sl := mls.Slice()
	sort.Sort(migration.LogsSlice(sl))

	return &sl[len(sl) - 1], nil
}

// LastTx mock
func (r *MigrationRepository) LastTx(ctx context.Context, t migration.Transaction, query *migration.QueryCondition) (*migration.Log, error) {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"LastTx",
		Params:		map[string]interface{}{
			"ctx":		ctx,
			"t":		t,
			"query":	query,
		},
	})

	mls := fixture.MigrationsLogsList
	if query != nil && query.Where != nil {
		tmls := FilterMigrationsLogsByStatus(*mls, query.Where.Status)
		mls = &tmls
	}

	if len(*mls) == 0 {
		return nil, apperror.ErrNotFound
	}

	sl := mls.Slice()
	sort.Sort(migration.LogsSlice(sl))

	return &sl[len(sl) - 1], nil
}

// ExecSQL mock
func (r *MigrationRepository) ExecSQL(ctx context.Context, sql string) error {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"ExecSQL",
		Params:		map[string]interface{}{
			"ctx":		ctx,
			"sql":		sql,
		},
	})
	return nil
}

// ExecSQLTx mock
func (r *MigrationRepository) ExecSQLTx(ctx context.Context, t migration.Transaction, sql string) error {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"ExecSQLTx",
		Params:		map[string]interface{}{
			"ctx":		ctx,
			"t":		t,
			"sql":		sql,
		},
	})
	return nil
}

// ExecFunc mock
func (r *MigrationRepository) ExecFunc(ctx context.Context, f migration.Func) (err error) {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"ExecFunc",
		Params:		map[string]interface{}{
			"ctx":		ctx,
			"f":		f,
		},
	})
	return nil
}

// ExecFuncTx mock
func (r *MigrationRepository) ExecFuncTx(ctx context.Context, t migration.Transaction, f migration.Func) (err error) {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"ExecFuncTx",
		Params:		map[string]interface{}{
			"ctx":		ctx,
			"t":		t,
			"f":		f,
		},
	})
	return nil
}

// BeginTx mock
func (r *MigrationRepository) BeginTx(ctx context.Context) (migration.Transaction, error) {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"BeginTx",
		Params:		map[string]interface{}{
			"ctx":		ctx,
		},
	})
	return Transaction{}, nil
}

// BatchCreateTx mock
func (r *MigrationRepository) BatchCreateTx(ctx context.Context, t migration.Transaction, list migration.LogsList) error {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"BatchCreateTx",
		Params:		map[string]interface{}{
			"ctx":		ctx,
			"t":		t,
			"list":		list,
		},
	})

	for id, ml := range list {
		if _, ok := (*fixture.MigrationsLogsList)[id]; ok {
			return apperror.ErrDuplicate
		}
		(*fixture.MigrationsLogsList)[id] = ml
	}

	return nil
}

// BatchUpdateTx mock
func (r *MigrationRepository) BatchUpdateTx(ctx context.Context, t migration.Transaction, list migration.LogsList) error {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"BatchUpdateTx",
		Params:		map[string]interface{}{
			"ctx":		ctx,
			"t":		t,
			"list":		list,
		},
	})

	for id, ml := range list {
		if _, ok := (*fixture.MigrationsLogsList)[id]; !ok {
			return apperror.ErrNotFound
		}
		(*fixture.MigrationsLogsList)[id] = ml
	}

	return nil
}





