package dbx

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	// pq is the driver for the postgres dialect
	_ "github.com/lib/pq"
)

// Configuration for connection to DB
type Configuration struct {
	DSN		string
	Dir		string
	Dialect	string
}

func (c *Configuration) clearQuotes() {
	c.DSN = strings.Trim(c.DSN, `"`)
}

// DBx is the interface for a DB connection
type DBx interface {
	DB() *sqlx.DB
}

// DB is the struct for a DB connection
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

var defaultTimeout = 10 * time.Second

// New creates a new DB connection
func New(conf Configuration, timeout *time.Duration) (*DB, error) {
	if timeout == nil {
		timeout = &defaultTimeout
	}
	conf.clearQuotes()
	db, err := connectLoop(conf.Dialect, conf.DSN, *timeout)

	if err != nil {
		return nil, err
	}

	dbobj := &DB{db: db}

	return dbobj, nil
}

// connectLoop is the func for connection in a loop with timeout
func connectLoop(dialect string, dsn string, timeout time.Duration) (*sqlx.DB, error) {
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
		}
	}
}

