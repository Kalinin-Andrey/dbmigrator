package migration

import (
	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator"
)

func init() {
	dbmigrator.Add(migration.Migration{
		ID:		1,
		Name:	"init",
		Up:		"CREATE TABLE IF NOT EXISTS public.test01(id int4)",
		Down:	"DROP TABLE public.test01",
	})
}


