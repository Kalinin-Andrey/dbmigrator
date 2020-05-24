package fixture

import (
	"github.com/jmoiron/sqlx"

	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
)

// MigrationsList fixture
var MigrationsList = &migration.MigrationsList{
	1: migration.Migration{
		ID:   1,
		Name: "first_migration",
		Up:   "CREATE TABLE IF NOT EXISTS public.test01(id int4)",
		Down: "DROP TABLE public.test01",
	},
	2: migration.Migration{
		ID:   2,
		Name: "second_migration",
		Up:   "CREATE TABLE IF NOT EXISTS public.test02(id int4)",
		Down: "DROP TABLE public.test02",
	},
	3: migration.Migration{
		ID:   3,
		Name: "third_migration",
		Up:		migration.Func(func(tx *sqlx.Tx) error {
			_, err := tx.Exec("CREATE TABLE IF NOT EXISTS public.test03(id int4)")
			return err
		}),
		Down:	migration.Func(func(tx *sqlx.Tx) error {
			_, err := tx.Exec("DROP TABLE public.test03")
			return err
		}),
	},
}


