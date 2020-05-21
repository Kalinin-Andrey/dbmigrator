package gomigration

import (
	"bytes"
	"encoding/json"
	"github.com/Kalinin-Andrey/dbmigrator/internal/domain/migration"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

type File struct {
	Dir			string
	Filename	string
}

func (f File) Validate() error {
	filePath := filepath.Join(f.Dir, filepath.Base(f.Filename))

	if _, err := os.Stat(filePath); err != nil {
		errors.Errorf("Not exists migration file %q", filePath)
	}
	return nil
}

func (f File) Run() (*migration.Migration, error) {
	var buf bytes.Buffer
	filePath := filepath.Join(f.Dir, filepath.Base(f.Filename))
	cmd := exec.Command("go", "run", filePath)
	//cmd.Stdout = os.Stdout
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		errors.Wrapf(err, "gomigration.File.Run() failed, migration file: %q", filePath)
	}
	m := &migration.Migration{}
	err := json.Unmarshal(buf.Bytes(), m)
	if err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal() error")
	}
	return m, nil
}

func DirExec(dir string) ([]migration.Migration, error) {
	files, err := GetFiles(dir)
	if err != nil {
		return nil, err
	}
	ms := make([]migration.Migration, 0)

	for _, f := range files {
		m, err := f.Run()
		if err != nil {
			return nil, err
		}
		ms = append(ms, *m)
	}
	return ms, nil
}

func GetFiles(dir string) ([]File, error) {
	files := make([]File, 0)
	err := filepath.Walk(dir, func(name string, info os.FileInfo, err error) error {
		ext := path.Ext(name)
		if ext != ".go" {
			return nil
		}

		f := File{
			Dir:      dir,
			Filename: name,
		}
		if err := f.Validate(); err != nil {
			return errors.Wrapf(err, "Invalid file %q", filepath.Join(f.Dir, filepath.Base(f.Filename)))
		}
		files = append(files, f)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
