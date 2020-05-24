package migration

import (
	"time"

)

// Log struct
type Log struct {
	ID				uint
	Status			uint
	Name			string
	Time			time.Time
}

// LogsList is a map of Log entities
type LogsList map[uint]Log
// LogsSlice is a slice of Log entities
type LogsSlice []Log


const (
	// TableName const
	TableName        = "dbmigrator_migration"
	// StatusNotApplied const
	StatusNotApplied = 0
	// StatusApplied const
	StatusApplied    = 1
	// StatusError const
	StatusError      = 2
)

// SQLCreateTable is the SQL text for creation table
var SQLCreateTable string = `CREATE TABLE IF NOT EXISTS public."` + TableName + `" (
	id int4 NOT NULL,
	status int4 NOT NULL DEFAULT 0,
	name varchar(100) NOT NULL,
	"time" timestamptz NOT NULL DEFAULT Now(),
	CONSTRAINT migration_pkey PRIMARY KEY (id)
);`

// QueryCondition struct for defining a query condition
type QueryCondition struct {
	Where	*WhereCondition
}
// WhereCondition struct
type WhereCondition struct {
	Status	uint
}

// Slice converts LogsList to slice
func (l LogsList) Slice() (mls []Log) {
	mls = make([]Log, 0, len(l))

	for _, ml := range l {
		mls = append(mls, ml)
	}

	return mls
}

// IDs returns slice of id
func (l LogsList) IDs() (ids []int) {
	ids = make([]int, 0, len(l))

	for id := range l {
		ids = append(ids, int(id))
	}
	return ids
}

// Copy one LogsList to another
func (l LogsList) Copy() (mls LogsList) {
	mls = make(LogsList, len(l))

	for id, ml := range l {
		mls[id] = ml
	}

	return mls
}

// GroupLogsByStatus groups logs applied/not applied
func GroupLogsByStatus(list []Log) (l map[uint]LogsList) {
	l = make(map[uint]LogsList, 2)
	l[StatusNotApplied] = make(LogsList)
	l[StatusApplied] = make(LogsList)

	for _, i := range list {
		if i.Status == StatusApplied {
			l[StatusApplied][i.ID] = i
		} else {
			l[StatusNotApplied][i.ID] = i
		}
	}

	return l
}

// MigrationsListFilterExceptByKeys returns MigrationsList with all entities from sourceList other than those represented in exceptList
func MigrationsListFilterExceptByKeys(sourceList MigrationsList, exceptList LogsList) (l MigrationsList) {
	l = make(MigrationsList)

	for id, m := range sourceList {
		if _, ok := exceptList[id]; !ok {
			l[id] = m
		}
	}

	return l
}

// MigrationsListFilterExistsByKeys returns MigrationsList with all entities from sourceList that represented in existList
func MigrationsListFilterExistsByKeys(sourceList MigrationsList, existList LogsList) (l MigrationsList) {
	l = make(MigrationsList)

	for id, m := range sourceList {
		if _, ok := existList[id]; ok {
			l[id] = m
		}
	}

	return l
}

// MigrationsLogsFilterExceptByKeys returns LogsList with all entities from sourceList other than those represented in exceptList
func MigrationsLogsFilterExceptByKeys(sourceList LogsList, exceptList LogsList) (l LogsList) {
	l = make(LogsList)

	for id, m := range sourceList {
		if _, ok := exceptList[id]; !ok {
			l[id] = m
		}
	}

	return l
}

// MigrationsLogsFilterExistsByKeys returns LogsList with all entities from sourceList that represented in existList
func MigrationsLogsFilterExistsByKeys(sourceList LogsList, existList LogsList) (l LogsList) {
	l = make(LogsList)

	for id, m := range sourceList {
		if _, ok := existList[id]; ok {
			l[id] = m
		}
	}

	return l
}

// Len returns length
func (s LogsSlice) Len() int {
	return len(s)
}

// Swap swaps elements
func (s LogsSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less returns true if i elements less than j elements, otherwise - false
func (s LogsSlice) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

