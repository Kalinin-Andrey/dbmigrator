package api

import (
	"os"
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)


const MigrationTypeSQL = "sql"
const MigrationTypeGo = "go"

var MigrationTypes []interface{} = []interface{}{MigrationTypeSQL, MigrationTypeGo}

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


var MigrationStatuses []string = []string{"not applied", "applied", "error"}


type MigrationsList map[uint]Migration


func (l MigrationsList) GetIDs() (ids []int) {
	ids = make([]int, 0, len(l))

	for id := range l {
		ids = append(ids, int(id))
	}
	return ids
}



