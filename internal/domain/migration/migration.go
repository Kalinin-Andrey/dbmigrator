package migration

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jmoiron/sqlx"

	"github.com/Kalinin-Andrey/dbmigrator/internal/pkg/apperror"
)

// MigrationTypeSQL - MigrationType gor SQL
const MigrationTypeSQL = "sql"
// MigrationTypeGo - MigrationType gor go
const MigrationTypeGo = "go"

// MigrationTypes is slice of migration types
var MigrationTypes = []interface{}{MigrationTypeSQL, MigrationTypeGo}

// CreateParams is struct for params for creation of migration
type CreateParams struct {
	ID		uint
	Type	string
	Name	string
}

// Validate method
func (p CreateParams) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ID, validation.Required),
		validation.Field(&p.Type, validation.Required, validation.In(MigrationTypes...)),
		validation.Field(&p.Name, validation.Required, validation.Length(1, 100), validation.Match(regexp.MustCompile("^[a-zA-Z0-9_-]+$"))),
	)
}

// Migration struct
// Up and Down is a Func or a string (plain SQL text)
type Migration struct {
	ID		uint
	Name	string
	Up		interface{}
	Down	interface{}
}

// Func is func for migrations Up/Down
type Func func(tx *sqlx.Tx) error

var migrationRule = []validation.Rule{
	validation.NotNil,
	validation.Required,
	validation.By(migrationFuncOrStringRule),
}

// Validate method
func (m Migration) Validate() error {

	err := validation.ValidateStruct(&m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.Name, validation.Required, validation.RuneLength(2, 100), is.UTFLetter),
		validation.Field(&m.Up, migrationRule...),
		validation.Field(&m.Down, migrationRule...),
	)
	return err
}

func migrationFuncOrStringRule(value interface{}) (err error) {
	switch value.(type) {
	case string:
	case Func:
	default:
		err = apperror.ErrUndefinedTypeOfAction
	}
	return err
}

// Log returns corresponding Log
func (m Migration) Log (status uint) *Log {
	return &Log{
		ID:		m.ID,
		Status:	status,
		Name:	m.Name,
	}
}

// MigrationsList ia a map of Migration
type MigrationsList map[uint]Migration

// IDs returns slice of IDs
func (l MigrationsList) IDs() (ids []int) {
	ids = make([]int, 0, len(l))

	for id := range l {
		ids = append(ids, int(id))
	}
	return ids
}


