package api

import (
	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
	"github.com/Kalinin-Andrey/dbmigrator/internal/pkg/dbx"
	"github.com/jmoiron/sqlx"
	"os"
)

// Logger interface for application
type Logger interface {
	Print(v ...interface{})
	Fatal(v ...interface{})
}

// Configuration struct
type Configuration struct {
	DSN		string
	Dir		string
	Dialect	string
}

// ExpandEnv reads env vars
func (c *Configuration) ExpandEnv() {
	c.Dir = os.ExpandEnv(c.Dir)
	c.DSN = os.ExpandEnv(c.DSN)
}

// DBxConf converts to the dbx configuration
func (c *Configuration) DBxConf() *dbx.Configuration {
	return &dbx.Configuration{
		DSN:     c.DSN,
		Dir:     c.Dir,
		Dialect: c.Dialect,
	}
}

// MigrationTypes is slice of migration types
var MigrationTypes = []interface{}{migration.MigrationTypeSQL, migration.MigrationTypeGo}

// MigrationCreateParams is struct for params for creation of migration
type MigrationCreateParams struct {
	ID		uint
	Type	string
	Name	string
}

// CoreParams converts to core params
func (p *MigrationCreateParams) CoreParams() *migration.CreateParams {
	return &migration.CreateParams{
		ID:		p.ID,
		Type:	p.Type,
		Name:	p.Name,
	}
}

// MigrationStatuses is the slice of the migration statuses
var MigrationStatuses = []string{"not applied", "applied", "error"}

// Migration struct
// Up and Down is a Func or a string (plain SQL text)
type Migration struct {
	ID		uint
	Name	string
	Up		interface{}
	Down	interface{}
}

// CoreMigration converts to core migration
func (m Migration) CoreMigration() *migration.Migration {
	var up, down interface{}
	up		= m.Up
	down	= m.Down

	if act, ok := (m.Up).(MigrationFunc); ok {
		up = (migration.Func)(act)
	}

	if act, ok := (m.Down).(MigrationFunc); ok {
		down = (migration.Func)(act)
	}

	return &migration.Migration{
		ID:   m.ID,
		Name: m.Name,
		Up:   up,
		Down: down,
	}
}

// MigrationFunc type
type MigrationFunc func(tx *sqlx.Tx) error

