package fixture

import (
	"time"

	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
)

// MigrationsLogsList fixture
var MigrationsLogsList = &migration.LogsList{
	1:	migration.Log{
			ID:     1,
			Status: migration.StatusApplied,
			Name:   "first_migration",
			Time:   time.Now(),
		},
	2:	migration.Log{
			ID:     2,
			Status: migration.StatusError,
			Name:   "second_migration",
			Time:   time.Now(),
		},
}

