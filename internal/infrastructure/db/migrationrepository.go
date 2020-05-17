package db

import (
	"context"
	"database/sql"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"
	"github.com/jmoiron/sqlx"
	"sort"

	"github.com/pkg/errors"
	// pq is the driver for the postgres dialect
	_ "github.com/lib/pq"

	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
)

// MigrationRepository is a repository for the migration entity
type MigrationRepository struct {
	repository
}

// MaxLIstLimit const
const MaxLIstLimit = 1000


var _ migration.IRepository = (*MigrationRepository)(nil)

// NewMigrationRepository creates a new Repository
func NewMigrationRepository(repository *repository) (*MigrationRepository, error) {
	return &MigrationRepository{repository: *repository}, nil
}

func (r MigrationRepository) SetLogger(logger api.Logger) {
	r.logger = logger
}

// get reads entities with the specified ID from the database.
func (r MigrationRepository) get(ctx context.Context, tx *sqlx.Tx, id uint) (*migration.MigrationLog, error) {
	entity := &migration.MigrationLog{}

	err := tx.GetContext(ctx, entity, "SELECT * FROM migration WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity, api.ErrNotFound
		}
	}
	return entity, err
}


func (r MigrationRepository) Query(ctx context.Context, offset, limit uint) ([]migration.MigrationLog, error) {
	var items []migration.MigrationLog

	if limit == 0 {
		limit = MaxLIstLimit
	}

	err := r.db.DB().SelectContext(ctx, &items, "SELECT * FROM migration ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrNotFound
		}
		return nil, errors.Wrapf(api.ErrInternal, "MigrationRepository.Query error: %v", err)
	}
	return items, nil
}

// QueryTx retrieves records with the specified offset and limit from the database.
func (r MigrationRepository) QueryTx(ctx context.Context, t migration.Transaction, query *migration.QueryCondition, offset, limit uint) ([]migration.MigrationLog, error) {
	var items []migration.MigrationLog
	var where string

	tx, ok := t.(*sqlx.Tx)
	if !ok {
		return items, errors.New("can not assert param t migration.Transaction to *sqlx.Tx")
	}

	if limit < 1 {
		limit = MaxLIstLimit
	}
	params := []interface{}{limit, offset}

	if query != nil && query.Where != nil {
		where = " WHERE status = $3 "
		params = append(params, query.Where.Status)
	}

	err := tx.SelectContext(ctx, &items, "SELECT * FROM migration" + where + " LIMIT $1 OFFSET $2", params...)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, api.ErrNotFound
		}
	}
	return items, err
}

// Last retrieves a last record with the specified query condition and limit 1 from the database.
func (r MigrationRepository) Last(ctx context.Context, query *migration.QueryCondition) (*migration.MigrationLog, error) {
	var where string

	params := []interface{}{}

	if query != nil && query.Where != nil {
		where = " WHERE status = $1 "
		params = append(params, query.Where.Status)
	}
	entity := &migration.MigrationLog{}

	err := r.db.DB().GetContext(ctx, entity, "SELECT * FROM migration" + where + " ORDER BY id DESC LIMIT 1", params...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrNotFound
		}
	}
	return entity, nil
}

// LastTx retrieves a last record with the specified query condition and limit 1 from the database.
func (r MigrationRepository) LastTx(ctx context.Context, t migration.Transaction, query *migration.QueryCondition) (*migration.MigrationLog, error) {
	var where string

	tx, ok := t.(*sqlx.Tx)
	if !ok {
		return nil, errors.New("can not assert param t migration.Transaction to *sqlx.Tx")
	}

	params := []interface{}{}

	if query != nil && query.Where != nil {
		where = " WHERE status = $1 "
		params = append(params, query.Where.Status)
	}
	entity := &migration.MigrationLog{}

	err := tx.GetContext(ctx, entity, "SELECT * FROM migration" + where + " ORDER BY id DESC LIMIT 1", params...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrNotFound
		}
	}
	return entity, nil
}

// BatchCreateTx saves a batch of a new entities in the database.
func (r MigrationRepository) BatchCreateTx(ctx context.Context, t migration.Transaction, list migration.MigrationsLogsList) error {
	tx, ok := t.(*sqlx.Tx)
	if !ok {
		return errors.New("can not assert param t migration.Transaction to *sqlx.Tx")
	}

	ids := list.GetIDs()
	sort.Ints(ids)

	for _, i := range ids {
		id := uint(i)
		item := list[id]
		err := r.create(ctx, tx, &item)
		if err != nil {
			return errors.Wrapf(api.ErrInternal, "error while creating log of migration #%v", id)
		}
		list[id] = item
	}
	return nil
}


// BatchUpdateTx updates records of a batch entities in the database.
func (r MigrationRepository) BatchUpdateTx(ctx context.Context, t migration.Transaction, list migration.MigrationsLogsList) error {
	tx, ok := t.(*sqlx.Tx)
	if !ok {
		return errors.New("can not assert param t migration.Transaction to *sqlx.Tx")
	}

	ids := list.GetIDs()
	sort.Ints(ids)

	for _, i := range ids {
		id := uint(i)
		item := list[id]
		err := r.update(ctx, tx, &item)
		if err != nil {
			return errors.Wrapf(api.ErrInternal, "error while updating log of migration #%v", id)
		}
		list[id] = item
	}
	return nil
}

// create saves a new entity in the database.
func (r MigrationRepository) create(ctx context.Context, tx *sqlx.Tx, entity *migration.MigrationLog) error {
	var lastInsertID uint

	err := tx.QueryRowContext(ctx, "INSERT INTO migration (id, status, \"name\", \"time\") VALUES ($1, $2, $3, Now()) RETURNING id", entity.ID, entity.Status, entity.Name).Scan(&lastInsertID)
	if err != nil {
		return errors.Wrapf(err, "MigrationRepository: error inserting entity %v", entity)
	}

	newEntity, err := r.get(ctx, tx, lastInsertID)
	if err != nil {
		return errors.Wrapf(err, "MigrationRepository: error inserting entity %v", entity)
	}
	*entity = *newEntity

	return nil
}

// update recoprd of entity in db
func (r MigrationRepository) update(ctx context.Context, tx *sqlx.Tx, entity *migration.MigrationLog) error {

	_, err := tx.ExecContext(ctx, "UPDATE migration SET status = $1, \"name\" = $2, \"time\" = Now() WHERE id = $3", entity.Status, entity.Name, entity.ID)
	if err != nil {
		return errors.Wrapf(err, "MigrationRepository: error updating entity %v", entity)
	}

	newEntity, err := r.get(ctx, tx, entity.ID)
	if err != nil {
		return errors.Wrapf(err, "MigrationRepository: error inserting entity %v", entity)
	}
	*entity = *newEntity

	return nil
}

// delete deletes a record with the specified ID from the database.
/*func (r MigrationRepository) delete(ctx context.Context, tx *sqlx.Tx, id uint) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM migration WHERE id = $1", id)
	if err != nil {
		return errors.Wrapf(err, "MigrationRepository: error deleting record id = %v", id)
	}
	return nil
}*/


func (r MigrationRepository) BeginTx(ctx context.Context) (migration.Transaction, error) {
	return r.db.DB().BeginTxx(ctx, nil)
}


func (r MigrationRepository) ExecSQL(ctx context.Context, sql string) error {
	tx, err := r.db.DB().BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "MigrationRepository.ExecSQL: transaction begin error")
	}

	_, err = tx.ExecContext(ctx, sql)
	if err != nil {
		er := errors.Wrapf(api.ErrUsersSQL, "MigrationRepository.ExecSQL error: %v", err)
		err = tx.Rollback()
		if err != nil {
			return errors.Wrapf(err, "MigrationRepository.ExecSQL tx.Rollback() error")
		}
		return er
	}
	err = tx.Commit()
	if err != nil {
		return errors.Wrapf(err, "MigrationRepository.ExecSQL tx.Commit() error")
	}

	return nil
}


func (r MigrationRepository) ExecFunc(ctx context.Context, f api.MigrationFunc) (returnErr error) {
	tx, err := r.db.DB().BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "MigrationRepository.ExecFunc: transaction begin error")
	}

	defer func() {
		if err := recover(); err != nil {
			returnErr = errors.Wrapf(api.ErrUsersFunc, "MigrationRepository.ExecFunc error: %v", err)
			err = tx.Rollback()
			if err != nil {
				returnErr = errors.Errorf("MigrationRepository.ExecFunc: transaction rollback error: %v", err)
			}
		}
	}()

	err = f(tx)
	if err != nil {
		er := errors.Wrapf(api.ErrUsersFunc, "MigrationRepository.ExecFunc error: %v", err)
		err = tx.Rollback()
		if err != nil {
			return errors.Wrapf(err, "MigrationRepository.ExecFunc tx.Rollback() error")
		}
		return er
	}
	err = tx.Commit()
	if err != nil {
		return errors.Wrapf(err, "MigrationRepository.ExecFunc tx.Commit() error")
	}

	return nil
}


func (r MigrationRepository) ExecSQLTx(ctx context.Context, t migration.Transaction, sql string) error {
	tx, ok := t.(*sqlx.Tx)
	if !ok {
		return errors.New("can not assert param t migration.Transaction to *sqlx.Tx")
	}

	_, err := tx.ExecContext(ctx, sql)
	if err != nil {
		return errors.Wrapf(api.ErrUsersSQL, "MigrationRepository.ExecSQL error: %v", err)
	}

	return nil
}


func (r MigrationRepository) ExecFuncTx(ctx context.Context, t migration.Transaction, f api.MigrationFunc) (returnErr error) {
	tx, ok := t.(*sqlx.Tx)
	if !ok {
		return errors.New("can not assert param t migration.Transaction to *sqlx.Tx")
	}

	defer func() {
		if err := recover(); err != nil {
			returnErr = errors.Wrapf(api.ErrUsersFunc, "MigrationRepository.ExecFunc error: %v", err)
		}
	}()

	err := f(tx)
	if err != nil {
		return errors.Wrapf(api.ErrUsersFunc, "MigrationRepository.ExecFunc error: %v", err)
	}

	return nil
}

