package mock

import (
	"context"
	"github.com/Kalinin-Andrey/dbmigrator/internal/infrastructure/db"
	"github.com/Kalinin-Andrey/dbmigrator/internal/test/fixture"
	"sort"

	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"
)

type Transaction struct {
}

type MigrationRepository struct {
	ExecutionLogs	[]MigrationRepositoryLog
}

type MigrationRepositoryLog struct {
	MethodName	string
	Params		map[string]interface{}
}

var _ migration.IRepository = (*MigrationRepository)(nil)
var _ migration.Transaction = (*Transaction)(nil)

func (t Transaction) Commit() error {
	return nil
}

func (t Transaction) Rollback() error {
	return nil
}

func NewMigrationRepository() *MigrationRepository {
	return &MigrationRepository{
		ExecutionLogs:	make([]MigrationRepositoryLog, 0, 10),
	}
}

func (r *MigrationRepository) Reset(logger api.Logger) {
	r.ExecutionLogs = make([]MigrationRepositoryLog, 0, 10)
}

func (r *MigrationRepository) SetLogger(logger api.Logger) {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"SetLogger",
		Params:		map[string]interface{}{
			"logger":	logger,
		},
	})
}

func (r *MigrationRepository) Query(ctx context.Context, offset, limit uint) ([]migration.MigrationLog, error) {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"Query",
		Params:		map[string]interface{}{
			"ctx":		ctx,
			"offset":	offset,
			"limit":	limit,
		},
	})

	if limit < 1 {
		limit = db.MaxLIstLimit
	}
	sl := fixture.MigrationsLogsList.GetSlice()
	sort.Sort(migration.MigrationsLogsSlice(sl))

	return sl, nil
}

func (r *MigrationRepository) QueryTx(ctx context.Context, t migration.Transaction, query *migration.QueryCondition, offset, limit uint) ([]migration.MigrationLog, error) {
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

	if limit < 1 {
		limit = db.MaxLIstLimit
	}

	mls := fixture.MigrationsLogsList
	if query != nil && query.Where != nil {
		tmls := FilterMigrationsLogsByStatus(*mls, query.Where.Status)
		mls = &tmls
	}

	if len(*mls) == 0 {
		return nil, api.ErrNotFound
	}

	sl := mls.GetSlice()
	sort.Sort(migration.MigrationsLogsSlice(sl))

	return sl, nil
}

func FilterMigrationsLogsByStatus(sourceMigrationsLogsList migration.MigrationsLogsList, status uint) migration.MigrationsLogsList {
	ml := make(migration.MigrationsLogsList)

	for id, i := range sourceMigrationsLogsList {
		if i.Status == status {
			ml[id] = i
		}
	}

	return ml
}

func (r *MigrationRepository) Last(ctx context.Context, query *migration.QueryCondition) (*migration.MigrationLog, error) {
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
		return nil, api.ErrNotFound
	}

	sl := mls.GetSlice()
	sort.Sort(migration.MigrationsLogsSlice(sl))

	return &sl[len(sl) - 1], nil
}

func (r *MigrationRepository) LastTx(ctx context.Context, t migration.Transaction, query *migration.QueryCondition) (*migration.MigrationLog, error) {
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
		return nil, api.ErrNotFound
	}

	sl := mls.GetSlice()
	sort.Sort(migration.MigrationsLogsSlice(sl))

	return &sl[len(sl) - 1], nil
}

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

func (r *MigrationRepository) ExecFunc(ctx context.Context, f api.MigrationFunc) (err error) {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"ExecFunc",
		Params:		map[string]interface{}{
			"ctx":		ctx,
			"f":		f,
		},
	})
	return nil
}

func (r *MigrationRepository) ExecFuncTx(ctx context.Context, t migration.Transaction, f api.MigrationFunc) (err error) {
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

func (r *MigrationRepository) BeginTx(ctx context.Context) (migration.Transaction, error) {
	r.ExecutionLogs = append(r.ExecutionLogs, MigrationRepositoryLog{
		MethodName:	"BeginTx",
		Params:		map[string]interface{}{
			"ctx":		ctx,
		},
	})
	return Transaction{}, nil
}

func (r *MigrationRepository) BatchCreateTx(ctx context.Context, t migration.Transaction, list migration.MigrationsLogsList) error {
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
			return api.ErrDuplicate
		}
		(*fixture.MigrationsLogsList)[id] = ml
	}

	return nil
}

func (r *MigrationRepository) BatchUpdateTx(ctx context.Context, t migration.Transaction, list migration.MigrationsLogsList) error {
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
			return api.ErrNotFound
		}
		(*fixture.MigrationsLogsList)[id] = ml
	}

	return nil
}





