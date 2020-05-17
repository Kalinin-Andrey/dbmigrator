package migration

import (
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"
)

func init() {
	sqlmigrator.Add(api.Migration{
		ID:		1,
		Name:	"init",
		Up:		"CREATE TABLE IF NOT EXISTS public.test01(id int4)",
		Down:	"DROP TABLE public.test01",
	})
}


