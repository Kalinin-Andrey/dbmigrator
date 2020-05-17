package migration

import (
	"context"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"
)

// IRepository encapsulates the logic to access albums from the data source.
type IRepository interface {
	SetLogger(logger api.Logger)
	// Get returns an entity with the specified ID.
	//Get(ctx context.Context, id uint) (*MigrationLog, error)
	// Count returns the number of entities.
	//Count(ctx context.Context) (uint, error)
	// Query returns the list of entities with the given offset and limit.
	Query(ctx context.Context, offset, limit uint) ([]MigrationLog, error)
	// QueryTx returns the list of entities with the given offset and limit.
	QueryTx(ctx context.Context, t Transaction, query *QueryCondition, offset, limit uint) ([]MigrationLog, error)
	// Last retrieves a last record with the specified query condition and limit 1 from the database.
	Last(ctx context.Context, query *QueryCondition) (*MigrationLog, error)
	// LastTx retrieves a last record with the specified query condition and limit 1 from the database.
	LastTx(ctx context.Context, t Transaction, query *QueryCondition) (*MigrationLog, error)
	// Create saves a new entity in the storage.
	//Create(ctx context.Context, entity *MigrationLog) error
	// Update updates an entity with given ID in the storage.
	//Update(ctx context.Context, entity *MigrationLog) error
	// Delete removes an entity with given ID from the storage.
	//Delete(ctx context.Context, id uint) error
	// ExecSQL executes an user's plain sql
	ExecSQL(ctx context.Context, sql string) error
	// ExecSQLTx executes an user's plain sql
	ExecSQLTx(ctx context.Context, t Transaction, sql string) error
	// ExecFunc executes an user's func
	ExecFunc(ctx context.Context, f api.MigrationFunc) (err error)
	// ExecFuncTx executes an user's func
	ExecFuncTx(ctx context.Context, t Transaction, f api.MigrationFunc) (err error)
	// BeginTx begins a transaction
	BeginTx(ctx context.Context) (Transaction, error)
	BatchCreateTx(ctx context.Context, t Transaction, list MigrationsLogsList) error
	BatchUpdateTx(ctx context.Context, t Transaction, list MigrationsLogsList) error
}

type Transaction interface {
	Commit() error
	Rollback() error
}

