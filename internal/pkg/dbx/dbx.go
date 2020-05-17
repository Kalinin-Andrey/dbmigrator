package dbx

import (
	"fmt"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"
	"time"

	"github.com/jmoiron/sqlx"
	// pq is the driver for the postgres dialect
	_ "github.com/lib/pq"
)

const ConnectionTimeout = time.Duration(30 * time.Second)

// DBx is the interface for a DB connection
type DBx interface {
	DB() *sqlx.DB
}

// DBx is the struct for a DB connection
type DB struct {
	db *sqlx.DB
}

// DB returns a db object
func (db *DB) DB() *sqlx.DB {
	return db.db
}

// Close connection
func (db *DB) Close() error {
	return db.db.Close()
}

var _ DBx = (*DB)(nil)

var defaultTimeout time.Duration = 10 * time.Second

// New creates a new DB connection
func New(conf api.Configuration, timeout *time.Duration) (*DB, error) {
	if timeout == nil {
		timeout = &defaultTimeout
	}
	db, err := ConnectLoop(conf.Dialect, conf.DSN, *timeout)

	if err != nil {
		return nil, err
	}

	dbobj := &DB{db: db}

	return dbobj, nil
}

// ConnectLoop is the func for connection in a loop with timeout
func ConnectLoop(dialect string, dsn string, timeout time.Duration) (*sqlx.DB, error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("DB connection failed after %s timeout", timeout)

		case <-ticker.C:
			db, err := sqlx.Connect(dialect, dsn)
			if err == nil {
				return db, nil
			}
			//errors.Wrapf(err, "Can not connect to DB %s by DSN: %q", dialect, dsn)
		}
	}
}

