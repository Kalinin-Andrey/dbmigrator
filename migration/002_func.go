package migration

import (
	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator"
	"github.com/jmoiron/sqlx"
)

func init() {
	dbmigrator.Add(migration.Migration{
		ID:		2,
		Name:	"func",
		Up:		migration.MigrationFunc(func(tx *sqlx.Tx) error {
			_, err := tx.Exec("CREATE TABLE IF NOT EXISTS public.test02(id int4)")
			return err
		}),
		Down:	migration.MigrationFunc(func(tx *sqlx.Tx) error {
			_, err := tx.Exec("DROP TABLE public.test02")
			return err
		}),
	})
}


