package fixture

import (
	"time"

	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
)

var MigrationsLogsList *migration.MigrationsLogsList = &migration.MigrationsLogsList{
	1:	migration.MigrationLog{
			ID:     1,
			Status: migration.StatusApplied,
			Name:   "first_migration",
			Time:   time.Now(),
		},
	2:	migration.MigrationLog{
			ID:     2,
			Status: migration.StatusError,
			Name:   "second_migration",
			Time:   time.Now(),
		},
}

