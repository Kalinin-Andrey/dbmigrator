package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
	"github.com/Kalinin-Andrey/dbmigrator/internal/test/fixture"
	"github.com/Kalinin-Andrey/dbmigrator/internal/test/mock"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator/api"
)

const Dir = "."


var SQLTpl = "\npackage migration\n\nimport (\n\t\"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator\"\n\t\"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api\"\n)\n\nfunc init() {\n\tsqlmigrator.Add(api.Migration{\n\t\tID:\t\t8,\n\t\tName:\t\"test_sql\",\n\t\tUp:\t\t\"\",\n\t\tDown:\t\"\",\n\t})\n}\n\n"
var GoTpl = "\npackage migration\n\nimport (\n\t\"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator\"\n\t\"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api\"\n\t\"github.com/jmoiron/sqlx\"\n)\n\nfunc init() {\n\tsqlmigrator.Add(api.Migration{\n\t\tID:\t\t9,\n\t\tName:\t\"test_go\",\n\t\tUp:\t\tapi.MigrationFunc(func(tx *sqlx.Tx) error {\n\t\t\t_, err := tx.Exec(\"CREATE TABLE IF NOT EXISTS public.test01(id int4)\")\t// for example\n\t\t\treturn err\n\t\t}),\n\t\tDown:\tapi.MigrationFunc(func(tx *sqlx.Tx) error {\n\t\t\t_, err := tx.Exec(\"DROP TABLE public.test01\")\t\t\t\t\t\t\t// for example\n\t\t\treturn err\n\t\t}),\n\t})\n}\n\n"


func getSQLMigrator() (*dbmigrator.DBMigrator, error) {
	return dbmigrator.NewSQLMigrator(
		context.Background(),
		api.Configuration{
			Dir:	Dir,
		},
		nil,
		mock.NewMigrationRepository(),
		*fixture.MigrationsList,
	)
}


func TestUp(t *testing.T) {
	mls := fixture.MigrationsLogsList.Copy()
	ms	:= fixture.MigrationsList

	for id, m := range *ms {
		if ml, ok := mls[id]; ok && ml.Status == migration.StatusApplied {
			continue
		}
		mls[id] = *m.Log(migration.StatusApplied)
	}

	m, err := getSQLMigrator()
	if err != nil {
		t.Fatalf("test.getSQLMigrator() error: %v", err)
	}

	err = m.Up(0)
	if err != nil {
		t.Fatalf("sqlmigrator.Up() error: %v", err)
	}

	if !reflect.DeepEqual(*fixture.MigrationsLogsList, mls) {
		t.Errorf("sqlmigrator.Up() result do not much; expected: %v, have: %v", mls, fixture.MigrationsLogsList)
	}
}


func TestDown(t *testing.T) {
	mls := fixture.MigrationsLogsList.Copy()

	if len(mls) == 0 {
		t.Fatalf("fixture.MigrationsLogsList is empty")
	}
	mls = mock.FilterMigrationsLogsByStatus(mls, migration.StatusApplied)
	sl := mls.GetSlice()
	sort.Sort(migration.MigrationsLogsSlice(sl))
	lastAppliedID := sl[len(sl) - 1].ID

	ml := mls[lastAppliedID]
	ml.Status = migration.StatusNotApplied
	mls[lastAppliedID] = ml

	m, err := getSQLMigrator()
	if err != nil {
		t.Fatalf("test.getSQLMigrator() error: %v", err)
	}

	err = m.Down(0)
	if err != nil {
		t.Fatalf("sqlmigrator.Down() error: %v", err)
	}

	if !reflect.DeepEqual(*fixture.MigrationsLogsList, mls) {
		t.Errorf("sqlmigrator.Down() result do not much; expected: %v, have: %v", mls, fixture.MigrationsLogsList)
	}
}


func TestRedo(t *testing.T) {
	mls := fixture.MigrationsLogsList.Copy()

	m, err := getSQLMigrator()
	if err != nil {
		t.Fatalf("test.getSQLMigrator() error: %v", err)
	}

	err = m.Redo()
	if err != nil {
		t.Fatalf("sqlmigrator.Redo() error: %v", err)
	}

	if !reflect.DeepEqual(*fixture.MigrationsLogsList, mls) {
		t.Errorf("sqlmigrator.Redo() result do not much; expected: %v, have: %v", mls, fixture.MigrationsLogsList)
	}
}


func TestStatus(t *testing.T) {
	var expectedList []migration.MigrationLog
	mls := *fixture.MigrationsLogsList

	if len(mls) > 0 {
		expectedList = mls.GetSlice()
		sort.Sort(migration.MigrationsLogsSlice(expectedList))
	}

	m, err := getSQLMigrator()
	if err != nil {
		t.Fatalf("test.getSQLMigrator() error: %v", err)
	}

	list, err := m.Status()
	if err != nil {
		t.Fatalf("sqlmigrator.Status() error: %v", err)
	}

	if !reflect.DeepEqual(list, expectedList) {
		t.Errorf("sqlmigrator.Status() result do not much; expected: %v, have: %v", expectedList, list)
	}
}


func TestDBVersion(t *testing.T) {
	var expectedDBVersion uint
	mls := *fixture.MigrationsLogsList

	if len(mls) > 0 {
		mls = mock.FilterMigrationsLogsByStatus(mls, migration.StatusApplied)
		sl := mls.GetSlice()
		sort.Sort(migration.MigrationsLogsSlice(sl))
		expectedDBVersion = sl[len(sl) - 1].ID
	}

	m, err := getSQLMigrator()
	if err != nil {
		t.Fatalf("test.getSQLMigrator() error: %v", err)
	}

	DBVersion, err := m.DBVersion()
	if err != nil {
		t.Fatalf("sqlmigrator.DBVersion() error: %v", err)
	}

	if DBVersion != expectedDBVersion {
		t.Errorf("sqlmigrator.DBVersion() result do not much; expected: %v, have: %v", expectedDBVersion, DBVersion)
	}
}


func TestCreateSQL(t *testing.T) {
	m, err := getSQLMigrator()
	if err != nil {
		t.Fatalf("test.getSQLMigrator() error: %v", err)
	}
	p := migration.MigrationCreateParams{
		ID:   8,
		Type: "sql",
		Name: "test_sql",
	}
	fileName := fmt.Sprintf("%03d", p.ID) + "_" + p.Name +".go"
	fileName = filepath.Join(Dir, fileName)

	err = m.Create(p)
	if err != nil {
		t.Fatalf("DBMigrator.Create() error: %v", err)
	}
	defer os.Remove(fileName)

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		t.Fatalf("os.OpenFile() error: %v", err)
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("ioutil.ReadAll() error: %v", err)
	}
	content := string(bs)

	if content != SQLTpl {
		t.Errorf("sqlmigrator.Create() result do not much; expected: %v, have: %v", SQLTpl, content)
	}
}



func TestCreateGo(t *testing.T) {
	m, err := getSQLMigrator()
	if err != nil {
		t.Fatalf("test.getSQLMigrator() error: %v", err)
	}
	p := migration.MigrationCreateParams{
		ID:   9,
		Type: "go",
		Name: "test_go",
	}
	fileName := fmt.Sprintf("%03d", p.ID) + "_" + p.Name +".go"
	fileName = filepath.Join(Dir, fileName)

	err = m.Create(p)
	if err != nil {
		t.Fatalf("DBMigrator.Create() error: %v", err)
	}
	defer os.Remove(fileName)

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		t.Fatalf("os.OpenFile() error: %v", err)
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("ioutil.ReadAll() error: %v", err)
	}
	content := string(bs)

	if content != GoTpl {
		t.Errorf("sqlmigrator.Create() result do not much; expected: %v, have: %v", GoTpl, content)
	}
}


