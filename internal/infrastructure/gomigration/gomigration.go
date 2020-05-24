package gomigration

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Args for execution of go migrations
type Args struct {
	DSN		string
	Action	string
}

// Strings returns representation in slice of strings
func (a Args) Strings() []string {
	return []string{fmt.Sprintf("--action=%s", a.Action), fmt.Sprintf("--dsn=%q", a.DSN)}
}

// Dir of migration for execution
type Dir struct {
	Path			string
}

// Validate Dir
func (d Dir) Validate() error {
	i, err := os.Stat(d.Path)

	if err != nil {
		return errors.Wrapf(err, "Not exists migration dir %q", d.Path)
	}

	if !i.IsDir() {
		return errors.Wrapf(err, "It is not a directory: %q", d.Path)
	}

	return nil
}

// Run Dir
func (d Dir) Run(a Args) (output string, err error) {
	var bufOut bytes.Buffer
	var bufErr bytes.Buffer
	var dir	string

	if !filepath.IsAbs(d.Path) {
		dir = "./" + strings.TrimPrefix(d.Path, "./")
	}
	args := a.Strings()

	cmd := exec.Command("go", "run", dir, args[0], args[1])
	cmd.Stdout = &bufOut
	cmd.Stderr = &bufErr

	if err = cmd.Run(); err != nil {
		err = errors.Wrapf(err, "gomigration.Dir.Run() execution error, migration dir: %q; Stderr: %q", d.Path, bufErr.String())
	} else if bufErr.String() != "" {
		err = errors.Errorf("Stderr: %q", bufErr.String())
	}

	return bufOut.String(), err
}

