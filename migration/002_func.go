package migration

import (
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"
	"github.com/jmoiron/sqlx"
)

func init() {
	sqlmigrator.Add(api.Migration{
		ID:		2,
		Name:	"func",
		Up:		api.MigrationFunc(func(tx *sqlx.Tx) error {
			_, err := tx.Exec("CREATE TABLE IF NOT EXISTS public.test02(id int4)")
			return err
		}),
		Down:	api.MigrationFunc(func(tx *sqlx.Tx) error {
			_, err := tx.Exec("DROP TABLE public.test02")
			return err
		}),
	})
}


