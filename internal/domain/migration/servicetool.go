package migration

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"
)

// ServiceTool is the service for DBMigrator as a tool
type ServiceTool struct {
	*Service
}

var _ IService = (*ServiceTool)(nil)

// NewServiceTool creates a new ServiceTool.
func NewServiceTool(is IService) (*ServiceTool, error) {
	s, ok := is.(*Service)
	if !ok {
		return nil, errors.Errorf("migration.NewServiceTool() error: can not assert to Service input parameter.")
	}
	return &ServiceTool{s}, nil
}

// CreateMainFile creates a main file for migrations execution
func (s ServiceTool) CreateMainFile(ctx context.Context, wr io.Writer) (err error) {
	_, err = fmt.Fprint(wr, mainFileContent)
	return err
}

// Create creates a file for migration
func (s ServiceTool) Create(ctx context.Context, wr io.Writer, p CreateParams) (err error) {
	if err = p.Validate(); err != nil {
		return errors.Wrapf(err, "Invalid create params")
	}

	return s.getTemplate().ExecuteTemplate(wr, "tpl", p)
}

func (s ServiceTool) getTemplate() (*template.Template) {
	return template.Must(template.New("tpl").Parse(`
package main

import (
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"{{if (eq .Type "go")}}
	"github.com/jmoiron/sqlx"{{end}}
)

func init() {
	sqlmigrator.Add(api.Migration{
		ID:		{{.ID}},
		Name:	"{{.Name}}",
		Up:		{{if (eq .Type "sql")}}""{{else}}api.MigrationFunc(func(tx *sqlx.Tx) error {
			_, err := tx.Exec("CREATE TABLE IF NOT EXISTS public.test01(id int4)")	// for example
			return err
		}){{end}},
		Down:	{{if (eq .Type "sql")}}""{{else}}api.MigrationFunc(func(tx *sqlx.Tx) error {
			_, err := tx.Exec("DROP TABLE public.test01")							// for example
			return err
		}){{end}},
	})
}

`))
}

const mainFileContent = `package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"os"

	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator/api"
)

const (
	actionUp		= "up"
	actionDown		= "down"
	actionRedo		= "redo"
)

type config struct {
	dsn		string
	action	string
}

var c config

func init() {
	flag.StringVar(&c.dsn, "dsn", "", "DSN of DB connection")
	flag.StringVar(&c.action, "action", "", "Migration action")
}

func main() {
	flag.Parse()
	conf := api.Configuration{
		DSN:     c.dsn,
		Dir:     ".",
	}
	err := dbmigrator.Init(context.Background(), conf, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch c.action {
	case actionUp:
		err = dbmigrator.Up(0)
	case actionDown:
		err = dbmigrator.Down(0)
	case actionRedo:
		err = dbmigrator.Redo()
	default:
		err = errors.Errorf("Invalid action %q.", c.action)
	}
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

`

