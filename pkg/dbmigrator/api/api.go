package api

import (
	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
	"github.com/Kalinin-Andrey/dbmigrator/internal/pkg/dbx"
	"github.com/jmoiron/sqlx"
	"os"
)


type Logger interface {
	Print(v ...interface{})
	Fatal(v ...interface{})
}

type Configuration struct {
	DSN		string
	Dir		string
	Dialect	string
}


func (config *Configuration) ExpandEnv() {
	config.Dir = os.ExpandEnv(config.Dir)
	config.DSN = os.ExpandEnv(config.DSN)
}

func (c *Configuration) DBxConf() *dbx.Configuration {
	return &dbx.Configuration{
		DSN:     c.DSN,
		Dir:     c.Dir,
		Dialect: c.Dialect,
	}
}


var MigrationTypes = []interface{}{migration.MigrationTypeSQL, migration.MigrationTypeGo}

type MigrationCreateParams struct {
	ID		uint
	Type	string
	Name	string
}

func (p *MigrationCreateParams) CoreParams() *migration.MigrationCreateParams {
	return &migration.MigrationCreateParams{
		ID:		p.ID,
		Type:	p.Type,
		Name:	p.Name,
	}
}

var MigrationStatuses = []string{"not applied", "applied", "error"}


type Migration struct {
	ID		uint
	Name	string
	Up		interface{}
	Down	interface{}
}

func (m Migration) DomainMigration() *migration.Migration {
	var up, down interface{}
	up		= m.Up
	down	= m.Down

	if act, ok := (m.Up).(MigrationFunc); ok {
		up = (migration.MigrationFunc)(act)
	}

	if act, ok := (m.Down).(MigrationFunc); ok {
		down = (migration.MigrationFunc)(act)
	}

	return &migration.Migration{
		ID:   m.ID,
		Name: m.Name,
		Up:   up,
		Down: down,
	}
}

type MigrationFunc func(tx *sqlx.Tx) error

