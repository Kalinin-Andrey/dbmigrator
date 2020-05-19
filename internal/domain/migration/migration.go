package migration

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jmoiron/sqlx"

	"github.com/Kalinin-Andrey/dbmigrator/internal/pkg/apperror"
)


const MigrationTypeSQL = "sql"
const MigrationTypeGo = "go"

var MigrationTypes = []interface{}{MigrationTypeSQL, MigrationTypeGo}

type MigrationCreateParams struct {
	ID		uint
	Type	string
	Name	string
}

func (p MigrationCreateParams) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ID, validation.Required),
		validation.Field(&p.Type, validation.Required, validation.In(MigrationTypes...)),
		validation.Field(&p.Name, validation.Required, validation.Length(1, 100), validation.Match(regexp.MustCompile("^[a-zA-Z0-9_-]+$"))),
	)
}

// Migration
// Up and Down is a MigrationFunc or a string (plain SQL text)
type Migration struct {
	ID		uint
	Name	string
	Up		interface{}
	Down	interface{}
}

// MigrationFunc is func for migrations Up/Down
type MigrationFunc func(tx *sqlx.Tx) error

var migrationRule = []validation.Rule{
	validation.NotNil,
	validation.Required,
	validation.By(migrationFuncOrStringRule),
}


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
	case MigrationFunc:
	default:
		err = apperror.ErrUndefinedTypeOfAction
	}
	return err
}


func (m Migration) Log (status uint) *MigrationLog {
	return &MigrationLog{
		ID:		m.ID,
		Status:	status,
		Name:	m.Name,
	}
}



type MigrationsList map[uint]Migration


func (l MigrationsList) GetIDs() (ids []int) {
	ids = make([]int, 0, len(l))

	for id := range l {
		ids = append(ids, int(id))
	}
	return ids
}


