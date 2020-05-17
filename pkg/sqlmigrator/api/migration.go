package api

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jmoiron/sqlx"

)

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
		err = ErrUndefinedTypeOfAction
	}
	return err
}



