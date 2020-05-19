package migration

import (
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator/api"
)

func init() {
	dbmigrator.Add(api.Migration{
		ID:		1,
		Name:	"init",
		Up:		"CREATE TABLE IF NOT EXISTS public.test01(id int4)",
		Down:	"DROP TABLE public.test01",
	})
}


