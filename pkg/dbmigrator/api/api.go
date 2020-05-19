package api

import (
	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
	"github.com/Kalinin-Andrey/dbmigrator/internal/pkg/dbx"
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

func (c *Configuration) DBConf() *dbx.Configuration {
	return &dbx.Configuration{
		DSN:     c.DSN,
		Dir:     c.Dir,
		Dialect: c.Dialect,
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
	return &migration.Migration{
		ID:   m.ID,
		Name: m.Name,
		Up:   m.Up,
		Down: m.Down,
	}
}



